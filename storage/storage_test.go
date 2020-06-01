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

// MustGetMockStorage creates test sqlite storage in file or in memory
func MustGetMockStorage(tb testing.TB, init bool) (storage.Storage, func()) {
	sqliteStorage, _ := mustGetSqliteStorage(tb, dataSource, init)

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

func initializeDatabase(tb testing.TB, connections *sql.DB) {
	statements := []string{
		`
create table cluster (
    ID      integer primary key asc,
    name    text not null
);
		`,
		`
create table configuration_profile (
    ID            integer primary key asc,
    configuration varchar not null,
    changed_at    datetime,
    changed_by    varchar,
    description   varchar
);
		`,
		`
create table operator_configuration (
    ID            integer primary key asc,
    cluster       integer not null,
    configuration integer not null,
    changed_at    datetime,
    changed_by    varchar,
    active        integer,
    reason        varchar,
    CONSTRAINT fk_cluster
        foreign key(cluster)
        references cluster(ID)
        on delete cascade
    CONSTRAINT fk_configuration
        foreign key (configuration)
        references configuration_profile(ID)
        on delete cascade
);
		`,
		`
create table trigger_type (
    ID            integer primary key asc,
    type          varchar not null,
    description   varchar
);
		`,
		`
create table trigger (
    ID            integer primary key asc,
    type          integer not null,
    cluster       integer not null,
    reason        varchar,
    link          varchar,
    triggered_at  datetime,
    triggered_by  varchar,
    acked_at      datetime,
    parameters    varchar,
    active        integer,
    CONSTRAINT fk_type
        foreign key (type)
        references trigger_type(ID)
        on delete cascade
    CONSTRAINT fk_cluster
        foreign key(cluster)
        references cluster(ID)
        on delete cascade
);
		`,
	}
	for _, s := range statements {
		statement, err := connections.Prepare(s)
		FailOnError(tb, err)

		_, err = statement.Exec()
		FailOnError(tb, err)

		err = statement.Close()
		FailOnError(tb, err)
	}
}

// mustGetSqliteStorage creates a mock DB storage based on SQLite engine
func mustGetSqliteStorage(tb testing.TB, datasource string, init bool) (storage.Storage, *sql.DB) {
	db, err := sql.Open(sqlite3, datasource)
	FailOnError(tb, err)

	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	FailOnError(tb, err)

	sqliteStorage := storage.NewFromConnection(db, sqlite3)

	if init {
		initializeDatabase(tb, db)

	}
	return sqliteStorage, db
}

func emptyDatabaseError(t *testing.T) {
	t.Fatal("Error should be reported on schemaless database")
}

func unexpectedDatabaseError(t *testing.T, err error) {
	t.Fatal("Error should not be reported on empty database", err)
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
	mockStorage, closer := MustGetMockStorage(t, false)
	defer closer()

	_, err := mockStorage.ListOfClusters()
	if err == nil {
		emptyDatabaseError(t)
	}
}

// TestDBStorageListOfClustersEmptyDB check the behaviour of method ListOfClusters on empty DB
func TestDBStorageListOfClustersEmptyDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, true)
	defer closer()

	_, err := mockStorage.ListOfClusters()
	if err != nil {
		unexpectedDatabaseError(t, err)
	}
}

// TestDBStorageGetClusterSchemalessDB check the behaviour of method GetCluster on DB without schema
func TestDBStorageGetClusterSchemalessDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, false)
	defer closer()

	_, err := mockStorage.GetCluster(0)
	if err == nil {
		emptyDatabaseError(t)
	}
}

// TestDBStorageGetClusterEmptyDB check the behaviour of method GetCluster on empty DB
func TestDBStorageGetClusterEmptyDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, true)
	defer closer()

	_, err := mockStorage.GetCluster(0)
	if err == nil {
		emptyDatabaseError(t)
	}
}

// TestDBStorageRegisterNewClusterSchemalessDB check the behaviour of method RegisterNewCluster on DB without schema
func TestDBStorageRegisterNewClusterSchemalessDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, false)
	defer closer()

	err := mockStorage.RegisterNewCluster("foobar")
	if err == nil {
		emptyDatabaseError(t)
	}
}

// TestDBStorageRegisterNewClusterEmptyDB check the behaviour of method RegisterNewCluster on empty DB
func TestDBStorageRegisterNewClusterEmptyDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, true)
	defer closer()

	err := mockStorage.RegisterNewCluster("foobar")
	if err != nil {
		unexpectedDatabaseError(t, err)
	}
}

// TestDBStorageCreateNewClusterSchemalessDB check the behaviour of method CreateNewCluster on DB without schema
func TestDBStorageCreateNewClusterSchemalessDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, false)
	defer closer()

	err := mockStorage.CreateNewCluster(0, "foobar")
	if err == nil {
		emptyDatabaseError(t)
	}
}

