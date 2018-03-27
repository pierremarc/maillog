package main

import "log"

func MakeMigration(store Store, v Volume) {

	store.Register("migrate/get-message",
		`SELECT  id, sender, topic, message 
		FROM {{.RawMails}} `)

	store.Register("mail/record",
		`INSERT INTO {{.Records}} 
		(sender, recipient, topic, domain, header_date, header_subject, body, payload)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id`)

	store.Register("attachment/record",
		`INSERT INTO {{.Attachments}}
		(record_id, content_type, file_name)
		VALUES ($1, $2, $3)`)

	// store.Register("migrate/get-answers",
	// 	`SELECT r.id as id, r.sender as sender, r.message as message
	// 	FROM {{.RawMails}} r
	// 		LEFT Join {{.Answers}} a ON r.id = a.child
	// 	WHERE a.parent = $1;`)

	q := store.QueryFunc("migrate/get-message")
	var (
		id      int
		sender  string
		topic   string
		message []byte
	)
	q(RowCallback(func() {
		migrateMessage(store, v, topic, sender, message).FoldF(
			func(err error) { log.Printf("Error %s", err.Error()) },
			func(i int) { log.Printf("Success %d", i) })
	}), &id, &sender, &topic, &message)
}

func migrateMessage(store Store, v Volume, recipient string, from string, data []byte) ResultInt {
	var (
		id            int
		sender        string
		topic         string
		domain        string
		dateHeader    string
		subjectHeader string
		body          string
	)

	sender = from
	domain = "no-domain"
	topic = getTopic(recipient).FoldString("na", IdString)

	return MakeSerializedMsg(&data).
		MapInt(func(sm SerializedMessage) int {
			var attachments []SerializedPart
			dateHeader = sm.Get("Date")
			subjectHeader = sm.Get("Subject")

			sm.Parse().
				Root().
				Walk(func(p SerializedPart) {
					if p != nil {
						if "text/plain" == p.ContentType() {
							body += p.ContentString()
						} else {
							attachments = append(attachments, p)
						}
					}
				})

			store.QueryFunc("mail/record",
				sender, recipient, topic, domain,
				dateHeader, subjectHeader, body, data).
				Exec(&id).
				FoldF(
					func(err error) { log.Printf("Error:mail/record %s", err.Error()) },
					func(r bool) { log.Printf("Recorded %d", id) })

			for _, at := range attachments {
				v.Write(WriteOptions{encodedSender(sender), topic, id, at.FileName(), at.DecodedContent()}).
					Map(func(fn string) {
						store.QueryFunc("attachment/record",
							id, at.ContentType(), fn).Exec()
					})
			}

			return id
		})
}
