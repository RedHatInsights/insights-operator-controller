#!/bin/sh

docker build -f .local/Dockerfile.postgres -t rhc-postgres -t rhc-postgres .

docker run --name rhc-postgres -p 5432:5432 --rm -d rhc-postgres