// TestDBStorageCreateNewClusterEmptyDB check the behaviour of method CreateNewCluster on empty DB
func TestDBStorageCreateNewClusterEmptyDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, true)
	defer closer()

	err := mockStorage.CreateNewCluster(0, "foobar")
	if err != nil {
		unexpectedDatabaseError(t, err)
	}
}

// TestDBStorageDeleteClusterSchemalessDB check the behaviour of method DeleteCluster on DB without schema
func TestDBStorageDeleteClusterSchemalessDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, false)
	defer closer()

	err := mockStorage.DeleteCluster(0)
	if err == nil {
		emptyDatabaseError(t)
	}
}

// TestDBStorageDeleteClusterEmptyDB check the behaviour of method DeleteCluster on empty DB
func TestDBStorageDeleteClusterEmptyDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, true)
	defer closer()

	err := mockStorage.DeleteCluster(0)
	if err == nil {
		emptyDatabaseError(t)
	}
}

// TestDBStorageDeleteClusterByNameSchemalessDB check the behaviour of method DeleteClusterByName on DB without schema
func TestDBStorageDeleteClusterByNameSchemalessDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, false)
	defer closer()

	err := mockStorage.DeleteClusterByName("foobar")
	if err == nil {
		emptyDatabaseError(t)
	}
}

// TestDBStorageDeleteClusterByNameEmptyDB check the behaviour of method DeleteClusterByName on empty DB
func TestDBStorageDeleteClusterByNameEmptyDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, true)
	defer closer()

	err := mockStorage.DeleteClusterByName("foobar")
	if err == nil {
		emptyDatabaseError(t)
	}
}

// TestDBStorageGetClusterByNameSchemalessDB check the behaviour of method GetClusterByName on DB without schema
func TestDBStorageGetClusterByNameSchemalessDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, false)
	defer closer()

	_, err := mockStorage.GetClusterByName("foobar")
	if err == nil {
		emptyDatabaseError(t)
	}
}

// TestDBStorageGetClusterByNameEmptyDB check the behaviour of method GetClusterByName on empty DB
func TestDBStorageGetClusterByNameEmptyDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, true)
	defer closer()

	_, err := mockStorage.GetClusterByName("foobar")
	if err == nil {
		emptyDatabaseError(t)
	}
}

// TestDBStorageListConfigurationProfiles check the behaviour of method ListOfConfigurationProfiles on DB without schema
func TestDBStorageListConfigurationProfilesSchemalessDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, false)
	defer closer()

	_, err := mockStorage.ListConfigurationProfiles()
	if err == nil {
		emptyDatabaseError(t)
	}
}

// TestDBStorageListConfigurationProfiles check the behaviour of method ListOfConfigurationProfiles on empty DB
func TestDBStorageListConfigurationProfilesEmptyDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, true)
	defer closer()

	_, err := mockStorage.ListConfigurationProfiles()
	if err != nil {
		unexpectedDatabaseError(t, err)
	}
}

// TestDBStorageGetConfigurationProfileSchemalessDB check the behaviour of method GetConfigurationProfile on DB without schema
func TestDBStorageGetConfigurationProfileSchemalessDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, false)
	defer closer()

	_, err := mockStorage.GetConfigurationProfile(0)
	if err == nil {
		emptyDatabaseError(t)
	}
}

// TestDBStorageGetConfigurationProfileEmptyDB check the behaviour of method GetConfigurationProfile on empty DB
func TestDBStorageGetConfigurationProfileEmptyDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, true)
	defer closer()

	_, err := mockStorage.GetConfigurationProfile(0)
	if err == nil {
		emptyDatabaseError(t)
	}
}

// TestDBStorageStoreConfigurationProfileSchemalessDB check the behaviour of method StoreConfigurationProfile on DB without schema
func TestDBStorageStoreConfigurationProfileSchemalessDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, false)
	defer closer()

	_, err := mockStorage.StoreConfigurationProfile("username0", "description0", "configuration0")
	if err == nil {
		emptyDatabaseError(t)
	}
}

// TestDBStorageStoreConfigurationProfileEmptyDB check the behaviour of method StoreConfigurationProfile on empty DB
func TestDBStorageStoreConfigurationProfileEmptyDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, true)
	defer closer()

	_, err := mockStorage.StoreConfigurationProfile("username0", "description0", "configuration0")
	if err != nil {
		unexpectedDatabaseError(t, err)
	}
}

// TestDBStorageChangeConfigurationProfileSchemalessDB check the behaviour of method ChangeConfigurationProfile on DB without schema
func TestDBStorageChangeConfigurationProfileSchemalessDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, false)
	defer closer()

	_, err := mockStorage.ChangeConfigurationProfile(42, "username0", "description0", "configuration0")
	if err == nil {
		emptyDatabaseError(t)
	}
}

