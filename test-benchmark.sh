#!/usr/bin/env bash
# Copyright 2020 Red Hat, Inc
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


export DBDRIVER=$1

echo "$DBDRIVER"

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
