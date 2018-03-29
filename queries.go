package main


const QueryInsertAttachment = "InsertAttachment"

const QueryInsertRecord = "InsertRecord"

const QueryInsertRecordParent = "InsertRecordParent"

const QuerySelectAnswerids = "SelectAnswerids"

const QuerySelectAnswers = "SelectAnswers"

const QuerySelectAttachment = "SelectAttachment"

const QuerySelectAttachments = "SelectAttachments"

const QuerySelectDomainMx = "SelectDomainMx"

const QuerySelectDomains = "SelectDomains"

const QuerySelectIntopic = "SelectIntopic"

const QuerySelectRecord = "SelectRecord"

const QuerySelectRecords = "SelectRecords"

const QuerySelectTopics = "SelectTopics"


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
	
	store.Register(QuerySelectRecords, `SELECT  
    id, ts, sender, topic, header_subject, body, parent 
FROM {{.Records}} 
WHERE 
    domain = $1
    id = $2`)
	
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
	
}