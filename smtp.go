package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/mhale/smtpd"
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

func getRecipent(to []string) OptionString {
	if len(to) > 0 {
		addr := to[0]
		parts := strings.Split(addr, "@")
		return SomeString(parts[0])
	}
	return NoneString()
}

func getDomain(to []string) OptionString {
	if len(to) > 0 {
		addr := to[0]
		parts := strings.Split(addr, "@")
		if len(parts) > 1 {
			return SomeString(parts[1])
		}
	}
	return NoneString()
}

func getTopic(recipient string) OptionString {
	parts := strings.Split(recipient, "+")
	if len(parts) > 0 {
		return SomeString(parts[0])
	}
	return NoneString()
}

func getAnswer(topic string) OptionUInt64 {
	parts := strings.Split(topic, "+")
	if len(parts) > 1 {
		i, err := strconv.ParseUint(parts[1], 10, 32)
		if err == nil {
			return SomeUInt64(i)
		}
	}
	return NoneUInt64()
}

func makeHandler(cont chan string, store Store, v Volume) smtpd.Handler {

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
			domain = getDomain(to).FoldString("", IdString)
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
						v.Write(WriteOptions{encodedSender(sender), topic, id, at.FileName(), at.DecodedContent()}).
							Map(func(fn string) {
								store.QueryFunc(QueryInsertAttachment,
									id, at.ContentType(), fn).Exec()
							})
					}

					return fmt.Sprintf("Received [%s] => [%s]: %s",
						from, recipient, subjectHeader)
				}).
				FoldString("Failed Processing Message", IdString)
		}

		cont <- getRecipent(to).FoldStringF(fail, success)
	}
}

func StartSMTP(cont chan string, iface string, store Store, v Volume) {
	handler := makeHandler(cont, store, v)
	cont <- fmt.Sprintf("SMTPD ready on %s", iface)
	smtpd.ListenAndServe(iface, handler, "MailLog", "Wow")
}
