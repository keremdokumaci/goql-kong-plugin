CREATE TABLE goql.queries (
   id  serial PRIMARY KEY,
   name VARCHAR(100) UNIQUE NOT NULL,
   created_at timestamp NOT NULL,
   updated_at timestamp
);