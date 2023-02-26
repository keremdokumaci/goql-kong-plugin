CREATE TABLE goql.whitelists (
   id  serial PRIMARY KEY NOT NULL,
   query_id int NOT NULL UNIQUE REFERENCES goql.queries(id),
   created_at timestamp NOT NULL,
   updated_at timestamp
);