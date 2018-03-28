-- records in topic

SELECT id, ts, sender, header_subject
FROM {{.Records}}  
WHERE 
    domain = $1
    topic = $2 
    AND parent IS NULL 
ORDER BY ts DESC