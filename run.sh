go build

if [ $? -eq 0 ]
then
    echo "Build ok"
    ./insights-operator-controller
else
    echo "Build failed"
fi
