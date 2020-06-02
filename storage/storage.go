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

package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/lib/pq"           // PostgreSQL database driver
	_ "github.com/mattn/go-sqlite3" // SQLite database driver
)

// Storage represents an interface to any relational database based on SQL language
type Storage struct {
	connections *sql.DB
	driver      string
	placeholder sq.PlaceholderFormat
}

// Column is typed reference to a sql column, which is further used by particular storage objects
type Column string

func enableForeignKeys(connections *sql.DB) {
	log.Println("Enabling foreign_keys pragma for sqlite")
	statement, err := connections.Prepare("PRAGMA foreign_keys = ON")
	if err != nil {
		log.Fatal("Can prepare statement set PRAGMA for sqlite", err)
	}

	// statement has to be closed at function exit
	defer func() {
		// try to close the statement
		err := statement.Close()
		// in case of error all we can do is to just log the error
		if err != nil {
			log.Println(err)
		}
	}()

	_, err = statement.Exec()
	if err != nil {
		log.Fatal("Can not set PRAGMA for sqlite", err)
	}
}

// New function creates and initializes a new instance of Storage structure
func New(driverName string, dataSourceName string) (Storage, error) {
	log.Printf("Making connection to data storage, driver=%s datasource=%s", driverName, dataSourceName)
	connections, err := sql.Open(driverName, dataSourceName)

	if err != nil {
		log.Println("Can not connect to data storage", err)
		return Storage{}, err
	}
	s := Storage{connections: connections, driver: driverName}

	switch driverName {
	case "sqlite3":
		enableForeignKeys(connections)
		s.placeholder = sq.Question
	case "postgres":
		s.placeholder = sq.Dollar
	}
	return s, nil
}

// NewFromConnection function creates and initializes a new instance of Storage interface from prepared connection
func NewFromConnection(connection *sql.DB, driverName string) Storage {
	return Storage{
		connections: connection,
		driver:      driverName,
	}
}

// Placeholder returns current query argument placeholder
// (?, or $).It depends on driver used. In squirrel format
func (storage Storage) Placeholder() sq.PlaceholderFormat {
	return storage.placeholder
}

// Connections is sql.DB connection
func (storage Storage) Connections() *sql.DB {
	return storage.connections
}

// Close method closes the connection to database. Needs to be called at the end of application lifecycle.
func (storage Storage) Close() {
	log.Println("Closing connection to data storage")
	if storage.connections != nil {
		err := storage.connections.Close()
		if err != nil {
			log.Fatal("Can not close connection to data storage", err)
		}
	}
}

// ID represents unique ID for any object.
type ID int

// Name represents common name of object stored in database.
type Name string

// ClusterID represents unique key of cluster stored in database.
type ClusterID ID

// ClusterName represents name of cluster in format c8590f31-e97e-4b85-b506-c45ce1911a12
type ClusterName Name

// Cluster represents cluster record in the controller service.
//     ID: unique key
//     Name: cluster GUID in the following format:
//         c8590f31-e97e-4b85-b506-c45ce1911a12
type Cluster struct {
	ID   ClusterID   `json:"id"`
	Name ClusterName `json:"name"`
}

// ConfigurationID represents unique key of configuration stored in database.
type ConfigurationID ID

// ConfigurationProfile represents configuration profile record in the controller service.
//     ID: unique key
//     Configuration: a JSON structure stored in a string
//     ChangeAt: username of admin that created or updated the configuration
//     ChangeBy: timestamp of the last configuration change
//     Description: a string with any comment(s) about the configuration
type ConfigurationProfile struct {
	ID            ConfigurationID `json:"id"`
	Configuration string          `json:"configuration"`
	ChangedAt     string          `json:"changed_at"`
	ChangedBy     string          `json:"changed_by"`
	Description   string          `json:"description"`
}

// ClusterConfigurationID represents unique key of cluster configuration stored in database.
type ClusterConfigurationID ID

// ClusterConfiguration represents cluster configuration record in the controller service.
//     ID: unique key
//     Cluster: cluster ID (not name)
//     Configuration: a JSON structure stored in a string
//     ChangeAt: timestamp of the last configuration change
//     ChangeBy: username of admin that created or updated the configuration
//     Active: flag indicating whether the configuration is active or not
//     Reason: a string with any comment(s) about the cluster configuration
type ClusterConfiguration struct {
	ID            ClusterConfigurationID `json:"id"`
	Cluster       string                 `json:"cluster"`
	Configuration string                 `json:"configuration"`
	ChangedAt     string                 `json:"changed_at"`
	ChangedBy     string                 `json:"changed_by"`
	Active        string                 `json:"active"`
	Reason        string                 `json:"reason"`
}

