package main

import (
	"bytes"
	"fmt"
	"net"
	"net/mail"
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

func makeHandler(cont chan string, store Store) smtpd.Handler {

	store.Register("mail/log",
		`INSERT INTO {{.RawMails}} (sender, topic, subject, message)
		VALUES ($1, $2, $3, $4)
		RETURNING id`)

	store.Register("mail/answer",
		`INSERT INTO {{.Answers}} (parent, child)
		VALUES ($1, $2)`)

	return func(origin net.Addr, from string, to []string, data []byte) {

		fail := func() string {
			return "Failed to parse a recipient"
		}

		success := func(recipient string) string {
			msg, _ := mail.ReadMessage(bytes.NewReader(data))
			subject := msg.Header.Get("Subject")
			var id int

			getAnswer(recipient).FoldF(
				func() {
					store.QueryFunc("mail/log", from, recipient, subject, data).Exec()
				},
				func(parentId uint64) {
					r := strings.Split(recipient, "+")[0]
					q := store.QueryFunc("mail/log", from, r, subject, data)
					q(RowCallback(func() {
						store.QueryFunc("mail/answer", parentId, id).Exec()
					}), &id)
				})

			return fmt.Sprintf("Received mail from %s for %s with subject %s", from, recipient, subject)
		}

		cont <- getRecipent(to).FoldStringF(fail, success)
	}
}

func StartSMTP(cont chan string, iface string, store Store) {
	handler := makeHandler(cont, store)
	cont <- fmt.Sprintf("SMTPD ready on %s", iface)
	smtpd.ListenAndServe(iface, handler, "MailLog", "Wow")
}
