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

// TestDBStorageListOfClustersSchemalessDB check the behaviour of method ListOfClusters on DB without schema
func TestDBStorageListOfClustersSchemalessDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, true)
	defer closer()

	_, err := mockStorage.ListOfClusters()
	if err == nil {
		emtpyDatabaseError(t)
	}
}

// TestDBStorageGetClusterSchemalessDB check the behaviour of method GetCluster on DB without schema
func TestDBStorageGetClusterSchemalessDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, true)
	defer closer()

	_, err := mockStorage.GetCluster(0)
	if err == nil {
		emtpyDatabaseError(t)
	}
}

// TestDBStorageRegisterNewClusterSchemalessDB check the behaviour of method RegisterNewCluster on DB without schema
func TestDBStorageRegisterNewClusterSchemalessDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, true)
	defer closer()

	err := mockStorage.RegisterNewCluster("foobar")
	if err == nil {
		emtpyDatabaseError(t)
	}
}

// TestDBStorageCreateNewClusterSchemalessDB check the behaviour of method CreateNewCluster on DB without schema
func TestDBStorageCreateNewClusterSchemalessDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, true)
	defer closer()

	err := mockStorage.CreateNewCluster(0, "foobar")
	if err == nil {
		emtpyDatabaseError(t)
	}
}

// TestDBStorageDeleteClusterSchemalessDB check the behaviour of method DeleteCluster on DB without schema
func TestDBStorageDeleteClusterSchemalessDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, true)
	defer closer()

	err := mockStorage.DeleteCluster(0)
	if err == nil {
		emtpyDatabaseError(t)
	}
}

// TestDBStorageDeleteClusterByNameSchemalessDB check the behaviour of method DeleteClusterByName on DB without schema
func TestDBStorageDeleteClusterByNameSchemalessDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, true)
	defer closer()

	err := mockStorage.DeleteClusterByName("foobar")
	if err == nil {
		emtpyDatabaseError(t)
	}
}

// TestDBStorageGetClusterByNameSchemalessDB check the behaviour of method GetClusterByName on DB without schema
func TestDBStorageGetClusterByNameSchemalessDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, true)
	defer closer()

	_, err := mockStorage.GetClusterByName("foobar")
	if err == nil {
		emtpyDatabaseError(t)
	}
}

// TestDBStorageListConfigurationProfiles check the behaviour of method ListOfConfigurationProfiles on DB without schema
func TestDBStorageListConfigurationProfilesSchemalessDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, true)
	defer closer()

	_, err := mockStorage.ListConfigurationProfiles()
	if err == nil {
		emtpyDatabaseError(t)
	}
}

// TestDBStorageGetConfigurationProfileSchemalessDB check the behaviour of method GetConfigurationProfile on DB without schema
func TestDBStorageGetConfigurationProfileSchemalessDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, true)
	defer closer()

	_, err := mockStorage.GetConfigurationProfile(0)
	if err == nil {
		emtpyDatabaseError(t)
	}
}

// TestDBStorageStoreConfigurationProfileSchemalessDB check the behaviour of method StoreConfigurationProfile on DB without schema
func TestDBStorageStoreConfigurationProfileSchemalessDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, true)
	defer closer()

	_, err := mockStorage.StoreConfigurationProfile("username0", "description0", "configuration0")
	if err == nil {
		emtpyDatabaseError(t)
	}
}

// TestDBStorageChangeConfigurationProfileSchemalessDB check the behaviour of method ChangeConfigurationProfile on DB without schema
func TestDBStorageChangeConfigurationProfileSchemalessDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, true)
	defer closer()

	_, err := mockStorage.ChangeConfigurationProfile(42, "username0", "description0", "configuration0")
	if err == nil {
		emtpyDatabaseError(t)
	}
}

// TestDBStorageDeleteConfigurationProfileSchemalessDB check the behaviour of method DeleteConfigurationProfile on DB without schema
func TestDBStorageDeleteConfigurationProfileSchemalessDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, true)
	defer closer()

	_, err := mockStorage.DeleteConfigurationProfile(42)
	if err == nil {
		emtpyDatabaseError(t)
	}
}