// TriggerID represents unique key of trigger stored in database.
type TriggerID ID

// Trigger represents trigger record in the controller service
//     ID: unique key
//     Type: ID of trigger type
//     Cluster: cluster ID (not name)
//     Reason: a string with any comment(s) about the trigger
//     Link: link to any document with customer ACK with the trigger
//     TriggeredAt: timestamp of the last configuration change
//     TriggeredBy: username of admin that created or updated the trigger
//     AckedAt: timestamp where the insights operator acked the trigger
//     Parameters: parameters that needs to be pass to trigger code
//     Active: flag indicating whether the trigger is still active or not
type Trigger struct {
	ID          TriggerID `json:"id"`
	Type        string    `json:"type"`
	Cluster     string    `json:"cluster"`
	Reason      string    `json:"reason"`
	Link        string    `json:"link"`
	TriggeredAt string    `json:"triggered_at"`
	TriggeredBy string    `json:"triggered_by"`
	AckedAt     string    `json:"acked_at"`
	Parameters  string    `json:"parameters"`
	Active      int       `json:"active"`
}

// ListOfClusters method selects all clusters from database.
func (storage Storage) ListOfClusters() ([]Cluster, error) {
	clusters := []Cluster{}

	rows, err := storage.connections.Query("SELECT id, name FROM cluster")
	if err != nil {
		return clusters, err
	}

	// close the query at function exit
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	for rows.Next() {
		var id int
		var name string

		err = rows.Scan(&id, &name)
		if err == nil {
			clusters = append(clusters, Cluster{ClusterID(id), ClusterName(name)})
		} else {
			log.Println("error", err)
		}
	}
	return clusters, nil
}

// GetCluster method selects the specified cluster from database. Also see GetClusterByName.
func (storage Storage) GetCluster(id int) (Cluster, error) {
	var cluster Cluster

	rows, err := storage.connections.Query("SELECT id, name FROM cluster WHERE id = $1", id)
	if err != nil {
		return cluster, err
	}
	defer rows.Close()

	if rows.Next() {
		var id int
		var name string

		err = rows.Scan(&id, &name)
		if err == nil {
			cluster.ID = ClusterID(id)
			cluster.Name = ClusterName(name)
		} else {
			log.Println("error", err)
		}
	} else {
		return cluster, &ItemNotFoundError{
			ItemID: id,
		}
	}
	return cluster, err
}

// RegisterNewCluster inserts information about new cluster into the database.
// It differs from CreateNewCluster, because ID is not specified explicitly here.
func (storage Storage) RegisterNewCluster(name string) error {
	statement, err := storage.connections.Prepare("INSERT INTO cluster(name) VALUES ($1)")
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(name)
	return err
}

// CreateNewCluster creates a new cluster with specified ID and name.
// It differs from RegisterNewCluster, because ID is specified explicitly here.
func (storage Storage) CreateNewCluster(id int64, name string) error {
	statement, err := storage.connections.Prepare("INSERT INTO cluster(id, name) VALUES ($1, $2)")
	if err != nil {
		log.Print(err)
		return err
	}

	// statement has to be closed at function exit
	defer func() {
		// try to close the statement
		err := statement.Close()
		// in case of error all we can do is to just log the error
		if err != nil {
			log.Println(err)
		}
	}()

	_, err = statement.Exec(id, name)
	return err
}

// DeleteCluster deletes cluster with specified ID from the database.
func (storage Storage) DeleteCluster(id int64) error {
	statement, err := storage.connections.Prepare("DELETE FROM cluster WHERE id = $1")
	if err != nil {
		log.Print(err)
		return err
	}

	// statement has to be closed at function exit
	defer func() {
		// try to close the statement
		err := statement.Close()
		// in case of error all we can do is to just log the error
		if err != nil {
			log.Println(err)
		}
	}()

	rowsAffected, err := execStatementAndGetRowsAffected(statement, id)
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return &ItemNotFoundError{
			ItemID: id,
		}
	}
	return nil
}

// DeleteClusterByName deletes cluster with specified name from the database.
func (storage Storage) DeleteClusterByName(name string) error {
	statement, err := storage.connections.Prepare("DELETE FROM cluster WHERE name = $1")
	if err != nil {
		log.Print(err)
		return err
	}

	// statement has to be closed at function exit
	defer func() {
		// try to close the statement
		err := statement.Close()
		// in case of error all we can do is to just log the error
		if err != nil {
			log.Println(err)
		}
	}()

	rowsAffected, err := execStatementAndGetRowsAffected(statement, name)
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return &ItemNotFoundError{
			ItemID: name,
		}
	}
	return nil
}