// TestDBStorageChangeConfigurationProfileEmptyDB check the behaviour of method ChangeConfigurationProfile on empty DB
func TestDBStorageChangeConfigurationProfileEmptyDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, true)
	defer closer()

	_, err := mockStorage.ChangeConfigurationProfile(42, "username0", "description0", "configuration0")
	if err == nil {
		emptyDatabaseError(t)
	}
}

// TestDBStorageDeleteConfigurationProfileSchemalessDB check the behaviour of method DeleteConfigurationProfile on DB without schema
func TestDBStorageDeleteConfigurationProfileSchemalessDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, false)
	defer closer()

	_, err := mockStorage.DeleteConfigurationProfile(42)
	if err == nil {
		emptyDatabaseError(t)
	}
}

// TestDBStorageDeleteConfigurationProfileEmptyDB check the behaviour of method DeleteConfigurationProfile on empty DB
func TestDBStorageDeleteConfigurationProfileEmptyDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, true)
	defer closer()

	_, err := mockStorage.DeleteConfigurationProfile(42)
	if err == nil {
		emptyDatabaseError(t)
	}
}

// TestDBStorageListAllClusterConfigurationsSchemalessDB check the behaviour of method ListAllClusterConfigurations on DB without schema
func TestDBStorageListAllClusterConfigurationsSchemalessDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, false)
	defer closer()

	_, err := mockStorage.ListAllClusterConfigurations()
	if err == nil {
		emptyDatabaseError(t)
	}
}

// TestDBStorageListAllClusterConfigurationsEmptyDB check the behaviour of method ListAllClusterConfigurations on empty DB
func TestDBStorageListAllClusterConfigurationsEmptyDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, true)
	defer closer()

	_, err := mockStorage.ListAllClusterConfigurations()
	if err != nil {
		unexpectedDatabaseError(t, err)
	}
}

// TestDBStorageListClusterConfigurationSchemalessDB check the behaviour of method ListClusterConfiguration on DB without schema
func TestDBStorageListClusterConfigurationSchemalessDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, false)
	defer closer()

	_, err := mockStorage.ListClusterConfiguration("000ffff")
	if err == nil {
		emptyDatabaseError(t)
	}
}

// TestDBStorageListClusterConfigurationEmptyDB check the behaviour of method ListClusterConfiguration on empty DB
func TestDBStorageListClusterConfigurationEmptyDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, true)
	defer closer()

	_, err := mockStorage.ListClusterConfiguration("000ffff")
	if err == nil {
		emptyDatabaseError(t)
	}
}

// TestDBStorageGetClusterConfigurationByIDSchemalessDB check the behaviour of method GetClusterConfigurationByID on DB without schema
func TestDBStorageGetClusterConfigurationByIDSchemalessDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, false)
	defer closer()

	_, err := mockStorage.GetClusterConfigurationByID(0x0001111)
	if err == nil {
		emptyDatabaseError(t)
	}
}

// TestDBStorageGetClusterConfigurationByIDEmptyDB check the behaviour of method GetClusterConfigurationByID on empty DB
func TestDBStorageGetClusterConfigurationByIDEmptyDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, true)
	defer closer()

	_, err := mockStorage.GetClusterConfigurationByID(0x0001111)
	if err == nil {
		emptyDatabaseError(t)
	}
}

// TestDBStorageGetClusterActiveConfigurationSchemalessDB check the behaviour of method GetClusterActiveConfiguration on DB without schema
func TestDBStorageGetClusterActiveConfigurationSchemalessDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, false)
	defer closer()

	_, err := mockStorage.GetClusterActiveConfiguration("0x0002222")
	if err == nil {
		emptyDatabaseError(t)
	}
}

// TestDBStorageGetClusterActiveConfigurationEmptyDB check the behaviour of method GetClusterActiveConfiguration on empty DB
func TestDBStorageGetClusterActiveConfigurationEmptyDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, true)
	defer closer()

	_, err := mockStorage.GetClusterActiveConfiguration("0x0002222")
	if err == nil {
		emptyDatabaseError(t)
	}
}

// TestDBStorageGetConfigurationIDForClusterSchemalessDB check the behaviour of method GetConfigurationIDForCluster on DB without schema
func TestDBStorageGetConfigurationIDForClusterSchemalessDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, false)
	defer closer()

	_, err := mockStorage.GetConfigurationIDForCluster("0x0002222")
	if err == nil {
		emptyDatabaseError(t)
	}
}

// TestDBStorageGetConfigurationIDForClusterEmptyDB check the behaviour of method GetConfigurationIDForCluster on empty DB
func TestDBStorageGetConfigurationIDForClusterEmptyDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, true)
	defer closer()

	_, err := mockStorage.GetConfigurationIDForCluster("0x0002222")
	if err == nil {
		emptyDatabaseError(t)
	}
}

