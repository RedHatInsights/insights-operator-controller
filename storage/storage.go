/*
Copyright © 2019 Red Hat, Inc.

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
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"time"
)

type Storage struct {
	connections *sql.DB
	driver      string
}

func enableForeignKeys(connections *sql.DB) {
	log.Println("Enabling foreign_keys pragma for sqlite")
	statement, err := connections.Prepare("PRAGMA foreign_keys = ON")
	if err != nil {
		log.Fatal("Can prepare statement set PRAGMA for sqlite", err)
	}
	defer statement.Close()
	_, err = statement.Exec()
	if err != nil {
		log.Fatal("Can not set PRAGMA for sqlite", err)
	}
}

func New(driverName string, dataSourceName string) Storage {
	log.Printf("Making connection to data storage, driver=%s datasource=%s", driverName, dataSourceName)
	connections, err := sql.Open(driverName, dataSourceName)

	if err != nil {
		log.Fatal("Can not connect to data storage", err)
	}

	if driverName == "sqlite3" {
		enableForeignKeys(connections)
	}

	return Storage{connections, driverName}
}

func (storage Storage) Close() {
	log.Println("Closing connection to data storage")
	if storage.connections != nil {
		err := storage.connections.Close()
		if err != nil {
			log.Fatal("Can not close connection to data storage", err)
		}
	}
}

// Representation of cluster record in the controller service.
//     ID: unique key
//     Name: cluster GUID in the following format:
//         c8590f31-e97e-4b85-b506-c45ce1911a12
type Cluster struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

// Representation of configuration profile record in the controller service.
//     ID: unique key
//     Configuration: a JSON structure stored in a string
//     ChangeAt: username of admin that created or updated the configuration
//     ChangeBy: timestamp of the last configuration change
//     Description: a string with any comment(s) about the configuration
type ConfigurationProfile struct {
	Id            int    `json:"id"`
	Configuration string `json:"configuration"`
	ChangedAt     string `json:"changed_at"`
	ChangedBy     string `json:"changed_by"`
	Description   string `json:"description"`
}

// Representation of cluster configuration record in the controller service.
//     ID: unique key
//     Cluster: cluster ID (not name)
//     Configuration: a JSON structure stored in a string
//     ChangeAt: timestamp of the last configuration change
//     ChangeBy: username of admin that created or updated the configuration
//     Active: flag indicating whether the configuration is active or not
//     Reason: a string with any comment(s) about the cluster configuration
type ClusterConfiguration struct {
	Id            int    `json:"id"`
	Cluster       string `json:"cluster"`
	Configuration string `json:"configuration"`
	ChangedAt     string `json:"changed_at"`
	ChangedBy     string `json:"changed_by"`
	Active        string `json:"active"`
	Reason        string `json:"reason"`
}

// Representation of trigger record in the controller service
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
	Id          int    `json:"id"`
	Type        string `json:"type"`
	Cluster     string `json:"cluster"`
	Reason      string `json:"reason"`
	Link        string `json:"link"`
	TriggeredAt string `json:"triggered_at"`
	TriggeredBy string `json:"triggered_by"`
	AckedAt     string `json:"acked_at"`
	Parameters  string `json:"parameters"`
	Active      int    `json:"active"`
}

func (storage Storage) ListOfClusters() ([]Cluster, error) {
	clusters := []Cluster{}

	rows, err := storage.connections.Query("SELECT id, name FROM cluster")
	if err != nil {
		return clusters, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name string

		err = rows.Scan(&id, &name)
		if err == nil {
			clusters = append(clusters, Cluster{id, name})
		} else {
			log.Println("error", err)
		}
	}
	return clusters, nil
}

func (storage Storage) GetCluster(id int) (Cluster, error) {
	var cluster Cluster

	rows, err := storage.connections.Query("SELECT id, name FROM cluster WHERE id = ?", id)
	if err != nil {
		return cluster, err
	}
	defer rows.Close()

	if rows.Next() {
		var id int
		var name string

		err = rows.Scan(&id, &name)
		if err == nil {
			cluster.Id = id
			cluster.Name = name
		} else {
			log.Println("error", err)
		}
	} else {
		return cluster, errors.New("Unknown cluster ID provided")
	}
	return cluster, err
}

func (storage Storage) RegisterNewCluster(name string) error {
	statement, err := storage.connections.Prepare("INSERT INTO cluster(name) VALUES (?)")
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(name)
	return err
}

func (storage Storage) CreateNewCluster(id string, name string) error {
	statement, err := storage.connections.Prepare("INSERT INTO cluster(id, name) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(id, name)
	return err
}

func (storage Storage) DeleteCluster(id string) error {
	statement, err := storage.connections.Prepare("DELETE FROM cluster WHERE id=?")
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(id)
	return err
}

func (storage Storage) GetClusterByName(name string) (Cluster, error) {
	var cluster Cluster

	rows, err := storage.connections.Query("SELECT id, name FROM cluster WHERE name = ?", name)
	if err != nil {
		return cluster, err
	}
	defer rows.Close()

	if rows.Next() {
		var id int
		var name string

		err = rows.Scan(&id, &name)
		if err == nil {
			cluster.Id = id
			cluster.Name = name
			log.Printf("Cluster name %s has id %d\n", name, id)
		} else {
			log.Println("error", err)
		}
	} else {
		return cluster, errors.New("Unknown cluster ID provided")
	}
	return cluster, err
}

func (storage Storage) ListConfigurationProfiles() ([]ConfigurationProfile, error) {
	profiles := []ConfigurationProfile{}

	rows, err := storage.connections.Query("SELECT id, configuration, changed_at, changed_by, description FROM configuration_profile")
	if err != nil {
		return profiles, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var configuration string
		var changed_at string
		var changed_by string
		var description string

		err = rows.Scan(&id, &configuration, &changed_at, &changed_by, &description)
		if err == nil {
			profiles = append(profiles, ConfigurationProfile{id, configuration, changed_at, changed_by, description})
		} else {
			log.Println("error", err)
		}
	}

	return profiles, nil
}

func (storage Storage) GetConfigurationProfile(id int) (ConfigurationProfile, error) {
	var profile ConfigurationProfile

	rows, err := storage.connections.Query("SELECT id, configuration, changed_at, changed_by, description FROM configuration_profile WHERE id = ?", id)
	if err != nil {
		return profile, err
	}
	defer rows.Close()

	if rows.Next() {
		var id int
		var configuration string
		var changed_at string
		var changed_by string
		var description string

		err = rows.Scan(&id, &configuration, &changed_at, &changed_by, &description)
		if err == nil {
			profile.Id = id
			profile.Configuration = configuration
			profile.ChangedAt = changed_at
			profile.ChangedBy = changed_by
			profile.Description = description
		} else {
			log.Println("error", err)
		}
	} else {
		return profile, errors.New("Unknown configuration profile ID provided")
	}
	return profile, err
}

func (storage Storage) StoreConfigurationProfile(username string, description string, configuration string) ([]ConfigurationProfile, error) {
	var profiles []ConfigurationProfile

	t := time.Now()

	statement, err := storage.connections.Prepare("INSERT INTO configuration_profile(configuration, changed_at, changed_by, description) VALUES (?, ?, ?, ?)")
	if err != nil {
		return profiles, err
	}
	defer statement.Close()

	_, err = statement.Exec(configuration, t, username, description)
	if err != nil {
		return profiles, err
	}

	return storage.ListConfigurationProfiles()
}

func (storage Storage) ChangeConfigurationProfile(id int, username string, description string, configuration string) ([]ConfigurationProfile, error) {
	var profiles []ConfigurationProfile

	t := time.Now()

	statement, err := storage.connections.Prepare("UPDATE configuration_profile SET configuration=?, changed_at=?, changed_by=?, description=? WHERE id=?")
	if err != nil {
		return profiles, err
	}
	defer statement.Close()

	_, err = statement.Exec(configuration, t, username, description, id)
	if err != nil {
		return profiles, err
	}

	return storage.ListConfigurationProfiles()
}

func (storage Storage) DeleteConfigurationProfile(id int) ([]ConfigurationProfile, error) {
	var profiles []ConfigurationProfile

	statement, err := storage.connections.Prepare("DELETE FROM configuration_profile WHERE id=?")
	if err != nil {
		return profiles, err
	}
	defer statement.Close()

	_, err = statement.Exec(id)
	if err != nil {
		return profiles, err
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
		var changed_at string
		var changed_by string
		var active string
		var reason string

		err := rows.Scan(&id, &cluster, &configuration, &changed_at, &changed_by, &active, &reason)
		if err == nil {
			configurations = append(configurations, ClusterConfiguration{id, cluster, configuration, changed_at, changed_by, active, reason})
		} else {
			log.Println("error", err)
		}
	}

	return configurations, nil
}

func (storage Storage) ListAllClusterConfigurations() ([]ClusterConfiguration, error) {
	rows, err := storage.connections.Query(`
SELECT operator_configuration.id, cluster.name, configuration, changed_at, changed_by, active, reason
  FROM operator_configuration, cluster
    ON cluster.id = operator_configuration.cluster`)

	if err != nil {
		return []ClusterConfiguration{}, err
	}
	return storage.readClusterConfigurations(rows)
}

func (storage Storage) ListClusterConfiguration(cluster string) ([]ClusterConfiguration, error) {
	rows, err := storage.connections.Query(`
SELECT operator_configuration.id, cluster.name, configuration, changed_at, changed_by, active, reason
  FROM operator_configuration, cluster
    ON cluster.id = operator_configuration.cluster
 WHERE cluster.name=?`, cluster)

	if err != nil {
		return []ClusterConfiguration{}, err
	}

	return storage.readClusterConfigurations(rows)
}

func (storage Storage) GetClusterConfigurationById(id string) (string, error) {
	var configuration string

	row, err := storage.connections.Query(`
SELECT configuration_profile.configuration
  FROM operator_configuration, configuration_profile
    ON configuration_profile.id = operator_configuration.configuration
 WHERE operator_configuration.id=?`, id)

	if err != nil {
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
	return configuration, errors.New("unable to read any active configuration")
}

func (storage Storage) GetClusterActiveConfiguration(cluster string) (string, error) {
	var configuration string

	row, err := storage.connections.Query(`
SELECT configuration_profile.configuration
  FROM operator_configuration, cluster, configuration_profile
    ON cluster.id = operator_configuration.cluster
   AND configuration_profile.id = operator_configuration.configuration
 WHERE operator_configuration.active = '1' AND cluster.name=?
 LIMIT 1`, cluster)

	if err != nil {
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
	return configuration, errors.New("unable to read any active configuration")
}

func (storage Storage) GetConfigurationIdForCluster(cluster string) (int, error) {
	rows, err := storage.connections.Query(`
SELECT operator_configuration.id
  FROM operator_configuration, cluster
    ON cluster.id = operator_configuration.cluster
 WHERE cluster.name=?`, cluster)

	if err != nil {
		return 0, err
	}
	defer rows.Close()

	if rows.Next() {
		var id int

		err = rows.Scan(&id)
		return id, err
	} else {
		return 0, errors.New("Unknown operator name provided")
	}
}

func (storage Storage) InsertNewConfigurationProfile(tx *sql.Tx, configuration string, username string, description string) bool {
	t := time.Now()

	statement, err := tx.Prepare("INSERT INTO configuration_profile(configuration, changed_at, changed_by, description) VALUES (?, ?, ?, ?)")
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

func (storage Storage) SelectConfigurationProfileId(tx *sql.Tx) (int, error) {
	rows, err := tx.Query(`SELECT rowid FROM configuration_profile ORDER BY rowid DESC limit 1`)
	if err != nil {
		return -1, err
	}
	defer rows.Close()

	if rows.Next() {
		var configurationId int
		err = rows.Scan(&configurationId)
		if err != nil {
			return -1, err
		}
		log.Printf("Configuration stored under ID=%d\n", configurationId)
		return configurationId, nil
	} else {
		return -1, errors.New("can not retrieve last configuration ID")
	}
}

func (storage Storage) DeactivatePreviousConfigurations(tx *sql.Tx, clusterId int) error {
	stmt, err := tx.Prepare("UPDATE operator_configuration SET active=0 WHERE cluster=?")
	defer stmt.Close()

	if err != nil {
		return err
	}
	_, err = stmt.Exec(clusterId)
	if err == nil {
		log.Printf("All previous configuration has been deactivated for clusterID %d\n", clusterId)
	}
	return err
}

func (storage Storage) InsertNewOperatorConfiguration(tx *sql.Tx, clusterId int, configurationId int, username string, reason string) error {
	t := time.Now()
	statement, err := tx.Prepare("INSERT INTO operator_configuration(cluster, configuration, changed_at, changed_by, active, reason) VALUES (?, ?, ?, ?, ?, ?)")
	defer statement.Close()
	if err != nil {
		return err
	}

	_, err = statement.Exec(clusterId, configurationId, t, username, "1", reason)
	if err == nil {
		log.Printf("New operator configuration %d has been assigned to cluster %d\n", configurationId, clusterId)
	}
	return err
}

func (storage Storage) CreateClusterConfiguration(cluster string, username string, reason string, description, configuration string) ([]ClusterConfiguration, error) {
	// retrieve cluster ID
	clusterInfo, err := storage.GetClusterByName(cluster)

	if err != nil {
		return []ClusterConfiguration{}, err
	}

	clusterId := clusterInfo.Id

	// begin transaction
	tx, err := storage.connections.Begin()
	if err != nil {
		log.Println("Transaction failed")
		return []ClusterConfiguration{}, err
	}

	// insert new configuration profile
	if !storage.InsertNewConfigurationProfile(tx, configuration, username, description) {
		_ = tx.Rollback()
		return []ClusterConfiguration{}, err
	}

	// retrieve configuration ID for newly created configuration
	configurationId, err := storage.SelectConfigurationProfileId(tx)
	if err != nil {
		_ = tx.Rollback()
		return []ClusterConfiguration{}, err
	}

	// deactivate all previous configurations
	err = storage.DeactivatePreviousConfigurations(tx, clusterId)
	if err != nil {
		_ = tx.Rollback()
		return []ClusterConfiguration{}, err
	}

	// and insert new one that will be activated
	err = storage.InsertNewOperatorConfiguration(tx, clusterId, configurationId, username, reason)
	if err != nil {
		_ = tx.Rollback()
		return []ClusterConfiguration{}, err
	}

	// end the transaction
	if err := tx.Commit(); err != nil {
		return []ClusterConfiguration{}, err
	}

	return storage.ListClusterConfiguration(cluster)
}

func (storage Storage) EnableClusterConfiguration(cluster string, username string, reason string) ([]ClusterConfiguration, error) {
	id, err := storage.GetConfigurationIdForCluster(cluster)
	if err != nil {
		return []ClusterConfiguration{}, err
	}

	statement, err := storage.connections.Prepare("UPDATE operator_configuration SET active=1, changed_at=?, changed_by=?, reason=? WHERE id=?")
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

// TODO: copy & paste, needs to be refactored later
func (storage Storage) DisableClusterConfiguration(cluster string, username string, reason string) ([]ClusterConfiguration, error) {
	id, err := storage.GetConfigurationIdForCluster(cluster)
	if err != nil {
		return []ClusterConfiguration{}, err
	}
	statement, err := storage.connections.Prepare("UPDATE operator_configuration SET active=0, changed_at=?, changed_by=?, reason=? WHERE id=?")
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

func (storage Storage) EnableOrDisableClusterConfigurationById(id string, active string) error {
	statement, err := storage.connections.Prepare("UPDATE operator_configuration SET active=?, changed_at=? WHERE id=?")
	if err != nil {
		return err
	}
	defer statement.Close()

	t := time.Now()

	_, err = statement.Exec(active, t, id)
	if err != nil {
		return err
	}
	return nil
}

// TODO: copy & paste, needs to be refactored later
func (storage Storage) DeleteClusterConfigurationById(id string) error {
	statement, err := storage.connections.Prepare("DELETE FROM operator_configuration WHERE id=?")
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(id)
	if err != nil {
		return err
	}
	return nil
}

func (storage Storage) getTriggers(rows *sql.Rows) ([]Trigger, error) {
	triggers := []Trigger{}
	defer rows.Close()

	for rows.Next() {
		var trigger Trigger

		err := rows.Scan(&trigger.Id, &trigger.Type, &trigger.Cluster,
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

func (storage Storage) GetTriggerById(id string) (Trigger, error) {
	rows, err := storage.connections.Query(`
SELECT trigger.id, trigger_type.type, cluster.name,
       trigger.reason, trigger.link, trigger.triggered_at, trigger.triggered_by,
       trigger.parameters, trigger.active, trigger.acked_at
  FROM trigger JOIN trigger_type ON trigger.type=trigger_type.id
               JOIN cluster ON trigger.cluster=cluster.id
 WHERE trigger.id=?`, id)

	if err != nil {
		return Trigger{}, err
	}

	triggers, err := storage.getTriggers(rows)
	if err != nil {
		return Trigger{}, err
	}

	if len(triggers) >= 1 {
		return triggers[0], nil
	} else {
		return Trigger{}, fmt.Errorf("No such trigger for ID=%s", id)
	}
}

func (storage Storage) DeleteTriggerById(id string) error {
	statement, err := storage.connections.Prepare(`
DELETE FROM trigger WHERE trigger.id=?`)
	if err != nil {
		return err
	}

	defer statement.Close()

	_, err = statement.Exec(id)
	return err
}

func (storage Storage) ChangeStateOfTriggerById(id string, active int) error {
	statement, err := storage.connections.Prepare(`
UPDATE trigger SET active=? WHERE trigger.id=?`)
	if err != nil {
		return err
	}

	defer statement.Close()

	_, err = statement.Exec(active, id)
	return err
}

func (storage Storage) ListAllTriggers() ([]Trigger, error) {
	triggers := []Trigger{}

	rows, err := storage.connections.Query(`
SELECT trigger.id, trigger_type.type, cluster.name,
       trigger.reason, trigger.link, trigger.triggered_at, trigger.triggered_by,
       trigger.parameters, trigger.active, trigger.acked_at
  FROM trigger JOIN trigger_type ON trigger.type=trigger_type.id
               JOIN cluster ON trigger.cluster=cluster.id`)

	if err != nil {
		return triggers, err
	}

	return storage.getTriggers(rows)
}

func (storage Storage) ListClusterTriggers(clusterName string) ([]Trigger, error) {
	triggers := []Trigger{}

	rows, err := storage.connections.Query(`
SELECT trigger.id, trigger_type.type, cluster.name,
       trigger.reason, trigger.link, trigger.triggered_at, trigger.triggered_by,
       trigger.parameters, trigger.active, trigger.acked_at
  FROM trigger JOIN trigger_type ON trigger.type=trigger_type.id
               JOIN cluster ON trigger.cluster=cluster.id
 WHERE cluster.name = ?`, clusterName)

	if err != nil {
		return triggers, err
	}

	return storage.getTriggers(rows)
}

func (storage Storage) ListActiveClusterTriggers(clusterName string) ([]Trigger, error) {
	triggers := []Trigger{}

	rows, err := storage.connections.Query(`
SELECT trigger.id, trigger_type.type, cluster.name,
       trigger.reason, trigger.link, trigger.triggered_at, trigger.triggered_by,
       trigger.parameters, trigger.active, trigger.acked_at
  FROM trigger JOIN trigger_type ON trigger.type=trigger_type.id
               JOIN cluster ON trigger.cluster=cluster.id
 WHERE trigger.active = 1
   AND cluster.name = ?`, clusterName)

	if err != nil {
		return triggers, err
	}

	return storage.getTriggers(rows)
}

func (storage Storage) GetTriggerId(triggerType string) (int, error) {
	var id int

	rows, err := storage.connections.Query("SELECT id FROM trigger_type WHERE type = ?", triggerType)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

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

func (storage Storage) NewTrigger(clusterName string, triggerType string, userName string, reason string, link string) error {
	// retrieve cluster ID
	clusterInfo, err := storage.GetClusterByName(clusterName)
	clusterId := clusterInfo.Id

	if err != nil {
		return err
	}

	triggerId, err := storage.GetTriggerId(triggerType)

	if err != nil {
		return err
	}
	t := time.Now()

	statement, err := storage.connections.Prepare("INSERT INTO trigger(type, cluster, reason, link, triggered_at, triggered_by, parameters, active, acked_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, '')")
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(triggerId, clusterId, reason, link, t, userName, "", 1)
	if err != nil {
		return err
	}
	return nil
}

func (storage Storage) AckTrigger(clusterName string, triggerId string) error {
	t := time.Now()

	// retrieve cluster ID
	clusterInfo, err := storage.GetClusterByName(clusterName)
	clusterId := clusterInfo.Id

	if err != nil {
		return err
	}

	statement, err := storage.connections.Prepare("UPDATE trigger SET acked_at=?, active=0 WHERE cluster = ? AND id = ?")

	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(t, clusterId, triggerId)
	if err != nil {
		return err
	}
	return nil
}
