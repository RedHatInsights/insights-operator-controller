package storage

import (
	"database/sql"
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

func ListOfClusters() {
	rows, err := connections.Query("SELECT id, cluster FROM cluster")

	var uid int
	var cluster string

	for rows.Next() {
		err = rows.Scan(&uid, &cluster)
		if err == nil {
			fmt.Println(uid)
			fmt.Println(cluster)
		} else {
			log.Println("error", err)
		}
	}

	rows.Close()
}
