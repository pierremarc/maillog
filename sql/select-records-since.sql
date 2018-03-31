SELECT  
    id, ts, sender, topic, header_subject, body 
FROM {{.Records}} 
WHERE 
    domain = $1
    AND ts >= $2::date