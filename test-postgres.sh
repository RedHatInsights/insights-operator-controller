#!/usr/bin/env bash

go build

if [ $? -eq 0 ]
then
    echo "Service build ok"
else
    echo "Build failed"
    exit 1
fi

go build -o rest-api-tests tests/rest_api_tests.go

if [ $? -eq 0 ]
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
PID=$(echo $!)

if [ $? -eq 0 ]
then
    echo "Service started, PID=$PID"
    sleep 2

    echo -e "------------------------------------------------------------------------------------------------"
    ./rest-api-tests
    EXIT_VALUE=$?
    echo -e "------------------------------------------------------------------------------------------------"

    kill $PID
    if [ $? -eq 0 ]
    then
        echo "Service killed, PID=$PID"
    else
        echo "Fatal, can not kill a service process"
    fi
else
    echo "Fatal, can not start service"
fi

exit $EXIT_VALUE
