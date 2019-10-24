package storage

import (
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"time"
)

type Storage struct {
	connections *sql.DB
}

func New(driverName string, dataSourceName string) Storage {
	log.Printf("Making connection to data storage, driver=%s datasource=%s", driverName, dataSourceName)
	connections, err := sql.Open(driverName, dataSourceName)

	if err != nil {
		log.Fatal("Can not connect to data storage", err)
	}
	return Storage{connections}
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

type Cluster struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type ConfigurationProfile struct {
	Id            int    `json:"id"`
	Configuration string `json:"configuration"`
	ChangedAt     string `json:"changed_at"`
	ChangedBy     string `json:"changed_by"`
	Description   string `json:"description"`
}

type ClusterConfiguration struct {
	Id            int    `json:"id"`
	Cluster       string `json:"cluster"`
	Configuration string `json:"configuration"`
	ChangedAt     string `json:"changed_at"`
	ChangedBy     string `json:"changed_by"`
	Active        string `json:"active"`
	Reason        string `json:"reason"`
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
	clusterId := clusterInfo.Id

	if err != nil {
		return []ClusterConfiguration{}, err
	}

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
