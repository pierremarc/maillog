package main

import (
	"crypto/md5"
	"fmt"
)

func SeedAttachments(store Store, v Volume) {
	store.QueryFunc(QueryTruncateAttachments).Exec()
	q := store.QueryFunc(QuerySelectAllPayloads)
	var (
		id      int
		sender  string
		topic   string
		payload []byte
	)
	q(RowCallback(func() {
		makeAttachment(store, v, sender, topic, id, payload)
	}), &id, &sender, &topic, &payload)
}

func makeAttachment(store Store, v Volume, sender string, topic string, id int, data []byte) {
	MakeSerializedMsg(&data).
		Map(func(sm SerializedMessage) {
			sm.Parse().
				Root().
				Walk(func(p SerializedPart) {
					if p != nil {
						if "text/plain" != p.ContentType() {
							fn := p.FileName()
							if "text/html" == p.MediaType() {
								fn = fmt.Sprintf("%x.html", md5.Sum(p.Content()))
							}
							v.Write(WriteOptions{encodedSender(sender), topic, id, fn, p.DecodedContent()}).
								Map(func(fn string) {
									store.QueryFunc(QueryInsertAttachment,
										id, p.ContentType(), fn).Exec()
								})
						}
					}
				})
		})
}
