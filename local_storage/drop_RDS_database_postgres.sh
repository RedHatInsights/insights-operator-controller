#!/bin/sh
if [[ -z "$RDS_MASTERUSER" ]]; then
  echo '$RDS_MASTERUSER not set !'
  exit 1
else
  SUPERUSER= $RDS_MASTERUSER #<master username>
fi

if [[ -z "$RDS_MASTERPASSWORD" ]]; then
  echo '$RDS_MASTEPASSWORD not set !'
  exit 1
else
  SU_PASSWORD= $RDS_MASTERPASSWORD #<master password> 
fi
if [[ -z "$RDS_ENDPOINT " ]]; then
  echo '$RDS_ENDPOINT not set !'
  exit 1
else
  DB_SERVER= $RDS_ENDPOINT # <DB_instance_endpoint:port>
fi


DATABASE=controller
USER=tester
USER_PASSWORD=tester


# drop database 
psql "postgresql://${SUPERUSER}:${SU_PASSWORD}@${DB_SERVER}/postgres" -c "DROP DATABASE IF EXISTS ${DATABASE};" 
psql  "postgresql://${SUPERUSER}:${SU_PASSWORD}@${DB_SERVER}/postgres" -c "REVOKE ALL PRIVILEGES ON SCHEMA public FROM ${USER};DROP ROLE ${USER};" 

