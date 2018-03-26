package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"net/mail"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo"
)

func formatMultipart(id int, msg *mail.Message, params map[string]string) Node {
	log.Printf("formatMultipart %v", id)
	node := Div(ClassAttr("message"))
	mc := make(chan SerializedPart)
	boundary := params["boundary"]
	go walkParts(msg.Body, boundary, mc, 0)

	parAttrs := ClassAttr("paragraph")
	for {
		part := <-mc
		if part == nil {
			break
		}
		switch mp := part.MainType(); mp {
		case "text":
			if "plain" == part.SubType() {
				ps := strings.Split(part.ContentString(), "\n")
				for _, p := range ps {
					node.Append(P(parAttrs, Text(p)))
				}
			}
		default:
			{
				fn := part.FileName()
				url := fmt.Sprintf("/attachment/%v/%s", id, fn)
				if "image" == mp {
					node.Append(
						A(ClassAttr("attachment image").Set("href", url).Set("title", fn),
							Img(ClassAttr("").Set("src", url))))
				} else {
					node.Append(A(ClassAttr("attachment link").Set("href", url), Text(fn)))
				}
			}
		}

	}

	return node
}

type MessageInfo struct {
	subject string
	t       time.Time
}

func messageInfoError(err error) MessageInfo {
	return MessageInfo{
		err.Error(),
		time.Now(),
	}
}

func formatMessage(topic string, id int, messageReader *bytes.Reader) (MessageInfo, Node) {
	log.Printf("formatMessage %v", id)
	msg, err := mail.ReadMessage(io.Reader(messageReader))
	if err != nil {
		return messageInfoError(err), Pre(ClassAttr("error"), Text(err.Error()))
	}

	cte := msg.Header.Get(HeaderContentTranferEncoding)
	body, err := ioutil.ReadAll(msg.Body)
	if err != nil {
		return messageInfoError(err), Pre(ClassAttr("error"), Text(err.Error()))
	}

	info := MessageInfo{
		msg.Header.Get("Subject"),
		OptionTimeFrom(msg.Header.Date()).FoldTime(time.Now(), IdTime),
	}
	mediatype, params, err := mime.ParseMediaType(msg.Header.Get(echo.HeaderContentType))

	plain := func() (MessageInfo, Node) {
		return info, Text(string(decodeContent(&body, cte)))
	}

	if err != nil {
		if "" == mediatype {
			return plain()
		}
		return messageInfoError(err), Pre(ClassAttr("error"), Text(err.Error()))
	}

	if isMultipart(mediatype) {
		messageReader.Seek(0, 0)
		return info, formatMultipart(id, msg, params)
	}
	return plain()
}

func sReader(s ...string) io.Reader {
	return strings.NewReader(strings.Join(s, ""))
}

func link(href string, label string) Node {
	return A(ClassAttr("link").Set("href", href), Text(label))
}

func link0(href string) Node {
	return link(href, href)
}

