-- answers
SELECT  
    id, ts, sender, topic, header_subject, body, parent 
FROM {{.Records}} 
WHERE parent = $1
ORDER BY ts ASC