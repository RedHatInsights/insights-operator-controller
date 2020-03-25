#!/usr/bin/env bash
# Copyright 2020 Red Hat, Inc
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


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

