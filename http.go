package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"net/mail"
	"strconv"
	"strings"

	"github.com/labstack/echo"
)

func mainContentType(p *multipart.Part) string {
	ct := p.Header.Get("Content-Type")
	ctParts := strings.Split(ct, "/")
	return ctParts[0]
}

type SerializedPart interface {
	ContentType() string
	Content() string
	FileName() string
	MainType() string
	SubType() string
}

type serPart struct {
	contentType string
	content     string
	filename    string
}

func (s serPart) ContentType() string {
	return s.contentType
}
func (s serPart) Content() string {
	return s.content
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
	content, err := ioutil.ReadAll(p)
	if err != nil {
		return nil
	}
	return serPart{
		contentType: p.Header.Get("Content-Type"),
		content:     fmt.Sprintf("%s", content),
		filename:    p.FileName(),
	}
}

func walkParts(r io.Reader, boundary string, c chan SerializedPart, depth int) {
	reader := multipart.NewReader(r, boundary)
	counter := 0
	for {
		counter = counter + 1
		part, err := reader.NextPart()
		if err != nil {
			break
		}

		contentType := part.Header.Get("Content-Type")
		mediaType, params, err := mime.ParseMediaType(contentType)
		if err != nil {
			break
		}

		log.Printf("mediatype  %v %v %v", depth, counter, mediaType)
		if !strings.HasPrefix(mediaType, "multipart/") {
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

func formatMultipart(msg *mail.Message, params map[string]string) string {
	results := []string{}
	mc := make(chan SerializedPart)
	boundary := params["boundary"]
	go walkParts(msg.Body, boundary, mc, 0)

	for {
		part := <-mc
		// log.Printf("Received Part %v", part)
		if part == nil {
			break
		}
		switch mp := part.MainType(); mp {
		case "text":
			if true {
				results = append(results, part.Content())
			}
		default:
			log.Printf("formatMultipart %s", mp)
			results = append(results,
				fmt.Sprintf("<a href=\"/attachments/%s\">%s</a>", part.FileName(), part.FileName()))
		}

	}

	return strings.Join(results, "\n")
}

func formatMessage(data string) string {
	msg, err := mail.ReadMessage(strings.NewReader(data))
	if err != nil {
		return "<h2>Error</h2><pre>" + err.Error() + "</pre>"
	}

	mediatype, params, err := mime.ParseMediaType(msg.Header.Get("Content-Type"))
	if err != nil {
		if "" == mediatype {
			body, err := ioutil.ReadAll(msg.Body)
			if err != nil {
				return "<h2>Error</h2><pre>" + err.Error() + "</pre>"
			}

			return fmt.Sprintf("%s", body)
		}
		return "<h2>Error</h2><pre>" + err.Error() + "</pre>"
	}

	return formatMultipart(msg, params)
}

func regHTTPHandlers(app *echo.Echo, store Store, cont chan string) {

	listTopics := func(c echo.Context) error {
		rows, err := store.Query("mail/topic-list")
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		var results = []string{"<h1>Topics</h1>"}
		for rows.Next() {
			var (
				topic string
				count int
			)
			scanError := rows.Scan(&topic, &count)
			if scanError != nil {
				errMesg := scanError.Error()
				cont <- errMesg
			} else {
				results = append(results, "<div><a href=\"/"+topic+"\">"+topic+"</a> ("+strconv.Itoa(count)+") </div>")
			}
		}
		return c.HTML(http.StatusOK, strings.Join(results, "\n"))
	}

	store.Register("mail/topic-list",
		`SELECT DISTINCT(topic) topic, count(id) FROM {{.RawMails}} GROUP BY topic;`)
	app.GET("/", listTopics)

	listInTopic := func(c echo.Context) error {
		paramTopic := c.Param("topic")
		rows, err := store.Query("mail/topic-messages", paramTopic)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		var results = []string{"<h1>" + paramTopic + "</h1>"}
		for rows.Next() {
			var (
				id      int
				sender  string
				subject string
			)
			scanError := rows.Scan(&id, &sender, &subject)
			if scanError != nil {
				errMesg := scanError.Error()
				cont <- errMesg
			} else {
				results = append(results, "<div><a href=\"/"+paramTopic+"/"+strconv.Itoa(id)+"\">"+subject+"</a></div>")
			}
		}
		return c.HTML(http.StatusOK, strings.Join(results, "\n"))
	}

	store.Register("mail/topic-messages",
		`SELECT id, sender, subject FROM {{.RawMails}} WHERE topic = $1;`)
	app.GET("/:topic", listInTopic)

	showMessage := func(c echo.Context) error {
		paramTopic := c.Param("topic")
		id := c.Param("id")
		rows, err := store.Query("mail/get-message", id, paramTopic)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		var results = []string{}
		for rows.Next() {
			var (
				id      int
				sender  string
				message string
			)
			scanError := rows.Scan(&id, &sender, &message)
			if scanError != nil {
				errMesg := scanError.Error()
				cont <- errMesg
			} else {
				results = append(results, formatMessage(message))
			}
		}
		return c.HTML(http.StatusOK, strings.Join(results, "\n"))
	}

	store.Register("mail/get-message",
		`SELECT id, sender, message FROM {{.RawMails}} WHERE id = $1 AND topic = $2;`)
	app.GET("/:topic/:id", showMessage)
}

func StartHTTP(cont chan string, store Store) {
	app := echo.New()
	regHTTPHandlers(app, store, cont)
	cont <- "Configured HTTP Server"
	app.Start(":8080")
}
