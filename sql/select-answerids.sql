-- answer ids
SELECT id, sender
FROM {{.Records}}
WHERE parent = $1
ORDER BY ts ASC