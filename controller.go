package main

import (
	"fmt"
	"github.com/redhatinsighs/insights-operator-controller/server"
	"github.com/redhatinsighs/insights-operator-controller/storage"
	"github.com/spf13/viper"
)

// Entry point to the Insights operator controller.
// It performs several tasks:
// - connect to the storage with basic test if storage is accessible
// - start the HTTP server with all required endpints
// - TODO: initialize connection to the logging service
func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	storage := storage.New("sqlite3", "./controller.db")
	defer storage.Close()

	server.Initialize(":8080", storage)
}
