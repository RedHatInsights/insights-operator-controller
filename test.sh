go build

if [ $? -eq 0 ]
then
    echo "Build ok"
    echo "Creating test database"
    rm -f test.db
    ./local_storage/create_test_database.sh
    echo "Done"
    echo "Starting service"
    ./insights-operator-controller --storage test.db &
    PID=$(echo $!)
    if [ $? -eq 0 ]
    then
        echo "Service started, PID=$PID"
        go test -count=1 github.com/redhatinsighs/insights-operator-controller/server
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
else
    echo "Build failed"
fi
