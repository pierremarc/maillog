-- insert attachment

INSERT INTO {{.Attachments}}
    (record_id, content_type, file_name)
VALUES 
    ($1, $2, $3)