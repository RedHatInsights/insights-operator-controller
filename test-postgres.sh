#!/usr/bin/env bash

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