// GetClusterByName selects a cluster specified by its name. Also see GetCluster.
func (storage Storage) GetClusterByName(name string) (Cluster, error) {
	var cluster Cluster

	rows, err := storage.connections.Query("SELECT id, name FROM cluster WHERE name = $1", name)
	if err != nil {
		log.Print(err)
		return cluster, err
	}

	// close the query at function exit
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	if rows.Next() {
		var id int
		var name string

		err = rows.Scan(&id, &name)
		if err == nil {
			cluster.ID = ClusterID(id)
			cluster.Name = ClusterName(name)
			log.Printf("Cluster name %s has id %d\n", name, id)
		} else {
			log.Println("error", err)
		}
	} else {
		return cluster, &ItemNotFoundError{
			ItemID: name,
		}
	}
	return cluster, err
}

// ListConfigurationProfiles selects list of all configuration profiles from database.
func (storage Storage) ListConfigurationProfiles() ([]ConfigurationProfile, error) {
	profiles := []ConfigurationProfile{}

	rows, err := storage.connections.Query("SELECT id, configuration, changed_at, changed_by, description FROM configuration_profile")
	if err != nil {
		log.Print(err)
		return profiles, err
	}

	// close the query at function exit
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	for rows.Next() {
		var id int
		var configuration string
		var changedAt string
		var changedBy string
		var description string

		err = rows.Scan(&id, &configuration, &changedAt, &changedBy, &description)
		if err == nil {
			profiles = append(profiles, ConfigurationProfile{ConfigurationID(id), configuration, changedAt, changedBy, description})
		} else {
			log.Println("error", err)
		}
	}

	return profiles, nil
}

// GetConfigurationProfile selects one configuration profile identified by its ID.
func (storage Storage) GetConfigurationProfile(id int) (ConfigurationProfile, error) {
	var profile ConfigurationProfile

	rows, err := storage.connections.Query("SELECT id, configuration, changed_at, changed_by, description FROM configuration_profile WHERE id = $1", id)
	if err != nil {
		return profile, err
	}
	defer rows.Close()

	if rows.Next() {
		var id int
		var configuration string
		var changedAt string
		var changedBy string
		var description string

		err = rows.Scan(&id, &configuration, &changedAt, &changedBy, &description)
		if err == nil {
			profile.ID = ConfigurationID(id)
			profile.Configuration = configuration
			profile.ChangedAt = changedAt
			profile.ChangedBy = changedBy
			profile.Description = description
		} else {
			log.Println("error", err)
		}
	} else {
		return profile, &ItemNotFoundError{
			ItemID: id,
		}
	}
	return profile, err
}

// StoreConfigurationProfile stores a given configuration profile (string ATM) into the database.
func (storage Storage) StoreConfigurationProfile(username string, description string, configuration string) ([]ConfigurationProfile, error) {
	var profiles []ConfigurationProfile

	t := time.Now()

	statement, err := storage.connections.Prepare("INSERT INTO configuration_profile(configuration, changed_at, changed_by, description) VALUES ($1, $2, $3, $4)")
	if err != nil {
		log.Print(err)
		return profiles, err
	}

	// statement has to be closed at function exit
	defer func() {
		// try to close the statement
		err := statement.Close()
		// in case of error all we can do is to just log the error
		if err != nil {
			log.Println(err)
		}
	}()

	_, err = statement.Exec(configuration, t, username, description)
	if err != nil {
		log.Print(err)
		return profiles, err
	}

	return storage.ListConfigurationProfiles()
}

// ChangeConfigurationProfile updates the existing configuration profile specified by its ID.
func (storage Storage) ChangeConfigurationProfile(id int, username string, description string, configuration string) ([]ConfigurationProfile, error) {
	var profiles []ConfigurationProfile

	t := time.Now()

	statement, err := storage.connections.Prepare("UPDATE configuration_profile SET configuration = $1, changed_at = $2, changed_by = $3, description = $4 WHERE id = $5")
	if err != nil {
		log.Print(err)
		return profiles, err
	}
	defer statement.Close()

	rowsAffected, err := execStatementAndGetRowsAffected(statement, configuration, t, username, description, id)
	if err != nil {
		log.Print(err)
		return profiles, err
	}
	if rowsAffected == 0 {
		return profiles, &ItemNotFoundError{
			ItemID: id,
		}
	}

	return storage.ListConfigurationProfiles()
}

// DeleteConfigurationProfile deletes a configuration profile specified by its name.
func (storage Storage) DeleteConfigurationProfile(id int) ([]ConfigurationProfile, error) {
	var profiles []ConfigurationProfile

	statement, err := storage.connections.Prepare("DELETE FROM configuration_profile WHERE id = $1")
	if err != nil {
		log.Print(err)
		return profiles, err
	}
	defer statement.Close()

	rowsAffected, err := execStatementAndGetRowsAffected(statement, id)
	if err != nil {
		log.Print(err)
		return profiles, err
	}
	if rowsAffected == 0 {
		return profiles, &ItemNotFoundError{
			ItemID: id,
		}
	}

	return storage.ListConfigurationProfiles()
}

