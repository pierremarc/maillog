SELECT  
    id, ts, sender, topic, header_subject, body, parent 
FROM {{.Records}} 
WHERE 
    id = $1