func header(title string, args ...string) Node {
	r := Div(ClassAttr("header"))
	bc := Div(ClassAttr("bc"), link("/", "root"))
	u := ""
	for i := 0; i < len(args); i++ {
		u += "/" + args[i]
		bc.Append(link(u, args[i]))
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

func listTopics(app *echo.Echo, store Store, cont chan string) {
	store.Register("mail/topic-list",
		`SELECT DISTINCT(topic) topic, count(id) FROM {{.RawMails}} GROUP BY topic;`)

	app.GET("/", func(c echo.Context) error {
		var doc = makeDocument()
		q := store.QueryFunc("mail/topic-list")
		var (
			topic string
			count int
		)

		doc.body.Append(header("Topics"))
		attrs := ClassAttr("topic")

		return q(RowCallback(func() {
			doc.body.Append(Div(attrs,
				A(ClassAttr("topic-link link").Set("href", "/"+topic),
					Text(topic)),
				Span(ClassAttr("topic-count"), Textf(" %v", count))))
		}), &topic, &count).
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

func listInTopics(app *echo.Echo, store Store, cont chan string) {
	store.Register("mail/topic-messages",
		`SELECT r.id, r.sender, r.subject, a.parent 
		FROM {{.RawMails}} r
			LEFT JOIN {{.Answers}} a ON r.id = a.child 
		WHERE topic = $1;`)

	handler := func(c echo.Context) error {
		paramTopic := c.Param("topic")
		q := store.QueryFunc("mail/topic-messages", paramTopic)
		var doc = makeDocument()
		var (
			id      int
			sender  string
			subject string
			parent  sql.NullInt64
		)

		doc.body.Append(header(paramTopic, paramTopic))
		attrs := ClassAttr("message-item")

		return q(RowCallback(func() {
			if parent.Valid {
				return
			}
			url := fmt.Sprintf("/%s/%v", paramTopic, id)

			doc.body.Append(Div(attrs,
				A(ClassAttr("message-link link").Set("href", url),
					Text(ensureSubject(subject))),
				Span(ClassAttr("message-item-sender"), Text(senderName(sender)))))
		}), &id, &sender, &subject, &parent).
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

func formatAnswers(topic string, pid string, store Store, c echo.Context, depth int) Node {
	var (
		id      int
		sender  string
		message sql.RawBytes
	)

	root := Div(ClassAttr(fmt.Sprintf("answer depth-%v", depth)))

	return store.QueryFunc(queryAnswers, pid)(RowCallback(func() {
		msgReader := bytes.NewReader(message)
		info, msgNode := formatMessage(topic, id, msgReader)
		block := Div(ClassAttr("answer-block"),
			Div(ClassAttr("answer-header-block"),
				H2(ClassAttr("answer-subject"),
					A(ClassAttr("link").
						Set("href", fmt.Sprintf("/%s/%v", topic, id)),
						Text(info.subject))),
				A(ClassAttr("answer-link").
					Set("href", fmt.Sprintf("mailto:%s+%v@%s", topic, id, c.Request().Host)),
					Text("answer"))),
			Div(ClassAttr("answer-body"), msgNode))
		root.Append(block, formatAnswers(topic, strconv.Itoa(id), store, c, depth+1))
	}), &id, &sender, &message).
		FoldNodeF(
			func(err error) Node { return Text(err.Error()) },
			func(_ bool) Node { return root })
}

func senderName(sender string) string {
	return strings.Split(sender, "@")[0]
}

func showMessage(app *echo.Echo, store Store, cont chan string) {
	store.Register(queryMessage,
		`SELECT  r.id, r.sender, r.message, a.parent 
		FROM {{.RawMails}} r 
			LEFT JOIN {{.Answers}} a ON r.id = a.child  
		WHERE r.id = $1;`)

	store.Register(queryAnswers,
		`SELECT r.id as id, r.sender as sender, r.message as message 
		FROM {{.RawMails}} r 
			LEFT Join {{.Answers}} a ON r.id = a.child
		WHERE a.parent = $1;`)

	handler := func(c echo.Context) error {
		paramId := c.Param("id")
		paramTopic := c.Param("topic")
		q := store.QueryFunc(queryMessage, paramId)
		var doc = makeDocument()
		var (
			id      int
			sender  string
			message sql.RawBytes
			parent  sql.NullInt64
		)

		defer func() { message = []byte{} }()

		return q(RowCallback(func() {
			msgReader := bytes.NewReader(message)
			info, msgNode := formatMessage(paramTopic, id, msgReader)
			pnode := NoDisplay
			if parent.Valid {
				pnode = Div(ClassAttr("parent"),
					A(ClassAttr("link").
						Set("href", fmt.Sprintf("/%s/%v", paramTopic, parent.Int64)),
						Text("parent")))
			}

			block := Div(ClassAttr("message-block"),
				Div(ClassAttr("message-header"),
					Span(ClassAttr("message-sender"),
						Text(fmt.Sprintf("%s - %s", senderName(sender), info.t))),
					pnode,
					A(ClassAttr("link").
						Set("href", fmt.Sprintf("mailto:%s+%v@%s", paramTopic, id, c.Request().Host)),
						Text("answer"))),
				Div(ClassAttr("message-body"), msgNode))

			doc.body.Append(
				header(info.subject, paramTopic, paramId),
				block, formatAnswers(paramTopic, paramId, store, c, 1))
		}), &id, &sender, &message, &parent).
			FoldErrorF(
				func(err error) error {
					return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
				},
				func(bool) error {
					return c.HTML(http.StatusOK, doc.Render())
				})

	}

	app.GET("/:topic/:id", handler)
}

func showAttachment(app *echo.Echo, store Store, cont chan string) {

	store.Register("mail/get-attachment",
		`SELECT id,  message FROM {{.RawMails}} WHERE id = $1;`)
	handler := func(c echo.Context) error {
		paramId := c.Param("id")
		name := c.Param("name")
		log.Printf("showAttachment Query Start")
		rows, err := store.Query("mail/get-attachment", paramId)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		log.Printf("showAttachment Query End")

		var (
			id      int
			message []byte
		)
		defer func() { rows.Close() }()
		defer func() { message = []byte{} }()

		for rows.Next() {
			log.Printf("showAttachment Row Start")

			scanError := rows.Scan(&id, &message)
			if scanError != nil {
				errMesg := scanError.Error()
				cont <- errMesg
				break
			}
			a := getAttachment(name, &message)
			// var gcs debug.GCStats
			// debug.ReadGCStats(&gcs)
			// log.Printf("GC  *v", gcs)
			debug.FreeOSMemory()
			return c.Stream(http.StatusOK, a.mediaType, a.r)
		}

		return c.String(http.StatusNotFound, name)
	}
	app.GET("/attachment/:id/:name", handler)
}

func regHTTPHandlers(app *echo.Echo, store Store, cont chan string) {
	listTopics(app, store, cont)
	listInTopics(app, store, cont)
	showMessage(app, store, cont)
	showAttachment(app, store, cont)
}

func StartHTTP(cont chan string, iface string, store Store) {
	app := echo.New()
	regHTTPHandlers(app, store, cont)
	cont <- fmt.Sprintf("HTTP ready on %s", iface)
	app.Start(iface)
}
