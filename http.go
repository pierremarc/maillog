package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/pgtype"
	"github.com/labstack/echo"
)

func link(href string, label string) Node {
	return A(ClassAttr("link").Set("href", href), Text(label))
}

func link0(href string) Node {
	return link(href, href)
}

func header(title string, args ...string) Node {
	r := Div(ClassAttr("header"))
	bc := Div(ClassAttr("bc"), link("/", "/"+GetSiteName()))
	u := ""
	for i := 0; i < len(args); i++ {
		part := "/" + args[i]
		u += part
		bc.Append(link(u, part))
	}
	r.Append(
		H1(ClassAttr("title"), Text(title)),
		Div(ClassAttr("bc-block"), bc))

	return r
}

func makeDocument() document {
	doc := NewDoc()
	doc.head.Append(HeadMeta(NewAttr().Set("charset", "utf-8")))
	doc.head.Append(HeadMeta(NewAttr().
		Set("name", "viewport").
		Set("content", "width=device-width, initial-scale=1.0")))
	doc.head.Append(Style(NewAttr(), Text(StyleSheet)))
	return doc
}

func listTopics(app *echo.Echo, store Store, v Volume, cont chan string) {
	store.Register("mail/topic-list",
		`SELECT DISTINCT(topic) topic, count(id), max(ts) as mts
		FROM {{.Records}} 
		WHERE strpos(topic, '_') <> 1
        GROUP BY topic
        ORDER BY topic ASC;`)

	app.GET("/", func(c echo.Context) error {
		var doc = makeDocument()
		q := store.QueryFunc("mail/topic-list")
		var (
			topic string
			count int
			mts   pgtype.Timestamp
		)

		doc.body.Append(header("Topics"))
		attrs := ClassAttr("topic")

		return q(RowCallback(func() {
			doc.body.Append(Div(attrs,
				A(ClassAttr("topic-link link").Set("href", "/"+topic),
					Text(topic)),
				Span(ClassAttr("topic-count"), Textf("(%d messages),", count)),
				Span(ClassAttr("topic-ts"), Textf("(last update: %s)", formatTimestamp(mts.Time)))))
		}), &topic, &count, &mts).
			FoldErrorF(
				func(err error) error {
					return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
				},
				func(b bool) error {
					return c.HTML(http.StatusOK, doc.Render())
				})
	})
}

const NoSubject = "***"

func ensureSubject(s string) string {
	ts := strings.Trim(s, " ")
	if 0 == len(ts) {
		return NoSubject
	}
	return ts
}

func formatTimestamp(t time.Time) string {
	d := time.Since(t)
	return d.Truncate(time.Second).String()
}

func listInTopics(app *echo.Echo, store Store, v Volume, cont chan string) {
	store.Register("mail/topic-messages",
		`SELECT id,ts, sender, header_subject
		FROM {{.Records}}  
        WHERE topic = $1 AND parent IS NULL 
        ORDER BY ts DESC;`)

	handler := func(c echo.Context) error {
		paramTopic := c.Param("topic")
		q := store.QueryFunc("mail/topic-messages", paramTopic)
		var doc = makeDocument()
		var (
			id      int
			ts      pgtype.Timestamptz
			sender  string
			subject string
		)

		doc.body.Append(header(paramTopic, paramTopic))
		attrs := ClassAttr("message-item")

		return q(RowCallback(func() {
			url := fmt.Sprintf("/%s/%v", paramTopic, id)

			doc.body.Append(Div(attrs,
				A(ClassAttr("message-link link").Set("href", url),
					Text(ensureSubject(subject))),
				Span(ClassAttr("message-item-sender"), Text(senderName(sender))),
				Span(ClassAttr("message-item-ts"), Textf(formatTimestamp(ts.Time)))))
		}), &id, &ts, &sender, &subject).
			FoldErrorF(
				func(err error) error {
					return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
				},
				func(b bool) error {
					return c.HTML(http.StatusOK, doc.Render())
				})
	}

	app.GET("/:topic", handler)
}

