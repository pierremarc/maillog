package main

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"mime/multipart"
	"mime/quotedprintable"
	"net/http"
	"net/mail"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo"
)

const HeaderContentTranferEncoding = "Content-Transfer-Encoding"

func mainContentType(p *multipart.Part) string {
	ct := p.Header.Get("Content-Type")
	ctParts := strings.Split(ct, "/")
	return ctParts[0]
}

type SerializedPart interface {
	ContentType() string
	ContentEncoding() string
	Content() []byte
	ContentString() string
	FileName() string
	MainType() string
	SubType() string
}

type serPart struct {
	contentType     string
	contentEncoding string
	content         []byte
	filename        string
}

func (s serPart) ContentType() string {
	return s.contentType
}
func (s serPart) ContentEncoding() string {
	return s.contentEncoding
}
func (s serPart) Content() []byte {
	return s.content
}
func (s serPart) ContentString() string {
	return fmt.Sprintf("%s", decodeContent(s.content, s.contentEncoding))
}
func (s serPart) FileName() string {
	return s.filename
}

func (s serPart) MainType() string {
	mediaType, _, err := mime.ParseMediaType(s.contentType)
	if err != nil {
		return "any"
	}
	ctParts := strings.Split(mediaType, "/")
	return ctParts[0]
}

func (s serPart) SubType() string {
	mediaType, _, err := mime.ParseMediaType(s.contentType)
	if err != nil {
		return "any"
	}
	ctParts := strings.Split(mediaType, "/")
	return ctParts[1]
}

func serPartF(p *multipart.Part) SerializedPart {
	ct := p.Header.Get(echo.HeaderContentType)
	cte := p.Header.Get(HeaderContentTranferEncoding)
	fn := p.FileName()
	log.Printf("serPartF %s %s", ct, fn)
	input, err := ioutil.ReadAll(p)
	if err != nil {
		return nil
	}

	// content := decodeContent(input, cte)

	return serPart{
		contentType:     ct,
		contentEncoding: cte,
		content:         input,
		filename:        fn,
	}
}

func walkParts(r io.Reader, boundary string, c chan SerializedPart, depth int) {
	log.Printf("walkParts %v", depth)
	reader := multipart.NewReader(r, boundary)
	counter := 0
	for {
		counter = counter + 1
		part, err := reader.NextPart()
		if err != nil {
			log.Printf("Err on NextPart %s", err.Error())
			break
		}

		contentType := part.Header.Get("Content-Type")
		mediaType, params, err := mime.ParseMediaType(contentType)
		if err != nil {
			log.Printf("Err on ParseMediaType %s", err.Error())
			break
		}

		if !isMultipart(mediaType) {
			c <- serPartF(part)
		} else {
			newBoundary := params["boundary"]
			walkParts(part, newBoundary, c, depth+1)
		}
		part.Close()
	}

	if 0 == depth {
		c <- nil
	}
}

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

func decodeContent(input []byte, cte string) []byte {
	log.Printf("decodeContent (%s)", cte)
	switch cte {

	case "base64":
		{
			content, err := base64.StdEncoding.DecodeString(string(input))
			if err != nil {
				return nil
			}
			return content
		}

	case "7bit":
		{
			r := quotedprintable.NewReader(bytes.NewReader(input))
			content, err := ioutil.ReadAll(r)
			if err != nil {
				return nil
			}
			return content
		}
	}

	return input
}

func isMultipart(mediaType string) bool {
	const rfc822 = "message/rfc822"
	const mp = "multipart"
	mainType := strings.Split(mediaType, "/")[0]
	if mediaType == rfc822 || mainType == mp {
		return true
	}
	return false
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
		return info, Text(string(decodeContent(body, cte)))
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

type attachment struct {
	r         io.Reader
	name      string
	mediaType string
}

func errAttach(r io.Reader) attachment {
	return attachment{
		r:         r,
		name:      "Error.txt",
		mediaType: "text/plain",
	}
}

func getAttachment(name string, data *[]byte) attachment {
	reader := bytes.NewReader(*data)
	msg, err := mail.ReadMessage(reader)
	if err != nil {
		return errAttach(sReader("Error Reading Message: ", err.Error()))
	}
	_, params, err := mime.ParseMediaType(msg.Header.Get("Content-Type"))
	if err != nil {
		return errAttach(sReader("Error MediaType: ", err.Error()))
	}
	mc := make(chan SerializedPart)
	boundary := params["boundary"]
	go walkParts(msg.Body, boundary, mc, 0)

	for {
		part := <-mc
		if part == nil {
			return errAttach(sReader("Attachment Not Found ", name))
		}
		mp := part.MainType()
		fn := part.FileName()
		content := decodeContent(part.Content(), part.ContentEncoding())
		if mp != "text" && fn == name {
			return attachment{
				r:         bytes.NewReader(content),
				name:      fn,
				mediaType: part.ContentType(),
			}
		}
	}

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
