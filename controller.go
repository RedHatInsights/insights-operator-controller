package main

import (
	"github.com/redhatinsighs/insights-operator-controller/server"
	"github.com/redhatinsighs/insights-operator-controller/storage"
)

func main() {
	storage := storage.New("sqlite3", "./controller.db")
	defer storage.Close()

	server.Initialize(":8080")
}
