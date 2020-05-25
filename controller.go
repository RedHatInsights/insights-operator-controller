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
	"github.com/RedHatInsights/insights-operator-controller/logging"
	"github.com/RedHatInsights/insights-operator-controller/server"
	"github.com/RedHatInsights/insights-operator-controller/storage"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
)

// ConfigurationEnvVarName contains name of environment variable with configuration file name settiongs
const ConfigurationEnvVarName = "INSIGHTS_CONTROLLER_CONFIG_FILE"

// Configuration represents service configuration
type Configuration struct {
	UseHTTPS             bool
	Address              string
	TlsCert              string
	TlsKey               string
	DbDriver             string
	StorageSpecification string
}

func initializeSplunk() logging.Client {
	splunkCfg := viper.Sub("splunk")
	enabled := splunkCfg.GetBool("enabled")
	address := splunkCfg.GetString("address")
	token := splunkCfg.GetString("token")
	source := splunkCfg.GetString("source")
	sourceType := splunkCfg.GetString("source_type")
	index := splunkCfg.GetString("index")
	return logging.NewClient(enabled, address, token, source, sourceType, index)
}

func readConfigurationFile(envVar string) error {
	configFile, specified := os.LookupEnv(envVar)
	if specified {
		// we need to separate the directory name and filename without extension
		directory, basename := filepath.Split(configFile)
		file := strings.TrimSuffix(basename, filepath.Ext(basename))
		// parse the configuration
		viper.SetConfigName(file)
		viper.AddConfigPath(directory)
	} else {
		// parse the configuration
		viper.SetConfigName("config")
		viper.AddConfigPath(".")
	}

	err := viper.ReadInConfig()
	return err
}

func readConfiguration(envVar string) (Configuration, error) {
	cfg := Configuration{}

	err := readConfigurationFile(envVar)
	if err != nil {
		return cfg, err
	}

	serviceCfg := viper.Sub("service")
	cfg.UseHTTPS = serviceCfg.GetBool("use_https")
	cfg.Address = serviceCfg.GetString("address")
	cfg.TlsCert = serviceCfg.GetString("tls_cert")
	cfg.TlsKey = serviceCfg.GetString("tls_key")

	// parse all command-line arguments
	dbDriver := flag.String("dbdriver", "sqlite3", "database driver specification")
	storageSpecification := flag.String("storage", "./controller.db", "storage specification")
	flag.Parse()

	cfg.DbDriver = *dbDriver
	cfg.StorageSpecification = *storageSpecification

	return cfg, nil
}

// Entry point to the Insights operator controller.
// It performs several tasks:
// - connect to the storage with basic test if storage is accessible
// - start the HTTP server with all required endpints
// - TODO: initialize connection to the logging service
func main() {
	cfg, err := readConfiguration(ConfigurationEnvVarName)
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s", err))
	}

	storage := storage.New(cfg.DbDriver, cfg.StorageSpecification)
	defer storage.Close()

	splunk := initializeSplunk()

	s := server.Server{
		Address:  cfg.Address,
		UseHTTPS: cfg.UseHTTPS,
		Storage:  storage,
		Splunk:   splunk,
		TLSCert:  cfg.TlsCert,
		TLSKey:   cfg.TlsKey,
	}

	s.Initialize()
}
