#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
	CREATE USER iot;
	CREATE DATABASE iot;
	GRANT ALL PRIVILEGES ON DATABASE iot TO iot;
EOSQL