func (storage Storage) readClusterConfigurations(rows *sql.Rows) ([]ClusterConfiguration, error) {
	configurations := []ClusterConfiguration{}

	defer rows.Close()

	for rows.Next() {
		var id int
		var cluster string
		var configuration string
		var changedAt string
		var changedBy string
		var active string
		var reason string

		err := rows.Scan(&id, &cluster, &configuration, &changedAt, &changedBy, &active, &reason)
		if err == nil {
			configurations = append(configurations, ClusterConfiguration{ClusterConfigurationID(id), cluster, configuration, changedAt, changedBy, active, reason})
		} else {
			log.Println("error", err)
		}
	}

	return configurations, nil
}

// ListAllClusterConfigurations selects all cluster configurations from the database.
func (storage Storage) ListAllClusterConfigurations() ([]ClusterConfiguration, error) {
	rows, err := storage.connections.Query(`
SELECT operator_configuration.id, cluster.name, configuration, changed_at, changed_by, active, reason
  FROM operator_configuration JOIN cluster
    ON (cluster.id = operator_configuration.cluster)
ORDER BY operator_configuration.id`)

	if err != nil {
		log.Print(err)
		return []ClusterConfiguration{}, err
	}
	return storage.readClusterConfigurations(rows)
}

// ListClusterConfiguration selects cluster configuration from the database for the specified cluster.
func (storage Storage) ListClusterConfiguration(cluster string) ([]ClusterConfiguration, error) {
	if _, err := storage.GetClusterByName(cluster); err != nil {
		return nil, err
	}

	rows, err := storage.connections.Query(`
SELECT operator_configuration.id, cluster.name, configuration, changed_at, changed_by, active, reason
  FROM operator_configuration JOIN cluster
    ON (cluster.id = operator_configuration.cluster)
 WHERE cluster.name = $1`, cluster)

	if err != nil {
		log.Print(err)
		return []ClusterConfiguration{}, err
	}

	return storage.readClusterConfigurations(rows)
}

// GetClusterConfigurationByID reads cluster configuration for the specified configuration ID.
func (storage Storage) GetClusterConfigurationByID(id int64) (string, error) {
	var configuration string

	row, err := storage.connections.Query(`
SELECT configuration_profile.configuration
  FROM operator_configuration JOIN configuration_profile
    ON (configuration_profile.id = operator_configuration.configuration)
 WHERE operator_configuration.id = $1`, id)

	if err != nil {
		log.Print(err)
		return configuration, err
	}
	defer row.Close()

	if row.Next() {
		err = row.Scan(&configuration)
		if err != nil {
			log.Println("error", err)
		}
		return configuration, err
	}
	return configuration, &ItemNotFoundError{
		ItemID: id,
	}
}

// GetClusterActiveConfiguration reads one active configuration for the selected cluster.
func (storage Storage) GetClusterActiveConfiguration(cluster string) (string, error) {
	var configuration string

	row, err := storage.connections.Query(`
SELECT configuration_profile.configuration
  FROM operator_configuration, cluster, configuration_profile
 WHERE cluster.id = operator_configuration.cluster
   AND configuration_profile.id = operator_configuration.configuration
   AND operator_configuration.active = '1' AND cluster.name = $1
 LIMIT 1`, cluster)

	if err != nil {
		log.Print(err)
		return configuration, err
	}
	defer row.Close()

	if row.Next() {
		err = row.Scan(&configuration)
		if err != nil {
			log.Println("error", err)
		}
		return configuration, err
	}
	return configuration, &ItemNotFoundError{
		ItemID: cluster,
	}
}

// GetConfigurationIDForCluster reads the ID for the specified cluster name.
func (storage Storage) GetConfigurationIDForCluster(cluster string) (int, error) {
	rows, err := storage.connections.Query(`
SELECT operator_configuration.id
  FROM operator_configuration, cluster
    ON cluster.id = operator_configuration.cluster
 WHERE cluster.name = $1`, cluster)

	if err != nil {
		return 0, err
	}
	defer rows.Close()

	if rows.Next() {
		var id int

		err = rows.Scan(&id)
		return id, err
	}
	return 0, errors.New("Unknown operator name provided")
}

