go build

if [ $? -eq 0 ]
then
    echo "Build ok"
    ./insights-operator-controller --dbdriver=postgres --storage=postgres://postgres:postgres@localhost/controller?sslmode=disable
else
    echo "Build failed"
fi
