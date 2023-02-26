DELETE FROM goql.whitelists where query_id=(SELECT id FROM goql.queries where name='getCountries');

DELETE FROM goql.queris where name='getCountries';