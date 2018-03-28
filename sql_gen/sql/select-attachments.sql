SELECT record_id, content_type, file_name
FROM {{.Attachments}}
WHERE record_id = $1