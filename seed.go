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
	"crypto/md5"
	"fmt"
	"strconv"

	"github.com/jackc/pgx/pgtype"
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

func SeedIndex(store Store, i Index) {
	q := store.QueryFunc(QuerySelectAllRecords)
	var (
		id            int
		ts            pgtype.Timestamptz
		sender        string
		topic         string
		headerSubject string
		body          string
	)
	q(RowCallback(func() {
		if !isSecretTopic(topic) {
			i.Push(strconv.Itoa(id), IndexRecord{headerSubject, body})
		}
	}), &id, &ts, &sender, &topic, &headerSubject, &body)
}
