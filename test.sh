#!/usr/bin/env bash

# check if file is locked (db is used by another process)
if fuser test.db &> /dev/null; then
	echo Database is locked, please kill or close the process using the file:
	fuser test.db
	exit 1
fi

if [ -z "$CONTROLLER_ENV" ]; then
    export CONTROLLER_ENV=test
fi

if go build -race
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
./local_storage/create_test_database_sqlite.sh
echo "Done"

echo "Starting service"
./insights-operator-controller --storage test.db &

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