// InsertNewConfigurationProfile inserts new configuration profile into a database (in transaction).
func (storage Storage) InsertNewConfigurationProfile(tx *sql.Tx, configuration string, username string, description string) bool {
	t := time.Now()

	statement, err := tx.Prepare("INSERT INTO configuration_profile(configuration, changed_at, changed_by, description) VALUES ($1, $2, $3, $4)")
	if err != nil {
		return false
	}
	defer statement.Close()

	_, err = statement.Exec(configuration, t, username, description)
	if err != nil {
		return false
	}
	return true
}

// SelectConfigurationProfileID selects the ID of lately inserted/created configuration profile. To be used in transaction.
func (storage Storage) SelectConfigurationProfileID(tx *sql.Tx) (int, error) {
	var rows *sql.Rows
	var err error

	// We need to get the ID from the last insert. Unfortunately it seems there is not
	// one existing solution that works for all databases.
	switch storage.driver {
	case "sqlite3":
		rows, err = tx.Query(`SELECT rowid FROM configuration_profile ORDER BY rowid DESC limit 1`)
	case "postgres":
		rows, err = tx.Query(`SELECT currval('configuration_profile_id_seq')`)
	default:
		return -1, errors.New("unknown DB driver:" + storage.driver)
	}
	if err != nil {
		log.Print(err)
		return -1, err
	}
	defer rows.Close()

	if rows.Next() {
		var configurationID int
		err = rows.Scan(&configurationID)
		if err != nil {
			return -1, err
		}
		log.Printf("Configuration stored under ID=%d\n", configurationID)
		return configurationID, nil
	}
	return -1, errors.New("can not retrieve last configuration ID")
}

// DeactivatePreviousConfigurations deactivate all previous configurations for the specified trigger.
// To be called inside transaction.
func (storage Storage) DeactivatePreviousConfigurations(tx *sql.Tx, clusterID ClusterID) error {
	stmt, err := tx.Prepare("UPDATE operator_configuration SET active=0 WHERE cluster = $1")
	defer stmt.Close()

	if err != nil {
		return err
	}
	_, err = stmt.Exec(clusterID)
	if err == nil {
		log.Printf("All previous configuration has been deactivated for clusterID %d\n", clusterID)
	}
	return err
}

// InsertNewOperatorConfiguration inserts the new configuration for selected operator/cluster.
// To be called inside transaction.
func (storage Storage) InsertNewOperatorConfiguration(tx *sql.Tx, clusterID ClusterID, configurationID int, username string, reason string) error {
	t := time.Now()
	statement, err := tx.Prepare("INSERT INTO operator_configuration(cluster, configuration, changed_at, changed_by, active, reason) VALUES ($1, $2, $3, $4, $5, $6)")
	defer statement.Close()
	if err != nil {
		return err
	}

	_, err = statement.Exec(clusterID, configurationID, t, username, "1", reason)
	if err == nil {
		log.Printf("New operator configuration %d has been assigned to cluster %d\n", configurationID, clusterID)
	}
	return err
}

// CreateClusterConfiguration creates new configuration for specified cluster.
func (storage Storage) CreateClusterConfiguration(cluster string, username string, reason string, description, configuration string) ([]ClusterConfiguration, error) {
	// retrieve cluster ID
	clusterInfo, err := storage.GetClusterByName(cluster)

	if err != nil {
		log.Print(err)
		return []ClusterConfiguration{}, err
	}

	clusterID := clusterInfo.ID

	// begin transaction
	tx, err := storage.connections.Begin()
	if err != nil {
		log.Print(err)
		log.Println("Transaction failed")
		return []ClusterConfiguration{}, err
	}

	// insert new configuration profile
	if !storage.InsertNewConfigurationProfile(tx, configuration, username, description) {
		log.Print(err)
		_ = tx.Rollback()
		return []ClusterConfiguration{}, err
	}

	// retrieve configuration ID for newly created configuration
	configurationID, err := storage.SelectConfigurationProfileID(tx)
	if err != nil {
		log.Print(err)
		_ = tx.Rollback()
		return []ClusterConfiguration{}, err
	}

	// deactivate all previous configurations
	err = storage.DeactivatePreviousConfigurations(tx, clusterID)
	if err != nil {
		log.Print(err)
		_ = tx.Rollback()
		return []ClusterConfiguration{}, err
	}

	// and insert new one that will be activated
	err = storage.InsertNewOperatorConfiguration(tx, clusterID, configurationID, username, reason)
	if err != nil {
		log.Print(err)
		_ = tx.Rollback()
		return []ClusterConfiguration{}, err
	}

	// end the transaction
	if err := tx.Commit(); err != nil {
		log.Print(err)
		return []ClusterConfiguration{}, err
	}

	return storage.ListClusterConfiguration(cluster)
}