// TestDBStorageListAllClusterConfigurationsSchemalessDB check the behaviour of method ListAllClusterConfigurations on DB without schema
func TestDBStorageListAllClusterConfigurationsSchemalessDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, true)
	defer closer()

	_, err := mockStorage.ListAllClusterConfigurations()
	if err == nil {
		emtpyDatabaseError(t)
	}
}

// TestDBStorageListClusterConfigurationSchemalessDB check the behaviour of method ListClusterConfiguration on DB without schema
func TestDBStorageListClusterConfigurationSchemalessDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, true)
	defer closer()

	_, err := mockStorage.ListClusterConfiguration("000ffff")
	if err == nil {
		emtpyDatabaseError(t)
	}
}

// TestDBStorageGetClusterConfigurationByIDSchemalessDB check the behaviour of method GetClusterConfigurationByID on DB without schema
func TestDBStorageGetClusterConfigurationByIDSchemalessDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, true)
	defer closer()

	_, err := mockStorage.GetClusterConfigurationByID(0x0001111)
	if err == nil {
		emtpyDatabaseError(t)
	}
}

// TestDBStorageGetClusterActiveConfigurationSchemalessDB check the behaviour of method GetClusterActiveConfiguration on DB without schema
func TestDBStorageGetClusterActiveConfigurationSchemalessDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, true)
	defer closer()

	_, err := mockStorage.GetClusterActiveConfiguration("0x0002222")
	if err == nil {
		emtpyDatabaseError(t)
	}
}

// TestDBStorageGetConfigurationIDForClusterSchemalessDB check the behaviour of method GetConfigurationIDForCluster on DB without schema
func TestDBStorageGetConfigurationIDForClusterSchemalessDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, true)
	defer closer()

	_, err := mockStorage.GetConfigurationIDForCluster("0x0002222")
	if err == nil {
		emtpyDatabaseError(t)
	}
}

// TestDBStorageCreateClusterConfigurationSchemalessDB check the behaviour of method CreateClusterConfiguration on DB without schema
func TestDBStorageCreateClusterConfigurationSchemalessDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, true)
	defer closer()

	_, err := mockStorage.CreateClusterConfiguration("cluster1", "user1", "reason1", "description1", "configuration1")
	if err == nil {
		emtpyDatabaseError(t)
	}
}

// TestDBStorageEnableClusterConfigurationSchemalessDB check the behaviour of method EnableClusterConfiguration on DB without schema
func TestDBStorageEnableClusterConfigurationSchemalessDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, true)
	defer closer()

	_, err := mockStorage.EnableClusterConfiguration("cluster1", "user1", "reason1")
	if err == nil {
		emtpyDatabaseError(t)
	}
}

// TestDBStorageDisableClusterConfigurationSchemalessDB check the behaviour of method DisableClusterConfiguration on DB without schema
func TestDBStorageDisableClusterConfigurationSchemalessDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, true)
	defer closer()

	_, err := mockStorage.DisableClusterConfiguration("cluster1", "user1", "reason1")
	if err == nil {
		emtpyDatabaseError(t)
	}
}

// TestDBStorageEnableOrDisableClusterConfigurationByIDSchemalessDB check the behaviour of method EnableOrDisableClusterConfiguration on DB without schema
func TestDBStorageEnableOrDisableClusterConfigurationByIDSchemalessDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, true)
	defer closer()

	err := mockStorage.EnableOrDisableClusterConfigurationByID(1, "1")
	if err == nil {
		emtpyDatabaseError(t)
	}

	err = mockStorage.EnableOrDisableClusterConfigurationByID(2, "0")
	if err == nil {
		emtpyDatabaseError(t)
	}
}

// TestDBStorageDeleteClusterConfigurationByIDSchemalessDB check the behaviour of method DeleteClusterConfiguration on DB without schema
func TestDBStorageDeleteClusterConfigurationByIDSchemalessDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, true)
	defer closer()

	err := mockStorage.DeleteClusterConfigurationByID(1)
	if err == nil {
		emtpyDatabaseError(t)
	}
}

