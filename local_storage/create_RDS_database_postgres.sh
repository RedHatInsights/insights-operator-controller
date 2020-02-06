#!/usr/bin/env bash

if [[ -z "$RDS_MASTERUSER" ]]; then
  echo '$RDS_MASTERUSER not set !'
  exit 1
else
  SUPERUSER=$RDS_MASTERUSER
fi

if [[ -z "$RDS_MASTERPASSWORD" ]]; then
  echo '$RDS_MASTEPASSWORD not set !'
  exit 1
else
  SU_PASSWORD=$RDS_MASTERPASSWORD
fi
if [[ -z "$RDS_ENDPOINT " ]]; then
  echo '$RDS_ENDPOINT not set !'
  exit 1
else
  DB_SERVER=$RDS_ENDPOINT
fi


DATABASE=controller


USER=tester
USER_PASSWORD=tester


SCRIPT_DIR="$( cd "$( dirname "$0" )" && pwd )"


# create user and database
psql "postgresql://${SUPERUSER}:${SU_PASSWORD}@${DB_SERVER}/postgres" -c "CREATE DATABASE ${DATABASE};" 
psql "postgresql://${SUPERUSER}:${SU_PASSWORD}@${DB_SERVER}/postgres" -c "CREATE USER ${USER} PASSWORD '${USER_PASSWORD}';" 

#create schema 
cat "${SCRIPT_DIR}/schema_postgres.sql" | psql  "postgresql://${SUPERUSER}:${SU_PASSWORD}@${DB_SERVER}/${DATABASE}" 

# grant priviliges to user 
psql "postgresql://${SUPERUSER}:${SU_PASSWORD}@${DB_SERVER}/${DATABASE}" -c "GRANT  SELECT, INSERT, UPDATE,  DELETE 
    ON  ALL TABLES IN SCHEMA public
    TO  ${USER}; 
	GRANT  SELECT,  UPDATE
    ON  ALL SEQUENCES IN SCHEMA public 
    TO  ${USER}; 
    GRANT USAGE ON SCHEMA public TO ${USER};"

# insert data as user 
cat "${SCRIPT_DIR}/test_data_postgres.sql" | psql "postgresql://${USER}:${USER_PASSWORD}@${DB_SERVER}/${DATABASE}" 

