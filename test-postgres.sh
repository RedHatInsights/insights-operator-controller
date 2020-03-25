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


if go-build
then
    echo "Service build ok"
else
    echo "Build failed"
    exit 1
fi

if go build -o rest-api-tests tests/rest_api_tests.go
then
    echo "REST API tests build ok"
else
    echo "Build failed"
    exit 1
fi

echo "Creating test database"
rm -f test.db
./local_storage/drop_test_database_postgres.sh
./local_storage/create_test_database_postgres.sh
echo "Done"

echo "Starting service"
./insights-operator-controller --dbdriver=postgres --storage=postgres://postgres:postgres@localhost/test_db?sslmode=disable &

# shellcheck disable=SC2116
PID=$(echo $!)

# shellcheck disable=SC2181
if [ $? -eq 0 ]
then
    echo "Service started, PID=$PID"
    sleep 2

    echo -e "------------------------------------------------------------------------------------------------"
    ./rest-api-tests
    EXIT_VALUE=$?
    echo -e "------------------------------------------------------------------------------------------------"

    if kill "$PID"
    then
        echo "Service killed, PID=$PID"
    else
        echo "Fatal, can not kill a service process"
    fi
else
    echo "Fatal, can not start service"
fi

exit $EXIT_VALUE
