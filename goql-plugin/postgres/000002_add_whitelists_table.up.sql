CREATE TABLE goql.whitelists (
   id  serial PRIMARY KEY,
   operation_name VARCHAR(100) UNIQUE NOT NULL,
   created_at timestamp NOT NULL,
   updated_at timestamp
);