#!/usr/bin/env bash

if go build
then
    echo "Build ok"
    ./insights-operator-controller --dbdriver=postgres --storage=postgres://postgres:postgres@localhost/controller?sslmode=disable
else
    echo "Build failed"
fi
