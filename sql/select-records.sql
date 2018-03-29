SELECT  
    id, ts, sender, topic, header_subject, body, parent 
FROM {{.Records}} 
WHERE 
    domain = $1
    id = $2