// TestDBStorageCreateClusterConfigurationSchemalessDB check the behaviour of method CreateClusterConfiguration on DB without schema
func TestDBStorageCreateClusterConfigurationSchemalessDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, false)
	defer closer()

	_, err := mockStorage.CreateClusterConfiguration("cluster1", "user1", "reason1", "description1", "configuration1")
	if err == nil {
		emptyDatabaseError(t)
	}
}

// TestDBStorageCreateClusterConfigurationEmptyDB check the behaviour of method CreateClusterConfiguration on empty DB
func TestDBStorageCreateClusterConfigurationEmptyDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, true)
	defer closer()

	_, err := mockStorage.CreateClusterConfiguration("cluster1", "user1", "reason1", "description1", "configuration1")
	if err == nil {
		emptyDatabaseError(t)
	}
}

// TestDBStorageEnableClusterConfigurationSchemalessDB check the behaviour of method EnableClusterConfiguration on DB without schema
func TestDBStorageEnableClusterConfigurationSchemalessDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, false)
	defer closer()

	_, err := mockStorage.EnableClusterConfiguration("cluster1", "user1", "reason1")
	if err == nil {
		emptyDatabaseError(t)
	}
}

// TestDBStorageEnableClusterConfigurationEmptyDB check the behaviour of method EnableClusterConfiguration on empty DB
func TestDBStorageEnableClusterConfigurationEmptyDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, true)
	defer closer()

	_, err := mockStorage.EnableClusterConfiguration("cluster1", "user1", "reason1")
	if err == nil {
		emptyDatabaseError(t)
	}
}

// TestDBStorageDisableClusterConfigurationSchemalessDB check the behaviour of method DisableClusterConfiguration on DB without schema
func TestDBStorageDisableClusterConfigurationSchemalessDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, false)
	defer closer()

	_, err := mockStorage.DisableClusterConfiguration("cluster1", "user1", "reason1")
	if err == nil {
		emptyDatabaseError(t)
	}
}

// TestDBStorageDisableClusterConfigurationEmptyDB check the behaviour of method DisableClusterConfiguration on empty DB
func TestDBStorageDisableClusterConfigurationEmptyDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, true)
	defer closer()

	_, err := mockStorage.DisableClusterConfiguration("cluster1", "user1", "reason1")
	if err == nil {
		emptyDatabaseError(t)
	}
}

// TestDBStorageEnableOrDisableClusterConfigurationByIDSchemalessDB check the behaviour of method EnableOrDisableClusterConfiguration on DB without schema
func TestDBStorageEnableOrDisableClusterConfigurationByIDSchemalessDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, false)
	defer closer()

	err := mockStorage.EnableOrDisableClusterConfigurationByID(1, "1")
	if err == nil {
		emptyDatabaseError(t)
	}

	err = mockStorage.EnableOrDisableClusterConfigurationByID(2, "0")
	if err == nil {
		emptyDatabaseError(t)
	}
}

// TestDBStorageEnableOrDisableClusterConfigurationByIDEmptyDB check the behaviour of method EnableOrDisableClusterConfiguration on empty DB
func TestDBStorageEnableOrDisableClusterConfigurationByIDEmptyDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, true)
	defer closer()

	err := mockStorage.EnableOrDisableClusterConfigurationByID(1, "1")
	if err == nil {
		emptyDatabaseError(t)
	}

	err = mockStorage.EnableOrDisableClusterConfigurationByID(2, "0")
	if err == nil {
		emptyDatabaseError(t)
	}
}

// TestDBStorageDeleteClusterConfigurationByIDSchemalessDB check the behaviour of method DeleteClusterConfiguration on DB without schema
func TestDBStorageDeleteClusterConfigurationByIDSchemalessDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, false)
	defer closer()

	err := mockStorage.DeleteClusterConfigurationByID(1)
	if err == nil {
		emptyDatabaseError(t)
	}
}

// TestDBStorageDeleteClusterConfigurationByIDEmptyDB check the behaviour of method DeleteClusterConfiguration on empty DB
func TestDBStorageDeleteClusterConfigurationByIDEmptyDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, true)
	defer closer()

	err := mockStorage.DeleteClusterConfigurationByID(1)
	if err == nil {
		emptyDatabaseError(t)
	}
}

// TestDBStorageGetTriggerByIDSchemalessDB check the behaviour of method GetTriggerByID on DB without schema
func TestDBStorageGetTriggerByIDSchemalessDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, false)
	defer closer()

	_, err := mockStorage.GetTriggerByID(1)
	if err == nil {
		emptyDatabaseError(t)
	}
}

