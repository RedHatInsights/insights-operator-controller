package storage

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var connections *sql.DB

func Initialize(driverName string, dataSourceName string) {
	var err error

	log.Println("Making connection to data storage")
	connections, err = sql.Open(driverName, dataSourceName)

	if err != nil {
		log.Fatal("Can not connect to data storage", err)
	}
}

func Close() {
	if connections != nil {
		err := connections.Close()
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
