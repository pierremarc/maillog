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
	"log"
	"net"
	"time"

	"github.com/pierremarc/smtpd"
)

type ErrorT struct {
	Time   time.Time
	Reason string
}

func (e ErrorT) Error() string {
	return fmt.Sprintf("Error [%v] %s", e.Time, e.Reason)
}

func makeError(reason string) ErrorT {
	return ErrorT{
		Time:   time.Now(),
		Reason: reason,
	}
}

func makeHandler(cont chan string, store Store, v Volume, n *Notifier) smtpd.Handler {

	return func(origin net.Addr, from string, to []string, data []byte) {

		fail := func() string {
			return "Failed to parse a recipient"
		}

		success := func(recipient string) string {
			var (
				id            int
				sender        string
				topic         string
				domain        string
				dateHeader    string
				subjectHeader string
				body          string
				parent        int

				queryName string
			)

			sender = from
			domain = getDomains(to).First().FoldString("", IdString)
			topic = getTopic(recipient).FoldString("na", IdString)

			getAnswer(recipient).FoldF(
				func() {
					log.Println("No Parent")
					queryName = "mail/record"
					parent = -1
				},
				func(i uint64) {
					queryName = "mail/recordp"
					parent = int(i)
					log.Printf("Parent !!! %v", parent)
				})

			return MakeSerializedMsg(&data).
				MapString(func(sm SerializedMessage) string {
					var attachments []SerializedPart
					dateHeader = sm.Get("Date")
					subjectHeader = sm.Get("Subject")

					sm.Parse().
						Root().
						Walk(func(p SerializedPart) {
							if p != nil {
								log.Printf("Part (%s)", p.MediaType())
								if "text/plain" == p.MediaType() {
									body += p.ContentString()
								} else {
									attachments = append(attachments, p)
								}
							}
						})

					qf := store.QueryFunc(QueryInsertRecord,
						sender, recipient, topic, domain,
						dateHeader, subjectHeader, body, data)

					if parent >= 0 {
						qf = store.QueryFunc(QueryInsertRecordParent,
							sender, recipient, topic, domain,
							dateHeader, subjectHeader, body, data, parent)
					}

					qf.Exec(&id).
						FoldF(
							func(err error) { log.Printf("Error:mail/record %s", err.Error()) },
							func(r bool) { log.Printf("Recorded %d", id) })

					for _, at := range attachments {
						fn := at.FileName()
						if "text/html" == at.MediaType() {
							fn = fmt.Sprintf("%d.html", time.Now().Nanosecond())
						}
						v.Write(WriteOptions{encodedSender(sender), topic, id, fn, at.DecodedContent()}).
							Map(func(fn string) {
								store.QueryFunc(QueryInsertAttachment,
									id, at.ContentType(), fn).Exec()
							})
					}

					n.Notify(MakeNotification(topic, id, parent))

					return fmt.Sprintf("Received [%s] => [%s]: %s",
						from, recipient, subjectHeader)
				}).
				FoldString("Failed Processing Message", IdString)
		}

		cont <- getRecipent(to).FoldStringF(fail, success)
	}
}

func makeHandlerRcpt(store Store) smtpd.HandlerRcpt {
	return func(remoteAddr net.Addr, from string, to string) bool {
		var (
			domain = getDomain(to)
			accept = false
			id     int
		)

		q := store.QueryFunc(QuerySelectDomainMx, domain)
		log.Printf("<<RCPT `%s`", domain)
		q(RowCallback(func() {
			log.Printf("GOT>> `%d`", id)
			accept = true
		}), &id)

		return accept
	}
}

func ListenAndServe(addr string, handler smtpd.Handler, rcpt smtpd.HandlerRcpt) error {
	srv := &smtpd.Server{
		Addr:        addr,
		Handler:     handler,
		HandlerRcpt: rcpt,
		MaxSize:     GetMaxSize(),
		Appname:     GetSiteName(),
		Hostname:    GetSiteName(),
	}
	return srv.ListenAndServe()
}

func StartSMTP(cont chan string, iface string, store Store, v Volume, n *Notifier) {
	handler := makeHandler(cont, store, v, n)
	rcpt := makeHandlerRcpt(store)
	cont <- fmt.Sprintf("SMTPD ready on %s", iface)
	ListenAndServe(iface, handler, rcpt)
}
