/*
 *  Copyright (C) 2018 Pierre Marchand <pierre.m@atelier-cartographique.be>
 *
 *  This program is free software: you can redistribute it and/or modify
 *  it under the terms of the GNU Affero General Public License as published by
 *  the Free Software Foundation, version 3 of the License.
 *
 *  This program is distributed in the hope that it will be useful,
 *  but WITHOUT ANY WARRANTY; without even the implied warranty of
 *  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *  GNU General Public License for more details.
 *
 *  You should have received a copy of the GNU General Public License
 *  along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package main

import (
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/pgtype"
	"github.com/labstack/echo"
)

func link(href string, label string) Node {
	return A(ClassAttr("link").Set("href", href), Text(label))
}

func link0(href string) Node {
	return link(href, href)
}

func rootHeader(c echo.Context) Node {
	return Div(ClassAttr("header"), H1(ClassAttr("title"), Text(getHostDomain(c))))
}

func header(c echo.Context, title string, args ...string) Node {
	r := Div(ClassAttr("header"))
	bc := Div(ClassAttr("bc"), link("/", "/"+getHostDomain(c)))
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

func makeDocument(page string, hn ...Node) document {
	doc := NewDoc(NewAttr().Set("data-page", page))
	doc.head.Append(HeadMeta(NewAttr().Set("charset", "utf-8")))
	doc.head.Append(HeadMeta(NewAttr().
		Set("name", "viewport").
		Set("content", "width=device-width, initial-scale=1.0")))
	doc.head.Append(Style(NewAttr(), Text(CssReset)))
	doc.head.Append(Style(NewAttr(), Text(CssStyle)))
	doc.head.Append(Script(NewAttr(), Text(JsBundle)))
	for _, n := range hn {
		doc.head.Append(n)
	}
	return doc
}

func listTopics(app *echo.Echo, store Store, v Volume, cont chan string) {

	app.GET("/", func(c echo.Context) error {
		log.Printf("List Topics(%s)", getHostDomain(c))
		var doc = makeDocument("root")
		q := store.QueryFunc(QuerySelectTopics, getHostDomain(c))
		var (
			topic string
			count int
			mts   pgtype.Timestamp
		)

		doc.body.Append(rootHeader(c))
		attrs := ClassAttr("topic")

		return q(RowCallback(func() {
			doc.body.Append(Div(attrs,
				A(ClassAttr("topic-link link").Set("href", "/"+topic),
					Text(topic)),
				Span(ClassAttr("topic-count"), Textf("(%d),", count)),
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
		qf := store.QueryFunc(QuerySelectFirstRecord,
			getHostDomain(c), paramTopic)
		rssLink := HeadLink(NewAttr().
			Set("rel", "alternate").
			Set("type", "application/rss+xml").
			Set("href", fmt.Sprintf("https://%s/.rss/%s", getHostDomain(c), paramTopic)).
			Set("title", fmt.Sprintf("News at %s in %s", getHostDomain(c), paramTopic)))
		var doc = makeDocument("thread", rssLink)
		var (
			id      int
			ts      pgtype.Timestamptz
			sender  string
			subject string
		)

		qf.Exec(&id, &ts, &sender, &subject)

		doc.body.Append(header(c, paramTopic, paramTopic))

		replyto := fmt.Sprintf("mailto:%s@%s", paramTopic, getHostDomain(c))
		rss := fmt.Sprintf("https://%s/.rss/%s", getHostDomain(c), paramTopic)
		firstUrl := fmt.Sprintf("/%s/%d", paramTopic, id)
		intro := P(ClassAttr("topic-replyto-description"),
			Textf(`
			This thread started with a message from %s on %s.
			The subject of it was `,
				senderName(sender), formatTimeDate(ts.Time)),
			A(ClassAttr("first-subject").Set("href", firstUrl),
				Text(decodeSubject(subject))),
			Text(". Send a message to this thread at the following address: "),
			A(ClassAttr("link").Set("href", replyto),
				Textf("%s@%s", paramTopic, getHostDomain(c))),
			Text("To keep track of this thread,you can subscribe to the "),
			A(ClassAttr("rss-link link").Set("href", rss), Text("rss feed")))

		doc.body.Append(
			Div(ClassAttr("topic-replyto-block").Set("data-topic", paramTopic), intro))

		return q(RowCallback(func() {
			url := fmt.Sprintf("/%s/%d", paramTopic, id)

			doc.body.Append(Div(ClassAttr("message-item"),
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

func replyBlock(c echo.Context, topic string, id int, subject string) Node {
	replySubject := fmt.Sprintf("Re: %s/%s", topic, decodeSubject(subject))
	replyto := fmt.Sprintf("mailto:%s+%d@%s?subject=%s",
		topic, id, getHostDomain(c), replySubject)

	return Div(ClassAttr("answer-link"),
		A(ClassAttr("link").Set("href", replyto), Text("Add Section")))
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
		hasAnswer   = false
	)

	qi := store.QueryFunc(QuerySelectAnswerids, pid)
	qm := store.QueryFunc(QuerySelectAnswers, pid)

	var nas = map[int]Node{}
	var vas = map[int]Node{}
	var sas = map[int]string{}

	qi(RowCallback(func() {
		aids = append(aids, id)
		nas[id] = Div(ClassAttr("attachment-block"))
		vas[id] = Div(ClassAttr("viz-block"))
		sas[id] = sender
	}), &id, &sender)

	root := Div(ClassAttr(fmt.Sprintf("answer depth-%v", depth)))

	return qm(RowCallback(func() {
		hasAnswer = true
		url := fmt.Sprintf("/%s/%v", topic, id)

		block := Div(
			ClassAttr("answer-block").Set("data-record", strconv.Itoa(id)))

		// headerBlock := Div(ClassAttr("message-header"),
		// 	Div(ClassAttr("message-sender"), Text(senderName(sender))),
		// 	Div(ClassAttr("message-date"), Text(formatTime(ts.Time))),
		// 	Div(ClassAttr("answer-view"),
		// 		A(ClassAttr("link").Set("href", url), Text("view"))),
		// 	replyBlock(c, topic, id, subject))

		headerBlock := Div(ClassAttr("message-header"),
			Div(ClassAttr("answer-view"),
				A(ClassAttr("section-link").
					Set("href", url).
					Set("title", "View Message"), Text("§"))))

		bodyBlock := Div(ClassAttr("answer-body"), NewRawNode(body), vas[id], nas[id])

		block.Append(headerBlock, bodyBlock)
		root.Append(block, formatAnswers(strconv.Itoa(id), store, c, depth+1))
	}), &id, &ts, &sender, &topic, &subject, &body, &parent).
		FoldNodeF(
			func(err error) Node { return Text(err.Error()) },
			func(_ bool) Node {
				attachments := []attachmentRecord{}
				for _, aid := range aids {
					attBlock := nas[aid]
					vizBlock := vas[aid]
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
						formatAttachments(attachments, attBlock, vizBlock)
						attachments = []attachmentRecord{}
					}), &id, &contentType, &fileName)

				}
				if hasAnswer {
					return root
				}
				return NoDisplay
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

func formatAttachments(rs []attachmentRecord, link Node, viz Node) {
	// node := Div(ClassAttr("attachment-block"))

	for _, r := range rs {
		mt := strings.Split(r.ct, "/")[0]
		url := fmt.Sprintf("/attachments/%s/%s/%d/%s",
			r.sender, r.topic, r.record, r.fn)
		if "image" == mt {
			thumbnailUrl := fmt.Sprintf("/thumbnail/%s/%s/%d/%d/%s",
				r.sender, r.topic, r.record, ThumbnailMedium, r.fn)
			viz.Append(
				A(ClassAttr("attachment image").Set("href", url).Set("title", r.fn),
					Img(ClassAttr("").Set("src", thumbnailUrl))))
		} else {
			if getMediaType(r.ct) != "text/html" { // html version is mostly noise
				log.Printf("link attachment %s", r.ct)
				link.Append(A(ClassAttr("attachment link").Set("href", url), Text(r.fn)))
			}
		}
	}

	// return node
}

func showMessage(app *echo.Echo, store Store, v Volume, cont chan string) {
	handler := func(c echo.Context) error {
		paramId := c.Param("id")
		paramTopic := c.Param("topic")
		qm := store.QueryFunc(QuerySelectRecord, paramId)
		qa := store.QueryFunc(QuerySelectAttachments, paramId)
		var doc = makeDocument("message")
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
		vizBlock := Div(ClassAttr("viz-block"))

		return qm(RowCallback(func() {
			pnode := NoDisplay
			if parent.Status == pgtype.Present {
				pnode = Div(ClassAttr("parent"),
					A(ClassAttr("link").
						Set("href", fmt.Sprintf("/%s/%d", paramTopic, parent.Int)),
						Text("parent")))
			}

			headerBlock := Div(ClassAttr("message-header"),
				Div(ClassAttr("message-sender"), Text(senderName(sender))),
				Div(ClassAttr("message-date"), Text(formatTime(ts.Time))),
				pnode,
				replyBlock(c, paramTopic, id, subject))

			bodyBlock := Div(ClassAttr("message-body"),
				NewRawNode(body), vizBlock, attBlock)

			block.Append(headerBlock, bodyBlock)

			doc.body.Append(
				header(c, ensureSubject(subject), paramTopic, paramId),
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

					formatAttachments(attachments, attBlock, vizBlock)

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

	thumbnailHandler := func(c echo.Context) error {
		sender := c.Param("sender")
		topic := c.Param("topic")
		id := c.Param("id")
		name := c.Param("name")
		sizeParam := c.Param("size")
		size := OptionIntFrom(strconv.Atoi(sizeParam)).FoldInt(ThumbnailSmall, IdInt)

		var (
			ct string
			fn string
		)

		store.QueryFunc(QuerySelectAttachment, id, name).Exec(&ct, &fn)
		fp := path.Join(sender, topic, id, fn)
		tp := GetThumbnail(v.GetPath(fp), uint(size))
		log.Printf("Thumnail %s", tp)

		return v.AbsReader(tp).FoldErrorF(
			func(err error) error {
				return echo.NewHTTPError(http.StatusNotFound, err.Error())
			},
			func(r io.Reader) error {
				return c.Stream(http.StatusOK, ct, r)
			})
	}

	app.GET("/thumbnail/:sender/:topic/:id/:size/:name", thumbnailHandler)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// allow all connections by default
		return true
	},
}

func createWSHandler(w http.ResponseWriter, r *http.Request, h func(*websocket.Conn)) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	h(conn)
}

func notifyHandler(app *echo.Echo, store Store, v Volume, cont chan string, n *Notifier) {
	handler := func(c echo.Context) error {
		w := c.Response()
		r := c.Request()
		createWSHandler(w, r, func(ws *websocket.Conn) {
			defer ws.Close()
			sub := n.Subscribe(func(n Notification) {
				log.Printf("Notification %v", n)
				err := ws.WriteJSON(n)
				if err != nil {
					ws.Close()
				}
			})
			defer n.Unsubscribe(sub)

			ponged := true
			pingDeadline := 12 * time.Second
			ws.SetPongHandler(func(appData string) error {
				ponged = true
				return nil
			})

			ticker := time.NewTicker(pingDeadline)
			defer ticker.Stop()
			go func() {
				for t := range ticker.C {
					if ponged {
						pingData := make([]byte, 0)
						deadline := time.Now().Add(pingDeadline)
						ws.WriteControl(websocket.PingMessage, pingData, deadline)
					} else {
						log.Println("Still waiting for pong", t)
					}
				}
			}()

			for {
				messageType, p, err := ws.ReadMessage()
				if err != nil {
					break
				}
				log.Printf("websocket:%v %s", messageType, p)
			}
		})
		return nil
	}
	app.GET("/.notifications", handler)
}

func lastweek() string {
	y := time.Now().Add(-24 * 7 * time.Hour)
	yy, my, dy := y.Date()
	return fmt.Sprintf("%d-%02d-%02d", yy, my, dy)
}

func rssHandler(app *echo.Echo, store Store, v Volume, cont chan string) {
	handler := func(c echo.Context) error {
		paramTopic := c.Param("topic")
		q := store.QueryFunc(QuerySelectRecordsTopicSince, getHostDomain(c), paramTopic, lastweek())

		var (
			id            int
			ts            pgtype.Timestamptz
			sender        string
			topic         string
			headerSubject string
			body          string
			maxTime       = time.Now().Add(-24 * 7 * time.Hour)
		)

		rss := MakeRSS()
		lastBuidDate := RssBuildDate(NewAttr())
		channel := MakeRssChannel(fmt.Sprintf("%s - %s", getHostDomain(c), paramTopic),
			fmt.Sprintf("https://%s/%s", getHostDomain(c), paramTopic),
			fmt.Sprintf("News from %s in %s", getHostDomain(c), paramTopic),
			fmt.Sprintf("https://%s/.rss/%s", getHostDomain(c), paramTopic))
		channel.Append(lastBuidDate)
		rss.Append(channel)

		q(RowCallback(func() {
			if ts.Time.After(maxTime) {
				maxTime = ts.Time
			}

			url := fmt.Sprintf("https://%s/%s/%d",
				getHostDomain(c), paramTopic, id)

			var (
				rid         int
				contentType string
				fn          string
				ats         = make([]Node, 0)
				thumbnail   Node
			)
			// attachments
			store.QueryFunc(QuerySelectAttachments, id)(RowCallback(func() {
				log.Printf("Got Attachment %v %s", rid, fn)
				mediaType := getMediaType(contentType)
				if getMainType(mediaType) == "image" && thumbnail == nil {
					thumbnail = MakeRssThumbnail(fmt.Sprintf("https://%s/thumbnail/%s/%s/%d/%d/%s",
						getHostDomain(c), encodedSender(sender), paramTopic, rid, ThumbnailSmall, fn))
				}
				mediaURL := fmt.Sprintf("https://%s/attachments/%s/%s/%d/%s",
					getHostDomain(c), encodedSender(sender), paramTopic, rid, fn)
				ats = append(ats, MakeRssMedia(mediaURL, mediaType, fn))
			}), &rid, &contentType, &fn)

			item := MakeRssItem(topic, senderName(sender), decodeSubject(headerSubject), url, body, ts.Time)
			if thumbnail != nil {
				item.Append(thumbnail)
			}
			item.Append(ats...)
			channel.Append(item)
		}), &id, &ts, &sender, &topic, &headerSubject, &body)

		lastBuidDate.Append(Text(maxTime.Format(time.RFC1123Z)))
		c.Response().Header().Set(echo.HeaderContentType, "application/rss+xml; charset=UTF-8")
		return c.String(http.StatusOK, RenderRss(rss))
	}

	handlerAll := func(c echo.Context) error {
		q := store.QueryFunc(QuerySelectRecordsSince, getHostDomain(c), lastweek())

		var (
			id            int
			ts            pgtype.Timestamptz
			sender        string
			topic         string
			headerSubject string
			body          string
			maxTime       = time.Now().Add(-24 * 7 * time.Hour)
		)

		rss := MakeRSS()
		channel := MakeRssChannel(getHostDomain(c),
			fmt.Sprintf("https://%s", getHostDomain(c)),
			fmt.Sprintf("News from %s", getHostDomain(c)),
			fmt.Sprintf("https://%s/.rss", getHostDomain(c)))
		rss.Append(channel)

		q(RowCallback(func() {
			if ts.Time.After(maxTime) {
				maxTime = ts.Time
			}

			url := fmt.Sprintf("https://%s/%s/%d",
				getHostDomain(c), topic, id)
			channel.Append(MakeRssItem(topic, senderName(sender), decodeSubject(headerSubject), url, body, ts.Time))
		}), &id, &ts, &sender, &topic, &headerSubject, &body)

		channel.Append(RssBuildDate(NewAttr(), Text(maxTime.Format(time.RFC1123Z))))
		c.Response().Header().Set(echo.HeaderContentType, "application/rss+xml; charset=UTF-8")
		return c.String(http.StatusOK, RenderRss(rss))
	}

	app.GET("/.rss", handlerAll)
	app.GET("/.rss/:topic", handler)
}

func searchHandler(app *echo.Echo, store Store, v Volume, cont chan string, i Index) {

	app.GET("/.search", func(c echo.Context) error {
		termQuery := c.QueryParam("q")
		domain := getHostDomain(c)
		log.Printf("Search (%s) ", termQuery)
		var (
			id            int
			ts            pgtype.Timestamptz
			sender        string
			topic         string
			headerSubject string
			body          string
			parent        pgtype.Int4
		)
		var doc = makeDocument("root")
		doc.body.Append(H1(ClassAttr("search-title"), Text(termQuery)))

		results := i.Query(termQuery)

		for _, r := range results {
			store.QueryFunc(QuerySelectRecordDomain, domain, r)(RowCallback(func() {
				url := fmt.Sprintf("/%s/%d", topic, id)
				peek := body[:int(math.Min(float64(200), float64(len(body))))]
				elem := Div(ClassAttr("search-result"),
					H2(NewAttr(), A(ClassAttr("link").Set("href", url),
						Textf("%s/%s", topic, headerSubject))),
					Text(peek))

				doc.body.Append(elem)
			}), &id, &ts, &sender, &topic, &headerSubject, &body, &parent)
		}

		return c.HTML(http.StatusOK, doc.Render())
	})
}

func regHTTPHandlers(app *echo.Echo, store Store, v Volume, cont chan string, n *Notifier, i Index) {
	notifyHandler(app, store, v, cont, n)
	rssHandler(app, store, v, cont)
	listTopics(app, store, v, cont)
	listInTopics(app, store, v, cont)
	showMessage(app, store, v, cont)
	showAttachment(app, store, v, cont)
	searchHandler(app, store, v, cont, i)
}

func StartHTTP(cont chan string, iface string, store Store, v Volume, n *Notifier, i Index) {
	app := echo.New()
	regHTTPHandlers(app, store, v, cont, n, i)
	cont <- fmt.Sprintf("HTTP ready on %s", iface)
	app.Start(iface)
}
