package main

import (
	"bytes"
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
	"strings"

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
	ctParts := strings.Split(s.contentType, "/")
	return ctParts[0]
}

func (s serPart) SubType() string {
	ctParts := strings.Split(s.contentType, "/")
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
	}

	if 0 == depth {
		c <- nil
	}
}

func formatMultipart(id int, msg *mail.Message, params map[string]string) Node {
	log.Printf("formatMultipart %v", id)
	node := Div(NewAttr())
	mc := make(chan SerializedPart)
	boundary := params["boundary"]
	go walkParts(msg.Body, boundary, mc, 0)

	parAttrs := NewAttr()
	for {
		part := <-mc
		if part == nil {
			break
		}
		switch mp := part.MainType(); mp {
		case "text":
			if true {
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
						A(NewAttr().Add("href", url).Add("title", fn)),
						Img(NewAttr().Add("src", url)))
				} else {
					node.Append(A(NewAttr().Add("href", url), Text(fn)))
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

func formatMessage(id int, data string) (string, Node) {
	log.Printf("formatMessage %v", id)
	messageReader := strings.NewReader(data)
	msg, err := mail.ReadMessage(messageReader)
	if err != nil {
		return "Server Error", Pre(NewAttr(), Text(err.Error()))
	}

	cte := msg.Header.Get(HeaderContentTranferEncoding)
	body, err := ioutil.ReadAll(msg.Body)
	if err != nil {
		return "Server Error", Pre(NewAttr(), Text(err.Error()))
	}

	subject := msg.Header.Get("Subject")
	mediatype, params, err := mime.ParseMediaType(msg.Header.Get(echo.HeaderContentType))
	if err != nil {
		if "" == mediatype {
			return subject, Text(string(decodeContent(body, cte)))
		}
		return "Message Parsing Error", Pre(NewAttr(), Text(err.Error()))
	}

	if isMultipart(mediatype) {
		messageReader.Seek(0, 0)
		return subject, formatMultipart(id, msg, params)
	}
	return subject, Text(string(decodeContent(body, cte)))

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

func getAttachment(name string, data string) attachment {
	msg, err := mail.ReadMessage(strings.NewReader(data))
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
	return A(NewAttr().Add("href", href), Text(label))
}

func link0(href string) Node {
	return link(href, href)
}

func header(title string, args ...string) Node {
	r := Div(NewAttr())
	bc := Div(NewAttr().Add("class", "bc"), link("/", "root"))
	u := ""
	for i := 0; i < len(args); i++ {
		u += "/" + args[i]
		bc.Append(link(u, args[i]))
	}
	r.Append(
		H1(NewAttr(), Text(title)),
		Div(NewAttr(), bc))

	return r
}

func listTopics(app *echo.Echo, store Store, cont chan string) {
	store.Register("mail/topic-list",
		`SELECT DISTINCT(topic) topic, count(id) FROM {{.RawMails}} GROUP BY topic;`)

	app.GET("/", func(c echo.Context) error {
		var doc = NewDoc()
		q := store.QueryFunc("mail/topic-list")
		var (
			topic string
			count int
		)

		doc.body.Append(header("Topics"))
		attrs := NewAttr().Add("class", "item")

		_, err := q(RowCallback(func() {
			doc.body.Append(Div(attrs,
				A(NewAttr().Add("href", "/"+topic),
					Text(topic)),
				Textf(" (%v)", count)))
		}), &topic, &count)

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.HTML(http.StatusOK, doc.Render())
	})
}

func listInTopics(app *echo.Echo, store Store, cont chan string) {
	store.Register("mail/topic-messages",
		`SELECT id, sender, subject FROM {{.RawMails}} WHERE topic = $1;`)

	handler := func(c echo.Context) error {
		paramTopic := c.Param("topic")
		q := store.QueryFunc("mail/topic-messages", paramTopic)
		var doc = NewDoc()
		var (
			id      int
			sender  string
			subject string
		)
		doc.body.Append(header(paramTopic, paramTopic))
		attrs := NewAttr().Add("class", "item")

		_, err := q(RowCallback(func() {
			url := fmt.Sprintf("/%s/%v", paramTopic, id)
			doc.body.Append(Div(attrs,
				A(NewAttr().Add("href", url),
					Text(subject)),
				Textf(" by %s", sender)))
		}), &id, &sender, &subject)

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.HTML(http.StatusOK, doc.Render())
	}

	app.GET("/:topic", handler)
}

func showMessage(app *echo.Echo, store Store, cont chan string) {
	store.Register("mail/get-message",
		`SELECT id, sender, message FROM {{.RawMails}} WHERE id = $1;`)

	handler := func(c echo.Context) error {
		paramId := c.Param("id")
		paramTopic := c.Param("topic")
		q := store.QueryFunc("mail/get-message", paramId)
		var doc = NewDoc()
		var (
			id      int
			sender  string
			message string
		)

		_, err := q(RowCallback(func() {
			subject, msgNode := formatMessage(id, message)
			doc.body.Append(
				header(subject, paramTopic, paramId),
				msgNode)
		}), &id, &sender, &message)

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.HTML(http.StatusOK, doc.Render())
	}

	app.GET("/:topic/:id", handler)
}

func showAttachment(app *echo.Echo, store Store, cont chan string) {

	store.Register("mail/get-attachment",
		`SELECT id,  message FROM {{.RawMails}} WHERE id = $1;`)
	handler := func(c echo.Context) error {
		id := c.Param("id")
		name := c.Param("name")
		rows, err := store.Query("mail/get-attachment", id)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		for rows.Next() {
			var (
				id      int
				message string
			)
			scanError := rows.Scan(&id, &message)
			if scanError != nil {
				errMesg := scanError.Error()
				cont <- errMesg
			} else {
				a := getAttachment(name, message)
				return c.Stream(http.StatusOK, a.mediaType, a.r)
			}
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