// EnableClusterConfiguration enables the specified cluster configuration (set the 'active' flag).
func (storage Storage) EnableClusterConfiguration(cluster string, username string, reason string) ([]ClusterConfiguration, error) {
	id, err := storage.GetConfigurationIDForCluster(cluster)
	if err != nil {
		return []ClusterConfiguration{}, err
	}

	statement, err := storage.connections.Prepare("UPDATE operator_configuration SET active=1, changed_at = $1, changed_by = $2, reason = $3 WHERE id = $4")
	if err != nil {
		return []ClusterConfiguration{}, err
	}
	defer statement.Close()

	t := time.Now()

	_, err = statement.Exec(t, username, reason, id)
	if err != nil {
		return []ClusterConfiguration{}, err
	}
	return storage.ListClusterConfiguration(cluster)
}

// DisableClusterConfiguration disables the specified cluster configuration (reset the 'active' flag).
// TODO: copy & paste, needs to be refactored later
func (storage Storage) DisableClusterConfiguration(cluster string, username string, reason string) ([]ClusterConfiguration, error) {
	id, err := storage.GetConfigurationIDForCluster(cluster)
	if err != nil {
		return []ClusterConfiguration{}, err
	}
	statement, err := storage.connections.Prepare("UPDATE operator_configuration SET active=0, changed_at = $1, changed_by = $2, reason = $3 WHERE id = $4")
	if err != nil {
		return []ClusterConfiguration{}, err
	}
	defer statement.Close()

	t := time.Now()

	_, err = statement.Exec(t, username, reason, id)
	if err != nil {
		return []ClusterConfiguration{}, err
	}
	return storage.ListClusterConfiguration(cluster)
}

// EnableOrDisableClusterConfigurationByID enables or disables the specified cluster configuration (set or reset the 'active' flag).
// Please see also EnableClusterConfiguration and DisableClusterConfiguration
func (storage Storage) EnableOrDisableClusterConfigurationByID(id int64, active string) error {
	statement, err := storage.connections.Prepare("UPDATE operator_configuration SET active = $1, changed_at = $2 WHERE id = $3")
	if err != nil {
		return err
	}
	defer statement.Close()

	t := time.Now()

	rowsAffected, err := execStatementAndGetRowsAffected(statement, active, t, id)
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return &ItemNotFoundError{
			ItemID: id,
		}
	}
	return nil
}

// DeleteClusterConfigurationByID deletes cluster configuration specified by its ID.
// TODO: copy & paste, needs to be refactored later
func (storage Storage) DeleteClusterConfigurationByID(id int64) error {
	statement, err := storage.connections.Prepare("DELETE FROM operator_configuration WHERE id = $1")
	if err != nil {
		return err
	}
	defer statement.Close()

	rowsAffected, err := execStatementAndGetRowsAffected(statement, id)
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return &ItemNotFoundError{
			ItemID: id,
		}
	}
	return nil
}

func (storage Storage) getTriggers(rows *sql.Rows) ([]Trigger, error) {
	triggers := []Trigger{}
	defer rows.Close()

	for rows.Next() {
		var trigger Trigger

		err := rows.Scan(&trigger.ID, &trigger.Type, &trigger.Cluster,
			&trigger.Reason, &trigger.Link,
			&trigger.TriggeredAt, &trigger.TriggeredBy,
			&trigger.Parameters, &trigger.Active, &trigger.AckedAt)
		if err == nil {
			triggers = append(triggers, trigger)
		} else {
			log.Println("error", err)
		}
	}

	return triggers, nil
}

// GetTriggerByID selects all informations about the trigger specified by its ID.
func (storage Storage) GetTriggerByID(id int64) (Trigger, error) {
	rows, err := storage.connections.Query(`
SELECT trigger.id, trigger_type.type, cluster.name,
       trigger.reason, trigger.link, trigger.triggered_at, trigger.triggered_by,
       trigger.parameters, trigger.active, trigger.acked_at
  FROM trigger JOIN trigger_type ON trigger.type=trigger_type.id
               JOIN cluster ON trigger.cluster=cluster.id
 WHERE trigger.id = $1`, id)

	if err != nil {
		return Trigger{}, err
	}

	triggers, err := storage.getTriggers(rows)
	if err != nil {
		return Trigger{}, err
	}

	if len(triggers) >= 1 {
		return triggers[0], nil
	}
	return Trigger{}, ErrNoSuchObj
}

