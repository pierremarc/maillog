-- attachment (single)
SELECT content_type, file_name
FROM {{.Attachments}}
WHERE 
    record_id = $1 
    AND file_name = $2