const queryMessage = "mail/get-message"
const queryAnswers = "mail/get-answers"
const queryAnswersIds = "mail/get-answers-ids"
const queryAttachments = "mail/get-attachments"

func bodyToP(body string) []Node {
	pars := NewArrayString(strings.Split(body, "\n"))

	return pars.MapNode(func(p string) Node {
		return P(ClassAttr("body-par"), Text(p))
	}).Slice()
}

func formatAnswers(pid string, store Store, c echo.Context, depth int) Node {
	var (
		id          int
		aids        []int
		ts          pgtype.Timestamptz
		sender      string
		topic       string
		subject     string
		body        string
		parent      pgtype.Int4
		contentType string
		fileName    string
	)

	qi := store.QueryFunc(queryAnswersIds, pid)
	qm := store.QueryFunc(queryAnswers, pid)

	var nas = map[int]Node{}
	var sas = map[int]string{}

	qi(RowCallback(func() {
		aids = append(aids, id)
		nas[id] = Div(ClassAttr("attachment-block"))
		sas[id] = sender
	}), &id, &sender)

	root := Div(ClassAttr(fmt.Sprintf("answer depth-%v", depth)))

	return qm(RowCallback(func() {
		block := Div(ClassAttr("answer-block"))
		block.Append(
			Div(ClassAttr("answer-header-block"),
				H2(ClassAttr("answer-subject"),
					A(ClassAttr("link").
						Set("href", fmt.Sprintf("/%s/%v", topic, id)),
						Text(senderName(sender)))),
				A(ClassAttr("answer-link link").
					Set("href", fmt.Sprintf("mailto:%s+%v@%s?subject=Re:%s/%s",
						topic, id, c.Request().Host, topic, pid)),
					Text("reply"))),
			Div(ClassAttr("answer-body"), bodyToP(body)...), nas[id])
		root.Append(block, formatAnswers(strconv.Itoa(id), store, c, depth+1))
	}), &id, &ts, &sender, &topic, &subject, &body, &parent).
		FoldNodeF(
			func(err error) Node { return Text(err.Error()) },
			func(_ bool) Node {
				attachments := []attachmentRecord{}
				for _, aid := range aids {
					attBlock := nas[aid]
					qa := store.QueryFunc(queryAttachments, aid)
					qa(RowCallback(func() {
						log.Printf("attach to answer %s(%s), %s, %d, %s",
							sas[aid], encodedSender(sas[aid]), topic, aid, fileName)
						attachments = append(attachments, attachmentRecord{
							sender: encodedSender(sas[aid]),
							topic:  topic,
							record: aid,
							ct:     contentType,
							fn:     fileName,
						})
					}), &id, &contentType, &fileName)

					attBlock.Append(formatAttachments(attachments))
					attachments = []attachmentRecord{}
				}
				return root
			})
}

func senderName(sender string) string {
	return strings.Split(sender, "@")[0]
}

type attachmentRecord struct {
	sender string
	topic  string
	record int
	ct     string
	fn     string
}

func formatAttachments(rs []attachmentRecord) Node {
	node := Div(ClassAttr("attachment-block"))

	for _, r := range rs {
		mt := strings.Split(r.ct, "/")[0]
		url := fmt.Sprintf("/attachments/%s/%s/%d/%s", r.sender, r.topic, r.record, r.fn)
		if "image" == mt {
			node.Append(
				A(ClassAttr("attachment image").Set("href", url).Set("title", r.fn),
					Img(ClassAttr("").Set("src", url))))
		} else {
			node.Append(A(ClassAttr("attachment link").Set("href", url), Text(r.fn)))
		}
	}

	return node
}

