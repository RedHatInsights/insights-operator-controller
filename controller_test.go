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

//
// Documentation in literate-programming-style is available at:
// https://redhatinsights.github.io/insights-operator-controller/packages/controller_test.html
package main_test

import (
	"os"
	"testing"

	main "github.com/RedHatInsights/insights-operator-controller"
)

const (
	ConfigFileEnvironmentVariable = "TEST_CONFIG_FILE"
	ExistingConfigFile            = "./config_tests.toml"
	IncorrectConfigFile           = "foobar"
)

func mustSetEnv(t *testing.T, key, val string) {
	err := os.Setenv(key, val)
	if err != nil {
		t.Fatal(err)
	}
}

// TestInitializeSplunkEnabledClient check whether the Splunk logging client can be initialized
func TestInitializeSplunkEnabledClient(t *testing.T) {
	// configuration don't have to be set fully, just Splunk part
	cfg := main.Configuration{}
	cfg.SplunkEnabled = true
	cfg.SplunkAddress = "address"
	cfg.SplunkToken = "token"
	cfg.SplunkSource = "source"
	cfg.SplunkSourceType = "source_type"
	cfg.SplunkIndex = "index"

	client := main.InitializeSplunk(&cfg)
	if client.ClientImpl == nil {
		t.Fatal("Splunk logging client has not been initialized")
	}
}

// TestInitializeSplunkDisabledClient check whether the Splunk logging client can be initialized
func TestInitializeSplunkDisabledClient(t *testing.T) {
	// configuration don't have to be set fully, just Splunk part
	cfg := main.Configuration{}
	cfg.SplunkEnabled = false
	cfg.SplunkAddress = "address"
	cfg.SplunkToken = "token"
	cfg.SplunkSource = "source"
	cfg.SplunkSourceType = "source_type"
	cfg.SplunkIndex = "index"

	client := main.InitializeSplunk(&cfg)
	if client.ClientImpl != nil {
		t.Fatal("Splunk logging client should not be initialized")
	}
}

// readConfigFileSpecifiedByEnvVar tries to read configuration file specified by environment variable
func readConfigFileSpecifiedByEnvVar(t *testing.T, filename string) error {
	mustSetEnv(t, ConfigFileEnvironmentVariable, filename)
	return main.ReadConfigurationFile(ConfigFileEnvironmentVariable)
}

// TestReadConfigurationFileViaEnvVariable checks whether it is possible to read configuration file specified by env. variable
func TestReadConfigurationFileViaEnvVariable(t *testing.T) {
	err := readConfigFileSpecifiedByEnvVar(t, ExistingConfigFile)
	if err != nil {
		t.Fatal("Error during config file reading", err)
	}

	// default config file should be read w/o any error
	os.Clearenv()
	err = main.ReadConfigurationFile(ConfigFileEnvironmentVariable)
	if err != nil {
		t.Fatal("Error during config file reading", err)
	}
}

// TestReadConfigurationFileNegative checks whether it is possible to read configuration file specified by env. variable
func TestReadConfigurationFileNegative(t *testing.T) {
	err := readConfigFileSpecifiedByEnvVar(t, IncorrectConfigFile)
	if err == nil {
		t.Fatal("Non-existing config file should not be processed w/o error")
	}
}

// TestReadConfigurationFromEnvVar check the ability to read configuration from file specified in environment variable
func TestReadConfigurationFromEnvVar(t *testing.T) {
	mustSetEnv(t, ConfigFileEnvironmentVariable, ExistingConfigFile)
	cfg, err := main.ReadConfiguration(ConfigFileEnvironmentVariable)
	if err != nil {
		t.Fatal("Error during config file reading", err)
	}
	if len(cfg.Address) == 0 {
		t.Fatal("The config is probably wrong", err)
	}
}

// TestReadConfigurationFromDefaultFile check the ability to read configuration from default config file
func TestReadConfigurationFromDefaultFile(t *testing.T) {
	mustSetEnv(t, ConfigFileEnvironmentVariable, IncorrectConfigFile)
	_, err := main.ReadConfiguration(ConfigFileEnvironmentVariable)
	if err == nil {
		t.Fatal("Error is expected for non existing configuration file")
	}
}
