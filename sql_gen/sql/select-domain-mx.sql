SELECT 
    id
FROM {{.Domains}}
WHERE mx_name = $1