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
	"golang.org/x/net/websocket"
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
	doc.head.Append(Style(NewAttr(), Text(CssReset)))
	doc.head.Append(Style(NewAttr(), Text(CssStyle)))
	doc.head.Append(Script(NewAttr(), Text(JsBundle)))
	return doc
}

func listTopics(app *echo.Echo, store Store, v Volume, cont chan string) {

	app.GET("/", func(c echo.Context) error {
		log.Printf("List Topics(%s)", getHostDomain(c))
		var doc = makeDocument()
		q := store.QueryFunc(QuerySelectTopics, getHostDomain(c))
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

const NoSubject = "[empty subject]"

func ensureSubject(s string) string {
	ts := strings.Trim(s, " ")
	if 0 == len(ts) {
		return NoSubject
	}
	return decodeSubject(ts)
}

func formatTimestamp(t time.Time) string {
	d := time.Since(t)
	return d.Truncate(time.Second).String()
}

func listInTopics(app *echo.Echo, store Store, v Volume, cont chan string) {

	handler := func(c echo.Context) error {
		paramTopic := c.Param("topic")
		q := store.QueryFunc(QuerySelectIntopic,
			getHostDomain(c), paramTopic)
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

	qi := store.QueryFunc(QuerySelectAnswerids, pid)
	qm := store.QueryFunc(QuerySelectAnswers, pid)

	var nas = map[int]Node{}
	var sas = map[int]string{}

	qi(RowCallback(func() {
		aids = append(aids, id)
		nas[id] = Div(ClassAttr("attachment-block"))
		sas[id] = sender
	}), &id, &sender)

	root := Div(ClassAttr(fmt.Sprintf("answer depth-%v", depth)))

	return qm(RowCallback(func() {
		block := Div(
			ClassAttr("answer-block").Set("data-record", strconv.Itoa(id)))
		block.Append(
			Div(ClassAttr("answer-header-block"),
				H2(ClassAttr("answer-subject"),
					A(ClassAttr("link").
						Set("href", fmt.Sprintf("/%s/%v", topic, id)),
						Text(senderName(sender)))),
				A(ClassAttr("answer-link link").
					Set("href", fmt.Sprintf("mailto:%s+%v@%s?subject=Re:%s/%s",
						topic, id, getHostDomain(c), topic, pid)),
					Text("reply"))),
			Div(ClassAttr("answer-body"), NewRawNode(body)), nas[id])
		root.Append(block, formatAnswers(strconv.Itoa(id), store, c, depth+1))
	}), &id, &ts, &sender, &topic, &subject, &body, &parent).
		FoldNodeF(
			func(err error) Node { return Text(err.Error()) },
			func(_ bool) Node {
				attachments := []attachmentRecord{}
				for _, aid := range aids {
					attBlock := nas[aid]
					qa := store.QueryFunc(QuerySelectAttachments, aid)
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
	handler := func(c echo.Context) error {
		paramId := c.Param("id")
		paramTopic := c.Param("topic")
		qm := store.QueryFunc(QuerySelectRecord, paramId)
		qa := store.QueryFunc(QuerySelectAttachments, paramId)
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

		block := Div(ClassAttr("message-block").Set("data-record", paramId))
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
						Set("href", fmt.Sprintf("mailto:%s+%v@%s?subject=Re:%s/%d", paramTopic, id, getHostDomain(c), topic, id)),
						Text("reply"))),
				Div(ClassAttr("message-body"), NewRawNode(body)),
				attBlock)

			doc.body.Append(
				header(ensureSubject(subject), paramTopic, paramId),
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

func showAttachment(app *echo.Echo, store Store, v Volume, cont chan string) {
	handler := func(c echo.Context) error {
		sender := c.Param("sender")
		topic := c.Param("topic")
		id := c.Param("id")
		name := c.Param("name")

		var (
			ct string
			fn string
		)

		store.QueryFunc(QuerySelectAttachment, id, name).Exec(&ct, &fn)
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

func notifyHandler(app *echo.Echo, store Store, v Volume, cont chan string, n *Notifier) {
	handler := func(c echo.Context) error {
		websocket.Handler(func(ws *websocket.Conn) {
			defer ws.Close()
			n.Subscribe(func(i interface{}) {
				switch i.(type) {
				case int:
					websocket.Message.Send(ws,
						fmt.Sprintf("{\"record\": %d}", i.(int)))
				}
			})

			for {
				msg := ""
				err := websocket.Message.Receive(ws, &msg)
				if err != nil {
					break
				}
				log.Printf("websocket: %s", msg)
			}
			// for {
			// 	// Write
			// 	err := websocket.Message.Send(ws, "HELO")
			// 	if err != nil {
			// 		c.Logger().Error(err)
			// 	}

			// 	// Read
			// 	msg := ""
			// 	err = websocket.Message.Receive(ws, &msg)
			// 	if err != nil {
			// 		c.Logger().Error(err)
			// 	}
			// 	fmt.Printf("%s\n", msg)
			// }
		}).ServeHTTP(c.Response(), c.Request())

		return nil
	}
	app.GET("/.notifications", handler)
}

func regHTTPHandlers(app *echo.Echo, store Store, v Volume, cont chan string, n *Notifier) {
	notifyHandler(app, store, v, cont, n)
	listTopics(app, store, v, cont)
	listInTopics(app, store, v, cont)
	showMessage(app, store, v, cont)
	showAttachment(app, store, v, cont)
}

func StartHTTP(cont chan string, iface string, store Store, v Volume, n *Notifier) {
	app := echo.New()
	regHTTPHandlers(app, store, v, cont, n)
	cont <- fmt.Sprintf("HTTP ready on %s", iface)
	app.Start(iface)
}
