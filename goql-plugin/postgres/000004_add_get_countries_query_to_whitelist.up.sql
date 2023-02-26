WITH new_query AS (
   INSERT INTO goql.queries (name, created_at) VALUES ('getCountries', current_date)
   RETURNING id
)
INSERT INTO goql.whitelists (query_id, created_at)
VALUES
((select id from new_query), current_date);