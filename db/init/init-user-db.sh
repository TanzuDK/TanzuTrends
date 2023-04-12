#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
	CREATE TABLE trends(
		id serial PRIMARY KEY,
		create_dtm VARCHAR,
		order_id VARCHAR,
		phone VARCHAR,
		name VARCHAR,
		address VARCHAR,
		menu VARCHAR,
		total_item VARCHAR,
		pay VARCHAR

	)
EOSQL