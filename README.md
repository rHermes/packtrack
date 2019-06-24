## INSERT INTO QUEUE
./packtrack -nodeid "node-123" -range -tracker "bring" -rangeStart 100000000 -rangeEnd 100500000



## Get types of errors
SELECT jsonb_path_query(resp, '$.consignmentSet[*].error') as ss, count(*) as n FROM scrape_jobs WHERE status = 'success' GROUP BY 1;

## Get who has sent the most packages
SELECT jsonb_path_query(resp, '$.consignmentSet[*].packageSet[*].brand') as brand, count(*) as n FROM scrape_jobs WHERE NOT jsonb_path_exists(resp, '$.consignmentSet[*].error') GROUP BY 1 LIMIT 100;