// TestDBStorageGetTriggerByIDEmptyDB check the behaviour of method GetTriggerByID on empty DB
func TestDBStorageGetTriggerByIDEmptyDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, true)
	defer closer()

	_, err := mockStorage.GetTriggerByID(1)
	if err == nil {
		emptyDatabaseError(t)
	}
}

// TestDBStorageDeleteTriggerByIDSchemalessDB check the behaviour of method DeleteTriggerByID on DB without schema
func TestDBStorageDeleteTriggerByIDSchemalessDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, false)
	defer closer()

	err := mockStorage.DeleteTriggerByID(1)
	if err == nil {
		emptyDatabaseError(t)
	}
}

// TestDBStorageDeleteTriggerByIDEmptyDB check the behaviour of method DeleteTriggerByID on empty DB
func TestDBStorageDeleteTriggerByIDEmptyDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, true)
	defer closer()

	err := mockStorage.DeleteTriggerByID(1)
	if err == nil {
		emptyDatabaseError(t)
	}
}

// TestDBStorageChangeOfStateTriggerByIDSchemalessDB check the behaviour of method ChangeStateTrigger on DB without schema
func TestDBStorageChangeStateOfTriggerByIDSchemalessDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, false)
	defer closer()

	err := mockStorage.ChangeStateOfTriggerByID(1, 0)
	if err == nil {
		emptyDatabaseError(t)
	}

	err = mockStorage.ChangeStateOfTriggerByID(1, 1)
	if err == nil {
		emptyDatabaseError(t)
	}
}

// TestDBStorageChangeOfStateTriggerByIDEmptyDB check the behaviour of method ChangeStateTrigger on empty DB
func TestDBStorageChangeStateOfTriggerByIDEmptyDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, true)
	defer closer()

	err := mockStorage.ChangeStateOfTriggerByID(1, 0)
	if err == nil {
		emptyDatabaseError(t)
	}

	err = mockStorage.ChangeStateOfTriggerByID(1, 1)
	if err == nil {
		emptyDatabaseError(t)
	}
}

// TestDBStorageListAllTriggersSchemalessDB check the behaviour of method ListAllTriggers on DB without schema
func TestDBStorageListAllTriggersSchemalessDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, false)
	defer closer()

	_, err := mockStorage.ListAllTriggers()
	if err == nil {
		emptyDatabaseError(t)
	}
}

// TestDBStorageListAllTriggersEmptyDB check the behaviour of method ListAllTriggers on empty DB
func TestDBStorageListAllTriggersEmptyDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, true)
	defer closer()

	_, err := mockStorage.ListAllTriggers()
	if err != nil {
		unexpectedDatabaseError(t, err)
	}
}

// TestDBStorageListClusterTriggersSchemalessDB check the behaviour of method ListClusterTriggers on DB without schema
func TestDBStorageListClusterTriggersSchemalessDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, false)
	defer closer()

	_, err := mockStorage.ListClusterTriggers("clusterX")
	if err == nil {
		emptyDatabaseError(t)
	}
}

// TestDBStorageListClusterTriggersEmptyDB check the behaviour of method ListClusterTriggers on empty DB
func TestDBStorageListClusterTriggersEmptyDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, true)
	defer closer()

	_, err := mockStorage.ListClusterTriggers("clusterX")
	if err == nil {
		emptyDatabaseError(t)
	}
}

// TestDBStorageListActiveClusterTriggersSchemalessDB check the behaviour of method ListActiveClusterTriggers on DB without schema
func TestDBStorageListActiveClusterTriggersSchemalessDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, false)
	defer closer()

	_, err := mockStorage.ListActiveClusterTriggers("clusterX")
	if err == nil {
		emptyDatabaseError(t)
	}
}

// TestDBStorageListActiveClusterTriggersEmptyDB check the behaviour of method ListActiveClusterTriggers on empty DB
func TestDBStorageListActiveClusterTriggersEmptyDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, true)
	defer closer()

	_, err := mockStorage.ListActiveClusterTriggers("clusterX")
	if err == nil {
		emptyDatabaseError(t)
	}
}

// TestDBStorageGetTriggerIDSchemalessDB check the behaviour of method GetTriggerID on DB without schema
func TestDBStorageGetTriggerIDSchemalessDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, false)
	defer closer()

	_, err := mockStorage.GetTriggerID("trigger1")
	if err == nil {
		emptyDatabaseError(t)
	}
}

// TestDBStorageGetTriggerIDEmptyDB check the behaviour of method GetTriggerID on empty DB
func TestDBStorageGetTriggerIDEmptyDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, true)
	defer closer()

	_, err := mockStorage.GetTriggerID("trigger1")
	if err == nil {
		emptyDatabaseError(t)
	}
}

// TestDBStorageNewTriggerSchemalessDB check the behaviour of method NewTrigger on DB without schema
func TestDBStorageNewTriggerSchemalessDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, false)
	defer closer()

	err := mockStorage.NewTrigger("clusterY", "triggerType1", "user3", "reason3", "link3")
	if err == nil {
		emptyDatabaseError(t)
	}
}

