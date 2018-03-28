-- topics

SELECT 
    DISTINCT(topic) 
    topic, 
    count(id) as count, 
    max(ts) as mts
FROM {{.Records}} 
WHERE
    domain = $ 1
    AND strpos(topic, '_') <> 1
GROUP BY topic
ORDER BY topic ASC