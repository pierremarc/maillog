SELECT  
    id, ts, sender, header_subject 
FROM {{.Records}} 
WHERE
    domain = $1
    AND topic = $2
ORDER BY ts  ASC
LIMIT 1