func execStatementAndGetRowsAffected(statement *sql.Stmt, args ...interface{}) (int64, error) {
	res, err := statement.Exec(args...)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

// DeleteTriggerByID deletes trigger specified by its ID
// returns ItemNotFoundError if trigger didn't exist
func (storage Storage) DeleteTriggerByID(id int64) error {
	statement, err := storage.connections.Prepare(`
DELETE FROM trigger WHERE trigger.id = $1`)
	if err != nil {
		log.Print(err)
		return err
	}

	defer statement.Close()

	rowsAffected, err := execStatementAndGetRowsAffected(statement, id)
	if err != nil {
		log.Print(err)
		return err
	}

	// non-existent trigger ID has been used
	if rowsAffected == 0 {
		// convert ID (numeric value) to string for proper logging
		IDstr := strconv.Itoa(int(id))
		return &ItemNotFoundError{
			ItemID: IDstr,
		}
	}

	return nil
}

// ChangeStateOfTriggerByID change the state ('active', 'inactive') of trigger specified by its ID.
// returns ItemNotFoundError if there weren't rows with such id
func (storage Storage) ChangeStateOfTriggerByID(id int64, active int) error {
	statement, err := storage.connections.Prepare(`
UPDATE trigger SET active = $1 WHERE trigger.id = $2`)
	if err != nil {
		return err
	}

	// statement has to be closed at function exit
	defer func() {
		// try to close the statement
		err := statement.Close()
		// in case of error all we can do is to just log the error
		if err != nil {
			log.Println(err)
		}
	}()

	rowsAffected, err := execStatementAndGetRowsAffected(statement, active, id)
	if err != nil {
		log.Print(err)
		return err
	}

	// non-existent trigger ID has been used
	if rowsAffected == 0 {
		// convert ID (numeric value) to string for proper logging
		IDstr := strconv.Itoa(int(id))
		return &ItemNotFoundError{
			ItemID: IDstr,
		}
	}

	return nil
}

// ListAllTriggers selects all triggers from the database.
func (storage Storage) ListAllTriggers() ([]Trigger, error) {
	triggers := []Trigger{}

	rows, err := storage.connections.Query(`
SELECT trigger.id, trigger_type.type, cluster.name,
       trigger.reason, trigger.link, trigger.triggered_at, trigger.triggered_by,
       trigger.parameters, trigger.active, trigger.acked_at
FROM trigger JOIN trigger_type ON trigger.type=trigger_type.id
               JOIN cluster ON trigger.cluster=cluster.id
ORDER BY trigger.id`)

	if err != nil {
		return triggers, err
	}

	return storage.getTriggers(rows)
}

// ListClusterTriggers selects all triggers assigned to the specified cluster.
func (storage Storage) ListClusterTriggers(clusterName string) ([]Trigger, error) {
	triggers := []Trigger{}

	// check that cluster exist
	if _, err := storage.GetClusterByName(clusterName); err != nil {
		return triggers, err
	}

	rows, err := storage.connections.Query(`
SELECT trigger.id, trigger_type.type, cluster.name,
       trigger.reason, trigger.link, trigger.triggered_at, trigger.triggered_by,
       trigger.parameters, trigger.active, trigger.acked_at
  FROM trigger JOIN trigger_type ON trigger.type=trigger_type.id
               JOIN cluster ON trigger.cluster=cluster.id
 WHERE cluster.name = $1
 ORDER BY trigger.id`, clusterName)

	if err != nil {
		return triggers, err
	}

	return storage.getTriggers(rows)
}

// ListActiveClusterTriggers selects all active triggers assigned to the specified cluster.
func (storage Storage) ListActiveClusterTriggers(clusterName string) ([]Trigger, error) {
	triggers := []Trigger{}

	// check that cluster exist
	if _, err := storage.GetClusterByName(clusterName); err != nil {
		return triggers, err
	}

	rows, err := storage.connections.Query(`
SELECT trigger.id, trigger_type.type, cluster.name,
       trigger.reason, trigger.link, trigger.triggered_at, trigger.triggered_by,
       trigger.parameters, trigger.active, trigger.acked_at
  FROM trigger JOIN trigger_type ON trigger.type=trigger_type.id
               JOIN cluster ON trigger.cluster=cluster.id
 WHERE trigger.active = 1
   AND cluster.name = $1`, clusterName)

	if err != nil {
		return triggers, err
	}

	return storage.getTriggers(rows)
}

// GetTriggerID select ID for specified trigger type (name).
func (storage Storage) GetTriggerID(triggerType string) (int, error) {
	var id int

	rows, err := storage.connections.Query("SELECT id FROM trigger_type WHERE type = $1", triggerType)
	if err != nil {
		return 0, err
	}

	// rows has to be closed at function exit
	defer func() {
		// try to close the statement
		err := rows.Close()
		// in case of error all we can do is to just log the error
		if err != nil {
			log.Println(err)
		}
	}()

	if rows.Next() {

		err = rows.Scan(&id)
		if err == nil {
			log.Printf("Trigger type %s has id %d\n", triggerType, id)
		} else {
			log.Println("error", err)
		}
	} else {
		return 0, errors.New("Unknown trigger type provided")
	}
	return id, err
}

// NewTrigger constructs new trigger in a database.
func (storage Storage) NewTrigger(clusterName string, triggerType string, userName string, reason string, link string) error {
	// retrieve cluster ID
	clusterInfo, err := storage.GetClusterByName(clusterName)
	clusterID := clusterInfo.ID

	if err != nil {
		log.Print(err)
		return err
	}

	triggerTypeID, err := storage.GetTriggerID(triggerType)

	if err != nil {
		log.Print(err)
		return err
	}
	t := time.Now()
	ackedAt := time.Unix(0, 0).UTC()

	statement, err := storage.connections.Prepare("INSERT INTO trigger(type, cluster, reason, link, triggered_at, triggered_by, parameters, active, acked_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)")
	if err != nil {
		log.Print(err)
		return err
	}

	// statement has to be closed at function exit
	defer func() {
		// try to close the statement
		err := statement.Close()
		// in case of error all we can do is to just log the error
		if err != nil {
			log.Println(err)
		}
	}()

	_, err = statement.Exec(triggerTypeID, clusterID, reason, link, t, userName, "", 1, ackedAt)
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}

// NewTriggerType inserts a trigger_type object in the database
func (storage Storage) NewTriggerType(ttype string, description string) error {
	statement, err := storage.connections.Prepare("INSERT INTO trigger_type(type, description) VALUES ($1, $2)")
	if err != nil {
		log.Print(err)
		return err
	}

	// statement has to be closed at function exit
	defer func() {
		// try to close the statement
		err := statement.Close()
		// in case of error all we can do is to just log the error
		if err != nil {
			log.Println(err)
		}
	}()

	_, err = statement.Exec(ttype, description)
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}

// AckTrigger sets a timestamp to the selected trigger + updates the 'active' flag.
// and returns error if trigger wasn't found
func (storage Storage) AckTrigger(clusterName string, triggerID int64) error {
	t := time.Now()

	// retrieve cluster ID
	clusterInfo, err := storage.GetClusterByName(clusterName)
	clusterID := clusterInfo.ID

	if err != nil {
		return err
	}

	statement, err := storage.connections.Prepare("UPDATE trigger SET acked_at = $1, active=0 WHERE cluster = $2 AND id = $3")

	if err != nil {
		return err
	}

	// statement has to be closed at function exit
	defer func() {
		// try to close the statement
		err := statement.Close()
		// in case of error all we can do is to just log the error
		if err != nil {
			log.Println(err)
		}
	}()

	rowsAffected, err := execStatementAndGetRowsAffected(statement, t, clusterID, triggerID)
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return &ItemNotFoundError{
			ItemID: fmt.Sprintf("%v/%v", clusterName, triggerID),
		}
	}
	return nil
}

