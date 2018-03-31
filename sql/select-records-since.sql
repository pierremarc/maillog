SELECT  
    id, ts, sender, topic, header_subject, body 
FROM {{.Records}} 
WHERE 
    domain = $1
    AND topic = $2
    AND ts >= $3::date