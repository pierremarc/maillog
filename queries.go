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


const QueryInsertAttachment = "InsertAttachment"

const QueryInsertRecord = "InsertRecord"

const QueryInsertRecordParent = "InsertRecordParent"

const QuerySelectAllPayloads = "SelectAllPayloads"

const QuerySelectAllRecords = "SelectAllRecords"

const QuerySelectAnswerids = "SelectAnswerids"

const QuerySelectAnswers = "SelectAnswers"

const QuerySelectAttachment = "SelectAttachment"

const QuerySelectAttachments = "SelectAttachments"

const QuerySelectDomainMx = "SelectDomainMx"

const QuerySelectDomains = "SelectDomains"

const QuerySelectFirstRecord = "SelectFirstRecord"

const QuerySelectIntopic = "SelectIntopic"

const QuerySelectRecord = "SelectRecord"

const QuerySelectRecordDomain = "SelectRecordDomain"

const QuerySelectRecordsSince = "SelectRecordsSince"

const QuerySelectRecordsTopicSince = "SelectRecordsTopicSince"

const QuerySelectTopics = "SelectTopics"

const QueryTruncateAttachments = "TruncateAttachments"


func RegisterQueries(store Store) {
	
	store.Register(QueryInsertAttachment, `-- insert attachment

INSERT INTO {{.Attachments}}
    (record_id, content_type, file_name)
VALUES 
    ($1, $2, $3)`)
	
	store.Register(QueryInsertRecord, `-- insert record

INSERT INTO {{.Records}} 
(
    sender, 
    recipient, 
    topic, 
    domain, 
    header_date, 
    header_subject, 
    body, 
    payload
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id`)
	
	store.Register(QueryInsertRecordParent, `-- insert record

INSERT INTO {{.Records}} 
(
    sender, 
    recipient, 
    topic, 
    domain, 
    header_date, 
    header_subject, 
    body, 
    payload, 
    parent
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING id`)
	
	store.Register(QuerySelectAllPayloads, `SELECT  
    id, sender, topic, payload
FROM {{.Records}}`)
	
	store.Register(QuerySelectAllRecords, `SELECT  
    id, ts, sender, topic, header_subject, body 
FROM {{.Records}};`)
	
	store.Register(QuerySelectAnswerids, `-- answer ids
SELECT id, sender
FROM {{.Records}}
WHERE parent = $1
ORDER BY ts ASC`)
	
	store.Register(QuerySelectAnswers, `-- answers
SELECT  
    id, ts, sender, topic, header_subject, body, parent 
FROM {{.Records}} 
WHERE parent = $1
ORDER BY ts ASC`)
	
	store.Register(QuerySelectAttachment, `-- attachment (single)
SELECT content_type, file_name
FROM {{.Attachments}}
WHERE 
    record_id = $1 
    AND file_name = $2`)
	
	store.Register(QuerySelectAttachments, `SELECT record_id, content_type, file_name
FROM {{.Attachments}}
WHERE record_id = $1`)
	
	store.Register(QuerySelectDomainMx, `SELECT 
    id
FROM {{.Domains}}
WHERE mx_name = $1`)
	
	store.Register(QuerySelectDomains, `SELECT 
    http_name, mx_name
FROM {{.Domains}}`)
	
	store.Register(QuerySelectFirstRecord, `SELECT  
    id, ts, sender, header_subject 
FROM {{.Records}} 
WHERE
    domain = $1
    AND topic = $2
ORDER BY ts  ASC
LIMIT 1`)
	
	store.Register(QuerySelectIntopic, `SELECT id, ts, sender, header_subject
FROM {{.Records}}  
WHERE 
    domain = $1
    AND topic = $2 
    AND parent IS NULL 
ORDER BY ts DESC`)
	
	store.Register(QuerySelectRecord, `SELECT  
    id, ts, sender, topic, header_subject, body, parent 
FROM {{.Records}} 
WHERE 
    id = $1`)
	
	store.Register(QuerySelectRecordDomain, `SELECT  
    id, ts, sender, topic, header_subject, body, parent 
FROM {{.Records}} 
WHERE 
    domain = $1
    AND id = $2`)
	
	store.Register(QuerySelectRecordsSince, `SELECT  
    id, ts, sender, topic, header_subject, body 
FROM {{.Records}} 
WHERE 
    domain = $1
    AND ts >= $2::date`)
	
	store.Register(QuerySelectRecordsTopicSince, `SELECT  
    id, ts, sender, topic, header_subject, body 
FROM {{.Records}} 
WHERE 
    domain = $1
    AND topic = $2
    AND ts >= $3::date
    AND parent IS NULL `)
	
	store.Register(QuerySelectTopics, `SELECT 
    DISTINCT(topic) 
    topic, 
    count(id) as count, 
    max(ts) as mts
FROM {{.Records}} 
WHERE
    domain = $1
    AND strpos(topic, '_') <> 1
GROUP BY topic
ORDER BY topic ASC`)
	
	store.Register(QueryTruncateAttachments, `-- truncate attachments table (generally before regenerating with sedd option)
TRUNCATE TABLE {{.Attachments}};`)
	
}