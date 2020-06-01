/*
Copyright © 2019, 2020 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package storage_test

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/RedHatInsights/insights-operator-controller/storage"
)

const sqlite3 = "sqlite3"
const dataSource = ":memory:"

// MustGetMockStorage creates test sqlite storage in file
func MustGetMockStorage(tb testing.TB, init bool) (storage.Storage, func()) {
	sqliteStorage := mustGetSqliteStorage(tb, dataSource, init)

	return sqliteStorage, func() {
		MustCloseStorage(tb, sqliteStorage)
	}
}

// MustCloseStorage closes the storage and calls t.Fatal on error
func MustCloseStorage(tb testing.TB, s storage.Storage) {
	s.Close()
}

// FailOnError wraps result of function with one argument
func FailOnError(t testing.TB, err error) {
	// assert.NoError is used to show human readable output
	assert.NoError(t, err)
	// assert.NoError doesn't stop next test execution which can cause strange panic because
	// there was error and some object was not constructed
	if err != nil {
		t.Fatal(err)
	}
}

// mustGetSqliteStorage creates a mock DB storage based on SQLite engine
func mustGetSqliteStorage(tb testing.TB, datasource string, init bool) storage.Storage {
	db, err := sql.Open(sqlite3, datasource)
	FailOnError(tb, err)

	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	FailOnError(tb, err)

	sqliteStorage := storage.NewFromConnection(db, sqlite3)

	return sqliteStorage
}

func emtpyDatabaseError(t *testing.T) {
	t.Fatal("Error should be reported on empty database")
}

// TestNewStorage checks whether constructor for new storage returns error for improper storage configuration
func TestNewStorageError(t *testing.T) {
	_, err := storage.New("non existing driver", "data source")
	assert.EqualError(t, err, "sql: unknown driver \"non existing driver\" (forgotten import?)")
}

// TestNewStorage checks whether constructor for new storage does not return error for proper storage configuration
func TestNewStorageNoError(t *testing.T) {
	_, err := storage.New("sqlite3", dataSource)
	if err != nil {
		t.Fatal("Error is not expected", err)
	}
}