// TestDBStorageNewTriggerEmptyDB check the behaviour of method NewTrigger on empty DB
func TestDBStorageNewTriggerEmptyDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, true)
	defer closer()

	err := mockStorage.NewTrigger("clusterY", "triggerType1", "user3", "reason3", "link3")
	if err == nil {
		emptyDatabaseError(t)
	}
}

// TestDBStorageNewTriggerTypeSchemalessDB check the behaviour of method NewTriggerType on DB without schema
func TestDBStorageNewTriggerTypeSchemalessDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, false)
	defer closer()

	err := mockStorage.NewTriggerType("trigger-type-X", "description-of-new-trigger-type")
	if err == nil {
		emptyDatabaseError(t)
	}
}

// TestDBStorageNewTriggerTypeEmptyDB check the behaviour of method NewTriggerType on empty DB
func TestDBStorageNewTriggerTypeEmptyDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, true)
	defer closer()

	err := mockStorage.NewTriggerType("trigger-type-X", "description-of-new-trigger-type")
	if err != nil {
		unexpectedDatabaseError(t, err)
	}
}

// TestDBStorageAckTriggerSchemalessDB check the behaviour of method AckTrigger on DB without schema
func TestDBStorageAckTriggerSchemalessDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, false)
	defer closer()

	err := mockStorage.AckTrigger("cluster-to-ack", 42)
	if err == nil {
		emptyDatabaseError(t)
	}
}

// TestDBStorageAckTriggerEmptyDB check the behaviour of method AckTrigger on empty DB
func TestDBStorageAckTriggerEmptyDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, true)
	defer closer()

	err := mockStorage.AckTrigger("cluster-to-ack", 42)
	if err == nil {
		emptyDatabaseError(t)
	}
}

// TestDBStoragePingSchemalessDB check the behaviour of method Ping on DB without schema
func TestDBStoragePingSchemalessDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, false)
	defer closer()

	err := mockStorage.Ping()
	if err == nil {
		emptyDatabaseError(t)
	}
}

// TestDBStoragePingEmptyDB check the behaviour of method Ping on empty DB
func TestDBStoragePingEmptyDB(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, true)
	defer closer()

	err := mockStorage.Ping()
	if err != nil {
		unexpectedDatabaseError(t, err)
	}
}

// TestDBStorageClusterOperation1 check the DB operations with clusers
func TestDBStorageClusterOperation1(t *testing.T) {
	const clusterName = "cluster_to_test_1"
	mockStorage, closer := MustGetMockStorage(t, true)
	defer closer()

	err := mockStorage.RegisterNewCluster(clusterName)
	if err != nil {
		unexpectedDatabaseError(t, err)
	}

	_, err = mockStorage.ListOfClusters()
	if err != nil {
		unexpectedDatabaseError(t, err)
	}

	_, err = mockStorage.GetCluster(1)
	if err != nil {
		unexpectedDatabaseError(t, err)
	}

	err = mockStorage.DeleteCluster(1)
	if err != nil {
		unexpectedDatabaseError(t, err)
	}
}

// TestDBStorageClusterOperation2 check the DB operations with clusers
func TestDBStorageClusterOperation2(t *testing.T) {
	const clusterName = "cluster_to_test_2"
	mockStorage, closer := MustGetMockStorage(t, true)
	defer closer()

	err := mockStorage.RegisterNewCluster(clusterName)
	if err != nil {
		unexpectedDatabaseError(t, err)
	}

	_, err = mockStorage.GetClusterByName(clusterName)
	if err != nil {
		unexpectedDatabaseError(t, err)
	}

	err = mockStorage.DeleteClusterByName(clusterName)
	if err != nil {
		unexpectedDatabaseError(t, err)
	}
}

// TestDBStorageClusterOperation3 check the DB operations with clusers
func TestDBStorageClusterOperation3(t *testing.T) {
	const clusterName = "cluster_to_test_3"
	mockStorage, closer := MustGetMockStorage(t, true)
	defer closer()

	err := mockStorage.RegisterNewCluster(clusterName)
	if err != nil {
		unexpectedDatabaseError(t, err)
	}

	clusters, err := mockStorage.ListOfClusters()
	if err != nil {
		unexpectedDatabaseError(t, err)
	}

	if len(clusters) != 1 {
		t.Fatal("Unexpected # of clusters:", len(clusters))
	}
}

// TestDBStorageConfigurationProfiles1
func TestDBStorageConfigurationProfiles1(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, true)
	defer closer()

	_, err := mockStorage.StoreConfigurationProfile("username1", "description1", "configuration1")
	if err != nil {
		unexpectedDatabaseError(t, err)
	}

	_, err = mockStorage.GetConfigurationProfile(1)
	if err != nil {
		unexpectedDatabaseError(t, err)
	}
}