// TestDBStorageGetTriggerByIDSchemalessDB check the behaviour of method GetTriggerByID on DB without schema
func TestDBStorageGetTriggerByIDSchemalessDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, true)
	defer closer()

	_, err := mockStorage.GetTriggerByID(1)
	if err == nil {
		emtpyDatabaseError(t)
	}
}

// TestDBStorageDeleteTriggerByIDSchemalessDB check the behaviour of method DeleteTriggerByID on DB without schema
func TestDBStorageDeleteTriggerByIDSchemalessDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, true)
	defer closer()

	err := mockStorage.DeleteTriggerByID(1)
	if err == nil {
		emtpyDatabaseError(t)
	}
}

// TestDBStorageChangeOfStateTriggerByIDSchemalessDB check the behaviour of method ChangeStateTrigger on DB without schema
func TestDBStorageChangeStateOfTriggerByIDSchemalessDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, true)
	defer closer()

	err := mockStorage.ChangeStateOfTriggerByID(1, 0)
	if err == nil {
		emtpyDatabaseError(t)
	}

	err = mockStorage.ChangeStateOfTriggerByID(1, 1)
	if err == nil {
		emtpyDatabaseError(t)
	}
}

// TestDBStorageListAllTriggersSchemalessDB check the behaviour of method ListAllTriggers on DB without schema
func TestDBStorageListAllTriggersSchemalessDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, true)
	defer closer()

	_, err := mockStorage.ListAllTriggers()
	if err == nil {
		emtpyDatabaseError(t)
	}
}

// TestDBStorageListClusterTriggersSchemalessDB check the behaviour of method ListClusterTriggers on DB without schema
func TestDBStorageListClusterTriggersSchemalessDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, true)
	defer closer()

	_, err := mockStorage.ListClusterTriggers("clusterX")
	if err == nil {
		emtpyDatabaseError(t)
	}
}

// TestDBStorageListActiveClusterTriggersSchemalessDB check the behaviour of method ListActiveClusterTriggers on DB without schema
func TestDBStorageListActiveClusterTriggersSchemalessDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, true)
	defer closer()

	_, err := mockStorage.ListActiveClusterTriggers("clusterX")
	if err == nil {
		emtpyDatabaseError(t)
	}
}

// TestDBStorageGetTriggerIDSchemalessDB check the behaviour of method GetTriggerID on DB without schema
func TestDBStorageGetTriggerIDSchemalessDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, true)
	defer closer()

	_, err := mockStorage.GetTriggerID("trigger1")
	if err == nil {
		emtpyDatabaseError(t)
	}
}

// TestDBStorageNewTriggerSchemalessDB check the behaviour of method NewTrigger on DB without schema
func TestDBStorageNewTriggerSchemalessDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, true)
	defer closer()

	err := mockStorage.NewTrigger("clusterY", "triggerType1", "user3", "reason3", "link3")
	if err == nil {
		emtpyDatabaseError(t)
	}
}

// TestDBStorageNewTriggerTypeSchemalessDB check the behaviour of method NewTriggerType on DB without schema
func TestDBStorageNewTriggerTypeSchemalessDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, true)
	defer closer()

	err := mockStorage.NewTriggerType("trigger-type-X", "description-of-new-trigger-type")
	if err == nil {
		emtpyDatabaseError(t)
	}
}

// TestDBStorageAckTriggerSchemalessDB check the behaviour of method AckTrigger on DB without schema
func TestDBStorageAckTriggerSchemalessDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, true)
	defer closer()

	err := mockStorage.AckTrigger("cluster-to-ack", 42)
	if err == nil {
		emtpyDatabaseError(t)
	}
}

// TestDBStoragePingSchemalessDB check the behaviour of method Ping on DB without schema
func TestDBStoragePingSchemalessDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, true)
	defer closer()

	err := mockStorage.Ping()
	if err == nil {
		emtpyDatabaseError(t)
	}
}
