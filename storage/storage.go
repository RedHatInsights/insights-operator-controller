package storage

import (
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

type Storage struct {
	connections *sql.DB
}

func New(driverName string, dataSourceName string) Storage {
	log.Println("Making connection to data storage")
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

func (storage Storage) ListOfClusters() []Cluster {
	clusters := []Cluster{}

	rows, err := storage.connections.Query("SELECT id, name FROM cluster")
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
	rows.Close()

	return clusters
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
	rows.Close()
	return cluster, err
}

func (storage Storage) ListConfigurationProfiles() []ConfigurationProfile {
	profiles := []ConfigurationProfile{}

	rows, err := storage.connections.Query("SELECT id, configuration, changed_at, changed_by, description FROM configuration_profile")
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

	return profiles
}

func (storage Storage) ListClusterConfiguration(cluster string) []ClusterConfiguration {
	configurations := []ClusterConfiguration{}

	rows, err := storage.connections.Query(`
SELECT operator_configuration.id, cluster.name, configuration, changed_at, changed_by, active, reason
  FROM operator_configuration, cluster
    ON cluster.id = operator_configuration.cluster
 WHERE cluster.name=?`, cluster)

	defer rows.Close()

	for rows.Next() {
		var id int
		var cluster string
		var configuration string
		var changed_at string
		var changed_by string
		var active string
		var reason string

		err = rows.Scan(&id, &cluster, &configuration, &changed_at, &changed_by, &active, &reason)
		if err == nil {
			configurations = append(configurations, ClusterConfiguration{id, cluster, configuration, changed_at, changed_by, active, reason})
		} else {
			log.Println("error", err)
		}
	}
	rows.Close()

	return configurations
}