// TestDBStorageConfigurationProfiles2
func TestDBStorageConfigurationProfiles2(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, true)
	defer closer()

	_, err := mockStorage.StoreConfigurationProfile("username1", "description1", "configuration1")
	if err != nil {
		unexpectedDatabaseError(t, err)
	}

	profiles, err := mockStorage.ListConfigurationProfiles()
	if err != nil {
		unexpectedDatabaseError(t, err)
	}

	if len(profiles) != 1 {
		t.Fatal("Incorrect number of configuration profiles", len(profiles))
	}
}

// TestDBStorageConfigurationProfiles3
func TestDBStorageConfigurationProfiles3(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, true)
	defer closer()

	_, err := mockStorage.StoreConfigurationProfile("username1", "description1", "configuration1")
	if err != nil {
		unexpectedDatabaseError(t, err)
	}

	_, err = mockStorage.ChangeConfigurationProfile(1, "username1", "description12", "configuration12")
	if err != nil {
		unexpectedDatabaseError(t, err)
	}
}

// TestDBStorageConfigurationProfiles4
func TestDBStorageConfigurationProfiles4(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, true)
	defer closer()

	_, err := mockStorage.StoreConfigurationProfile("username1", "description1", "configuration1")
	if err != nil {
		unexpectedDatabaseError(t, err)
	}

	_, err = mockStorage.DeleteConfigurationProfile(1)
	if err != nil {
		unexpectedDatabaseError(t, err)
	}
}

// TestDBClusterConfiguration1
func TestDBClusterConfiguration1(t *testing.T) {
	const clusterName = "cluster2"
	mockStorage, closer := MustGetMockStorage(t, true)
	defer closer()

	err := mockStorage.RegisterNewCluster(clusterName)
	if err != nil {
		unexpectedDatabaseError(t, err)
	}

	_, err = mockStorage.CreateClusterConfiguration(clusterName, "user1", "reason1", "description1", "configuration1")
	if err != nil {
		unexpectedDatabaseError(t, err)
	}
}

// TestDBClusterConfiguration2
func TestDBClusterConfiguration2(t *testing.T) {
	const clusterName = "cluster2"
	mockStorage, closer := MustGetMockStorage(t, true)
	defer closer()

	err := mockStorage.RegisterNewCluster(clusterName)
	if err != nil {
		unexpectedDatabaseError(t, err)
	}

	_, err = mockStorage.CreateClusterConfiguration(clusterName, "user1", "reason1", "description1", "configuration1")
	if err != nil {
		unexpectedDatabaseError(t, err)
	}

	err = mockStorage.EnableOrDisableClusterConfigurationByID(1, "0")
	if err != nil {
		unexpectedDatabaseError(t, err)
	}
}

// TestDBClusterConfiguration3
func TestDBClusterConfiguration3(t *testing.T) {
	const clusterName = "cluster2"
	mockStorage, closer := MustGetMockStorage(t, true)
	defer closer()

	err := mockStorage.RegisterNewCluster(clusterName)
	if err != nil {
		unexpectedDatabaseError(t, err)
	}

	_, err = mockStorage.CreateClusterConfiguration(clusterName, "user1", "reason1", "description1", "configuration1")
	if err != nil {
		unexpectedDatabaseError(t, err)
	}

	_, err = mockStorage.EnableClusterConfiguration(clusterName, "user1", "reason1")
	if err != nil {
		unexpectedDatabaseError(t, err)
	}
}

// TestDBClusterConfiguration4
func TestDBClusterConfiguration4(t *testing.T) {
	const clusterName = "cluster2"
	mockStorage, closer := MustGetMockStorage(t, true)
	defer closer()

	err := mockStorage.RegisterNewCluster(clusterName)
	if err != nil {
		unexpectedDatabaseError(t, err)
	}

	_, err = mockStorage.CreateClusterConfiguration(clusterName, "user1", "reason1", "description1", "configuration1")
	if err != nil {
		unexpectedDatabaseError(t, err)
	}

	_, err = mockStorage.DisableClusterConfiguration(clusterName, "user2", "reason2")
	if err != nil {
		unexpectedDatabaseError(t, err)
	}
}

// TestDBClusterConfiguration5
func TestDBClusterConfiguration5(t *testing.T) {
	const clusterName = "cluster2"
	mockStorage, closer := MustGetMockStorage(t, true)
	defer closer()

	err := mockStorage.RegisterNewCluster(clusterName)
	if err != nil {
		unexpectedDatabaseError(t, err)
	}

	_, err = mockStorage.CreateClusterConfiguration(clusterName, "user1", "reason1", "description1", "configuration1")
	if err != nil {
		unexpectedDatabaseError(t, err)
	}

	_, err = mockStorage.GetClusterConfigurationByID(1)
	if err != nil {
		unexpectedDatabaseError(t, err)
	}
}

