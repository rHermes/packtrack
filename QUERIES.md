## INSERT INTO QUEUE
./packtrack -nodeid "node-123" -range -tracker "bring" -rangeStart 100000000 -rangeEnd 100500000



## Get types of errors
SELECT jsonb_path_query(resp, '$.consignmentSet[*].error') as ss, count(*) as n FROM scrape_jobs WHERE status = 'success' GROUP BY 1;

## Get who has sent the most packages
SELECT jsonb_path_query(resp, '$.consignmentSet[*].packageSet[*].brand') as brand, count(*) as n FROM scrape_jobs WHERE NOT jsonb_path_exists(resp, '$.consignmentSet[*].error') GROUP BY 1 LIMIT 100;

## Get percentage of jobs done
SELECT count(*) filter (where status <> 'created') / count(*)::numeric as per_done FROM scrape_jobs;

## Number of packges sent to a country
SELECT jsonb_path_query(resp, '$.consignmentSet[*].packageSet[*].recipientAddress.country') as country, count(*) as n FROM scrape_jobs GROUP BY 1 ORDER BY 2 DESC;


## INDEXES FOR SPEED

CREATE INDEX idx_scrape_jobs_error_found ON scrape_jobs (jsonb_path_exists(resp, '$."consignmentSet"[*]."error"'::jsonpath, '{}'::jsonb, false));

## Views for easier work

 SELECT scrape_jobs.id,
    jsonb_path_query(scrape_jobs.resp, '$."consignmentSet"[*]."consignmentId"'::jsonpath) #>> '{}'::text[] AS c,
    jsonb_path_query(scrape_jobs.resp, '$."consignmentSet"[*]'::jsonpath) AS b
   FROM scrape_jobs
  WHERE NOT jsonb_path_exists(scrape_jobs.resp, '$."consignmentSet"[*]."error"'::jsonpath);


 SELECT consignments.c AS con,
    jsonb_path_query(consignments.b, '$."packageSet"[*]."packageNumber"'::jsonpath) #>> '{}'::text[] AS pn,
    jsonb_path_query(consignments.b, '$."packageSet"[*]'::jsonpath) AS pj
   FROM consignments;

 SELECT packages.pn,
    (jsonb_path_query(packages.pj, '$."eventSet"[*]."dateIso"'::jsonpath) #>> '{}'::text[])::timestamp with time zone AS date,
    jsonb_path_query(packages.pj, '$."eventSet"[*]."status"'::jsonpath) #>> '{}'::text[] AS status,
    jsonb_path_query(packages.pj, '$."eventSet"[*]."description"'::jsonpath) #>> '{}'::text[] AS description,
    jsonb_path_query(packages.pj, '$."eventSet"[*]."unitId"'::jsonpath) #>> '{}'::text[] AS unit_id,
    jsonb_path_query(packages.pj, '$."eventSet"[*]."unitType"'::jsonpath) #>> '{}'::text[] AS unit_type,
    jsonb_path_query(packages.pj, '$."eventSet"[*]."countryCode"'::jsonpath) #>> '{}'::text[] AS country_code,
    jsonb_path_query(packages.pj, '$."eventSet"[*]."country"'::jsonpath) #>> '{}'::text[] AS country,
    jsonb_path_query(packages.pj, '$."eventSet"[*]."city"'::jsonpath) #>> '{}'::text[] AS city,
    jsonb_path_query(packages.pj, '$."eventSet"[*]."postalCode"'::jsonpath) #>> '{}'::text[] AS postal_code
   FROM packages;

## Indexes on these views

## First and last events
SELECT pn, first_value(status) OVER w as first_events, last_value(status) OVER w as last_event FROM events WINDOW w as (PARTITION BY pn ORDER BY date ASC)

## Count of first to last events
 SELECT first_event, last_event, count(*) as n FROM (SELECT first_value(status) OVER w as first_event, last_value(status) OVER w as last_event FROM events WINDOW w as (PARTITION BY pn ORDER BY date ASC)) as a GROUP BY 1, 2 ORDER BY 3 DESC;

## Transitions from one event to the next
SELECT prev_status, status, count(*) as n FROM (SELECT pn, date, lag(status) OVER w as prev_status,  status FROM events WINDOW w as (PARTITION BY pn ORDER BY date ASC)) as a GROUP BY 1, 2 ORDER BY 3 DESC;

## First and last countries
SELECT first_country, last_country, count(*) as n FROM (SELECT first_value(country) OVER w as first_country, last_value(country) OVER w as last_country FROM events WINDOW w as (PARTITION BY pn ORDER BY date ASC)) as a GROUP BY 1, 2 ORDER BY 3 DESC;

### Same without nulls

SELECT first_country, last_country, count(*) as n FROM (SELECT first_value(country) OVER w as first_country, last_value(country) OVER w as last_country FROM events WINDOW w as (PARTITION BY pn ORDER BY date ASC)) as a WHERE first_country <> '' AND last_country <> '' GROUP BY 1, 2 ORDER BY 3 DESC;



## Grafana very slow query, but shows how many packages where left at a given time:


Select
  a.time,
  b.n as "numbers.line"
from
  (
    SELECT
      TIMESTAMP WITH TIME ZONE 'epoch' + $__timeGroup(end_time, $myinterval) * INTERVAL '1 second' as time
    FROM
      scrape_jobs
    WHERE
       $__timeFilter(end_time)
    GROUP BY time
    ORDER BY time
  ) as a LEFT JOIN LATERAL (
    SELECT
      count(*) as n
    from
      scrape_jobs sj
    where
      sj.created_at <= a.time
      and
      (sj.end_time IS NULL OR sj.end_time > a.time)
  ) as b ON TRUE


