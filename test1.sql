SELECT COUNT(*),strftime('%Y', A.date)
FROM article A
WHERE A.body LIKE '%2%'
GROUP BY strftime('%Y', A.date) ;