// TestDBClusterConfiguration6
func TestDBClusterConfiguration6(t *testing.T) {
	const clusterName = "cluster2"
	mockStorage, closer := MustGetMockStorage(t, true)
	defer closer()

	err := mockStorage.RegisterNewCluster(clusterName)
	if err != nil {
		unexpectedDatabaseError(t, err)
	}

	_, err = mockStorage.CreateClusterConfiguration(clusterName, "user1", "reason1", "description1", "configuration1")
	if err != nil {
		unexpectedDatabaseError(t, err)
	}

	_, err = mockStorage.GetClusterActiveConfiguration(clusterName)
	if err != nil {
		unexpectedDatabaseError(t, err)
	}
}

// TestDBStorageTrigger1
func TestDBStorageTrigger1(t *testing.T) {
	const clusterName = "cluster3"
	const triggerType = "trigger-type-Y"

	mockStorage, closer := MustGetMockStorage(t, true)
	defer closer()

	err := mockStorage.RegisterNewCluster(clusterName)
	if err != nil {
		unexpectedDatabaseError(t, err)
	}

	err = mockStorage.NewTriggerType(triggerType, "description-of-new-trigger-type")
	if err != nil {
		unexpectedDatabaseError(t, err)
	}

	err = mockStorage.NewTrigger(clusterName, triggerType, "user3", "reason3", "link3")
	if err != nil {
		unexpectedDatabaseError(t, err)
	}
}

// TestDBStorageTrigger2
func TestDBStorageTrigger2(t *testing.T) {
	const clusterName = "cluster3"
	const triggerType = "trigger-type-Y"

	mockStorage, closer := MustGetMockStorage(t, true)
	defer closer()

	err := mockStorage.RegisterNewCluster(clusterName)
	if err != nil {
		unexpectedDatabaseError(t, err)
	}

	err = mockStorage.NewTriggerType(triggerType, "description-of-new-trigger-type")
	if err != nil {
		unexpectedDatabaseError(t, err)
	}

	err = mockStorage.NewTrigger(clusterName, triggerType, "user3", "reason3", "link3")
	if err != nil {
		unexpectedDatabaseError(t, err)
	}

	_, err = mockStorage.ListClusterTriggers(clusterName)
	if err != nil {
		unexpectedDatabaseError(t, err)
	}
}

// TestDBStorageTrigger3
func TestDBStorageTrigger3(t *testing.T) {
	const clusterName = "cluster3"
	const triggerType = "trigger-type-Y"

	mockStorage, closer := MustGetMockStorage(t, true)
	defer closer()

	err := mockStorage.RegisterNewCluster(clusterName)
	if err != nil {
		unexpectedDatabaseError(t, err)
	}

	err = mockStorage.NewTriggerType(triggerType, "description-of-new-trigger-type")
	if err != nil {
		unexpectedDatabaseError(t, err)
	}

	err = mockStorage.NewTrigger(clusterName, triggerType, "user3", "reason3", "link3")
	if err != nil {
		unexpectedDatabaseError(t, err)
	}

	_, err = mockStorage.ListActiveClusterTriggers(clusterName)
	if err != nil {
		unexpectedDatabaseError(t, err)
	}
}

// TestDBStorageTrigger4
func TestDBStorageTrigger4(t *testing.T) {
	const clusterName = "cluster3"
	const triggerType = "trigger-type-Y"

	mockStorage, closer := MustGetMockStorage(t, true)
	defer closer()

	err := mockStorage.RegisterNewCluster(clusterName)
	if err != nil {
		unexpectedDatabaseError(t, err)
	}

	err = mockStorage.NewTriggerType(triggerType, "description-of-new-trigger-type")
	if err != nil {
		unexpectedDatabaseError(t, err)
	}

	err = mockStorage.NewTrigger(clusterName, triggerType, "user3", "reason3", "link3")
	if err != nil {
		unexpectedDatabaseError(t, err)
	}

	err = mockStorage.AckTrigger(clusterName, 1)
	if err != nil {
		unexpectedDatabaseError(t, err)
	}

	err = mockStorage.ChangeStateOfTriggerByID(1, 0)
	if err != nil {
		unexpectedDatabaseError(t, err)
	}
}

func TestDBPlaceholder(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, true)
	defer closer()

	p := mockStorage.Placeholder()
	if p == nil {
		t.Fatal("Placeholder should not be empty")
	}

}

func TestDBConnections(t *testing.T) {
	mockStorage, closer := MustGetMockStorage(t, true)
	defer closer()

	c := mockStorage.Connections()
	if c == nil {
		t.Fatal("Connections should not be empty")
	}
}
