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
package main

import (
	"flag"
	"fmt"
	"github.com/redhatinsighs/insights-operator-controller/logging"
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
	// parse the configuration
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	// parse all command-line arguments
	dbDriver := flag.String("dbdriver", "sqlite3", "database driver specification")
	storageSpecification := flag.String("storage", "./controller.db", "storage specification")
	flag.Parse()

	storage := storage.New(*dbDriver, *storageSpecification)
	defer storage.Close()

	splunkCfg := viper.Sub("splunk")
	address := splunkCfg.GetString("address")
	token := splunkCfg.GetString("token")
	source := splunkCfg.GetString("source")
	source_type := splunkCfg.GetString("source_type")
	index := splunkCfg.GetString("index")
	splunk := logging.NewClient(address, token, source, source_type, index)

	server.Initialize(":8080", storage)
}
