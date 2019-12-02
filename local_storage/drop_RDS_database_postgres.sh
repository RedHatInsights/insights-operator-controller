#!/bin/sh

SUPERUSER= #<master username>
SU_PASSWORD= #<master password> 

DATABASE=controller
DB_SERVER= # <DB instance endpoint>

USER=tester
USER_PASSWORD=tester


# drop database 
psql "postgresql://${SUPERUSER}:${SU_PASSWORD}@${DB_SERVER}/postgres" -c "DROP DATABASE IF EXISTS ${DATABASE};" 
psql  "postgresql://${SUPERUSER}:${SU_PASSWORD}@${DB_SERVER}/postgres" -c "REVOKE ALL PRIVILEGES ON SCHEMA public FROM ${USER};DROP ROLE ${USER};" 

