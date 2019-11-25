#!/bin/sh


SUPERUSER=samuelRHadmin
SU_PASSWORD=v3ry53cur3

DATABASE=controller_test_db
DB_SERVER=rhtestinstance.cux1erificun.us-east-1.rds.amazonaws.com

USER=tester
USER_PASSWORD=tester


SCRIPT_DIR="$( cd "$( dirname "$0" )" && pwd )"


# create user and database
psql "postgresql://${SUPERUSER}:${SU_PASSWORD}@${DB_SERVER}/postgres" -c "CREATE DATABASE ${DATABASE};" 
psql "postgresql://${SUPERUSER}:${SU_PASSWORD}@${DB_SERVER}/postgres" -c "CREATE USER ${USER} PASSWORD '${USER_PASSWORD}';" 

#create schema 
cat "${SCRIPT_DIR}/schema.sql" | psql  "postgresql://${SUPERUSER}:${SU_PASSWORD}@${DB_SERVER}/${DATABASE}" 

# grant priviliges to user 
psql "postgresql://${SUPERUSER}:${SU_PASSWORD}@${DB_SERVER}/${DATABASE}" -c "GRANT  SELECT, INSERT, UPDATE,  DELETE 
    ON  ALL TABLES IN SCHEMA public
    TO  ${USER}; 
	GRANT  SELECT,  UPDATE
    ON  ALL SEQUENCES IN SCHEMA public 
    TO  ${USER}; 
    GRANT USAGE ON SCHEMA public TO ${USER};"

# insert data as user 
cat "${SCRIPT_DIR}/test_data.sql" | psql "postgresql://${USER}:${USER_PASSWORD}@${DB_SERVER}/${DATABASE}" 



#questions for Pavel :
# 		do i use integer or serial ? 





