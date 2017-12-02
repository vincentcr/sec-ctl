#!/bin/bash
set -e

DB_PASSWORD=md5$(cat /pgsql-bootstrap/md5_passwd)
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" -f /pgsql-bootstrap/db.sql -v DB_NAME=$DB_NAME  -v DB_PASSWORD=$DB_PASSWORD  -v DB_USER=$DB_USER
psql -v ON_ERROR_STOP=1 --username "$DB_USER" -d $DB_NAME  -f /pgsql-bootstrap/tables.sql
psql -v ON_ERROR_STOP=1 --username "$DB_USER" -d $DB_NAME  -f /pgsql-bootstrap/functions.sql

