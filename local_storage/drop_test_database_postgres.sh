#!/usr/bin/env bash

SUPERUSER=postgres
SU_PASSWORD=postgres

DATABASE=test_db
DB_SERVER=localhost

USER=tester

# drop database 
psql "postgresql://${SUPERUSER}:${SU_PASSWORD}@${DB_SERVER}" -c "DROP DATABASE IF EXISTS ${DATABASE};" 
psql  "postgresql://${SUPERUSER}:${SU_PASSWORD}@${DB_SERVER}" -c "REVOKE ALL PRIVILEGES ON SCHEMA public FROM ${USER};DROP ROLE ${USER};" 

