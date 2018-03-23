package main

import (
	"bytes"
	"fmt"
	"net"
	"net/mail"
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

func getRecipent(to []string) (string, error) {
	if len(to) > 0 {
		addr := to[0]
		parts := strings.Split(addr, "@")
		return parts[0], nil
	}
	return "", makeError("Empty recipients list")
}

func withAnswer(topic string) (string, string) {
	return "", nil
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

		recipient, err := getRecipent(to)
		if err != nil {
			cont <- err.Error()
			return
		}

		msg, _ := mail.ReadMessage(bytes.NewReader(data))
		subject := msg.Header.Get("Subject")
		cont <- fmt.Sprintf("Received mail from %s for %s with subject %s", from, recipient, subject)

		rows, err = store.Query("mail/log", from, recipient, subject, data)
		if err != nil {
			cont <- err.Error()
		}
	}
}

func StartSMTP(cont chan string, iface string, store Store) {
	handler := makeHandler(cont, store)
	cont <- fmt.Sprintf("SMTPD ready on %s", iface)
	smtpd.ListenAndServe(iface, handler, "MailLog", "Wow")
}
