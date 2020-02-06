#!/usr/bin/env bash

export DBDRIVER=$1

echo $DBDRIVER
if [[ $DBDRIVER != 'sqlite3' ]] && [[ $DBDRIVER != 'postgres' ]]
then
	echo 'usage test-benchmark.sh sqlite3|postgres'
	exit 1
fi

echo 'starting benchmark testing'

if [[ $DBDRIVER == 'sqlite3' ]]
then
	export STORAGE='./../../test.db'

	echo 'creating sqlite db'

	rm ./test.db
	./local_storage/create_test_database_sqlite.sh
	echo 'performing tests, please wait'
	go test -json -bench=. ./tests/benchmark > benchmark_test_out.json
	echo 'benchmark tests completed!'
fi

if [[ $DBDRIVER == 'postgres' ]]
then
	export STORAGE='postgres://postgres:postgres@localhost/test_db?sslmode=disable'
	echo 'creating postgres db'

	./local_storage/drop_test_database_postgres.sh
	./local_storage/create_test_database_postgres.sh

	echo 'performing tests, please wait'
	go test -json -bench=. ./tests/benchmark > benchmark_test_out.json
	echo 'benchmark tests completed!'
fi
