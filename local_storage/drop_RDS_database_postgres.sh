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

