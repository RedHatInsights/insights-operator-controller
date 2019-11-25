#!/bin/sh

SUPERUSER=samuelRHadmin
SU_PASSWORD=v3ry53cur3

DATABASE=controller_test_db
DB_SERVER=rhtestinstance.cux1erificun.us-east-1.rds.amazonaws.com

USER=tester
USER_PASSWORD=tester


# drop database 
psql "postgresql://${SUPERUSER}:${SU_PASSWORD}@${DB_SERVER}/postgres" -c "DROP DATABASE IF EXISTS ${DATABASE};" 
psql  "postgresql://${SUPERUSER}:${SU_PASSWORD}@${DB_SERVER}/postgres" -c "REVOKE ALL PRIVILEGES ON SCHEMA public FROM ${USER};DROP ROLE ${USER};" 

