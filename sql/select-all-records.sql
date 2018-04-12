SELECT  
    id, ts, sender, topic, header_subject, body 
FROM {{.Records}};