func showMessage(app *echo.Echo, store Store, v Volume, cont chan string) {
	store.Register(queryMessage,
		`SELECT  id, ts, sender, topic, header_subject, body, parent 
		FROM {{.Records}} 
        WHERE id = $1;`)

	store.Register(queryAnswers,
		`SELECT  id, ts, sender, topic, header_subject, body, parent 
		FROM {{.Records}} 
        WHERE parent = $1
        ORDER BY ts ASC;`)

	store.Register(queryAnswersIds,
		`SELECT id, sender
		FROM {{.Records}}
        WHERE parent = $1
        ORDER BY ts ASC`)

	store.Register(queryAttachments,
		`SELECT record_id, content_type, file_name
		FROM {{.Attachments}}
		WHERE record_id = $1`)

	handler := func(c echo.Context) error {
		paramId := c.Param("id")
		paramTopic := c.Param("topic")
		qm := store.QueryFunc(queryMessage, paramId)
		qa := store.QueryFunc(queryAttachments, paramId)
		var doc = makeDocument()
		var (
			id          int
			ts          pgtype.Timestamptz
			sender      string
			topic       string
			subject     string
			body        string
			parent      pgtype.Int4
			contentType string
			fileName    string
			attachments []attachmentRecord
		)

		block := Div(ClassAttr("message-block"))
		attBlock := Div(ClassAttr("attachment-block"))

		return qm(RowCallback(func() {
			pnode := NoDisplay
			if parent.Status == pgtype.Present {
				pnode = Div(ClassAttr("parent"),
					A(ClassAttr("link").
						Set("href", fmt.Sprintf("/%s/%d", paramTopic, parent.Int)),
						Text("parent")))
			}

			block.Append(
				Div(ClassAttr("message-header"),
					Span(ClassAttr("message-sender"),
						Text(fmt.Sprintf("%s - %s", senderName(sender), ts.Time))),
					pnode,
					A(ClassAttr("link").
						Set("href", fmt.Sprintf("mailto:%s+%v@%s?subject=Re:%s/%d", paramTopic, id, c.Request().Host, topic, id)),
						Text("reply"))),
				Div(ClassAttr("message-body"), bodyToP(body)...),
				attBlock)

			doc.body.Append(
				header(subject, paramTopic, paramId),
				block, formatAnswers(paramId, store, c, 1))
		}), &id, &ts, &sender, &topic, &subject, &body, &parent).
			FoldErrorF(
				func(err error) error {
					return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
				},
				func(bool) error {
					qa(RowCallback(func() {
						log.Printf("attach to record %s(%s), %s, %d, %s",
							sender, encodedSender(sender), topic, id, fileName)
						attachments = append(attachments, attachmentRecord{
							sender: encodedSender(sender),
							topic:  topic,
							record: id,
							ct:     contentType,
							fn:     fileName,
						})
					}), &id, &contentType, &fileName)

					attBlock.Append(formatAttachments(attachments))

					return c.HTML(http.StatusOK, doc.Render())
				})

	}

	app.GET("/:topic/:id", handler)
}

const queryAttachment = "attachment/get-one"

func showAttachment(app *echo.Echo, store Store, v Volume, cont chan string) {
	store.Register(queryAttachment,
		`SELECT content_type, file_name
		FROM {{.Attachments}}
		WHERE record_id = $1 AND file_name = $2`)

	handler := func(c echo.Context) error {
		sender := c.Param("sender")
		topic := c.Param("topic")
		id := c.Param("id")
		name := c.Param("name")

		var (
			ct string
			fn string
		)

		store.QueryFunc(queryAttachment, id, name).Exec(&ct, &fn)
		fp := path.Join(sender, topic, id, fn)
		return v.Reader(fp).FoldErrorF(
			func(err error) error {
				return echo.NewHTTPError(http.StatusNotFound, err.Error())
			},
			func(r io.Reader) error {
				return c.Stream(http.StatusOK, ct, r)
			})

	}
	app.GET("/attachments/:sender/:topic/:id/:name", handler)
}

func regHTTPHandlers(app *echo.Echo, store Store, v Volume, cont chan string) {
	listTopics(app, store, v, cont)
	listInTopics(app, store, v, cont)
	showMessage(app, store, v, cont)
	showAttachment(app, store, v, cont)
}

func StartHTTP(cont chan string, iface string, store Store, v Volume) {
	app := echo.New()
	regHTTPHandlers(app, store, v, cont)
	cont <- fmt.Sprintf("HTTP ready on %s", iface)
	app.Start(iface)
}
