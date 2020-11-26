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

// Entry point to the Insights Controller service
package main

// Generated documentation is available at:
// https://godoc.org/github.com/RedHatInsights/insights-operator-controller
//
// Documentation in literate-programming-style is available at:
// https://redhatinsights.github.io/insights-operator-controller/packages/controller.html

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"

	"github.com/RedHatInsights/insights-operator-controller/logging"
	"github.com/RedHatInsights/insights-operator-controller/server"
	"github.com/RedHatInsights/insights-operator-controller/storage"
)

// ConfigurationEnvVarName contains name of environment variable with configuration file name settiongs
const ConfigurationEnvVarName = "INSIGHTS_CONTROLLER_CONFIG_FILE"

// Configuration represents service configuration
type Configuration struct {
	UseHTTPS             bool
	Address              string
	TLSCert              string
	TLSKey               string
	DbDriver             string
	StorageSpecification string
	SplunkEnabled        bool
	SplunkAddress        string
	SplunkToken          string
	SplunkSource         string
	SplunkSourceType     string
	SplunkIndex          string
}

func initializeSplunk(cfg *Configuration) logging.Client {
	return logging.NewClient(cfg.SplunkEnabled,
		cfg.SplunkAddress,
		cfg.SplunkToken,
		cfg.SplunkSource,
		cfg.SplunkSourceType,
		cfg.SplunkIndex)
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
	cfg.TLSCert = serviceCfg.GetString("tls_cert")
	cfg.TLSKey = serviceCfg.GetString("tls_key")

	splunkCfg := viper.Sub("splunk")
	cfg.SplunkEnabled = splunkCfg.GetBool("enabled")
	cfg.SplunkAddress = splunkCfg.GetString("address")
	cfg.SplunkToken = splunkCfg.GetString("token")
	cfg.SplunkSource = splunkCfg.GetString("source")
	cfg.SplunkSourceType = splunkCfg.GetString("source_type")
	cfg.SplunkIndex = splunkCfg.GetString("index")

	storageCfg := viper.Sub("storage")
	cfg.DbDriver = storageCfg.GetString("driver")
	cfg.StorageSpecification = splunkCfg.GetString("source")

	// parse all command-line arguments
	dbDriver := flag.String("dbdriver", "sqlite3", "database driver specification")
	storageSpecification := flag.String("storage", "./controller.db", "storage specification")
	flag.Parse()

	// override configuration by CLI parameter
	if dbDriver != nil {
		cfg.DbDriver = *dbDriver
	}

	// override configuration by CLI parameter
	if storageSpecification != nil {
		cfg.StorageSpecification = *storageSpecification
	}

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

	// try to initialize the storage
	storageInstance, err := storage.New(cfg.DbDriver, cfg.StorageSpecification)
	if err != nil {
		panic(err)
	}
	defer storageInstance.Close()

	// try to check if storage is really configured properly
	err = storageInstance.Ping()
	if err != nil {
		panic(err)
	}

	splunk := initializeSplunk(&cfg)

	s := server.Server{
		Address:  cfg.Address,
		UseHTTPS: cfg.UseHTTPS,
		Storage:  storageInstance,
		Splunk:   splunk,
		TLSCert:  cfg.TLSCert,
		TLSKey:   cfg.TLSKey,
	}

	s.Initialize()
}