// QueryOne is generating Sql query using squirell sql builder, querying it with db store and mapping result to destination object with provided mapper
func (storage Storage) QueryOne(ctx context.Context, selectCols []Column, selectBuilder sq.SelectBuilder, mapper func(Column, interface{}) (interface{}, error), res interface{}) error {
	q, args, err := selectBuilder.ToSql()
	if err != nil {
		return err
	}
	rowScanner := storage.connections.QueryRowContext(ctx, q, args...)
	if rowScanner == nil {
		return ErrUnknown
	}

	resMap, err := storage.Map(selectCols, mapper, res)
	if err != nil {
		return err
	}
	err = rowScanner.Scan(resMap...)
	if err == sql.ErrNoRows {
		return ErrNoSuchObj
	}
	if err != nil {
		return err
	}
	return nil
}

// Map creates a list of destination struct fields using columns to select
func (storage Storage) Map(cols []Column, mapper func(Column, interface{}) (interface{}, error), r interface{}) ([]interface{}, error) {
	var mappedCols []interface{}
	for _, c := range cols {
		mc, err := mapper(c, r)
		if err != nil {
			return nil, err
		}
		mappedCols = append(mappedCols, mc)
	}
	return mappedCols, nil
}

// Ping checks whether the database connection is really configured properly
func (storage Storage) Ping() error {
	rows, err := storage.connections.Query("SELECT id, name FROM cluster LIMIT 1")
	if err != nil {
		return err
	}

	err = rows.Close()
	if err != nil {
		return err
	}

	return nil
}

// ErrNoSuchObj is indicating no result returned from db
var ErrNoSuchObj = fmt.Errorf("no such object")

// ErrUnknown indicates db query failed without details (QueryRow returning Row wasn't populated)
var ErrUnknown = fmt.Errorf("unknown error during querying db")
