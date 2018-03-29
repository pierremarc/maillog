-- insert record

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
RETURNING id