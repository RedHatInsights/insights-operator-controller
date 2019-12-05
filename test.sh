
if [ -z "$CONTROLLER_ENV" ]; then
    export CONTROLLER_ENV=test
fi

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
./local_storage/create_test_database_sqlite.sh
echo "Done"

echo "Starting service"
./insights-operator-controller --storage test.db &
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
