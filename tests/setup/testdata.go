// Copyright 2020 Red Hat, Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package testdata provides utilities for inserting dummy data for
// testing the insight-operator-controller .
package testdata

import (
	"fmt"

	"github.com/RedHatInsights/insights-operator-controller/storage"
	"github.com/brianvoe/gofakeit"
	"github.com/spf13/viper"
)

// DataGenerator is a wrapper for the Storage type providing methods to
// insert dummy data into the database, instantiate with NewFactory(..)
type DataGenerator struct {
	storage storage.Storage
	config  dataConfiguration
}

// DataConfiguration contains fields for test setup
type dataConfiguration struct {
	OperatorConfigurationNo int
	ClusterNo               int
	TriggerNo               int
	ConfProfileNo           int
}

// NewDataGenerator provides an instance of Factory type, allowing access to
// methods for creating dummy data for tests
// Insert the same dbDriver and storageSpecification you expect to use
// in your test to have meaningful and coherent tests
func NewDataGenerator(dbDriver string, storageSpecification string) DataGenerator {
	var returnMe DataGenerator
	returnMe.storage = storage.New(dbDriver, storageSpecification)

	viper.SetConfigName("testconfig")

	viper.AddConfigPath("../setup")
	viper.AddConfigPath("./setup")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s ", err))
	}

	conf := viper.Sub("benchmark")
	returnMe.config = dataConfiguration{
		OperatorConfigurationNo: conf.GetInt("operator_configuration_no"),
		ClusterNo:               conf.GetInt("cluster_no"),
		TriggerNo:               conf.GetInt("trigger_no"),
		ConfProfileNo:           conf.GetInt("conf_profile_no"),
	}
	return returnMe
}

// Close closes a connection to the DB created to insert the test data
func (g DataGenerator) Close() {
	g.storage.Close()
}

//PopulateCluster inserts clusterNo cluster objects in the database
func (g DataGenerator) PopulateCluster() []error {

	var err error
	var errs []error

	for i := 0; i < g.config.ClusterNo; i++ {

		err = g.storage.RegisterNewCluster(GetClusterName(i))
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}

// GetClusterName returns a string in the format
// "00000000-0000-0000-0000-num" where num is padded with trailing zeros
// until it has a length of 12 char
func GetClusterName(num int) string {
	prefix := "00000000-0000-0000-0000-"
	end := "000000000000"
	clusterName := fmt.Sprintf("%s%d", end, num)
	clusterName = clusterName[len(clusterName)-len(end):]
	return fmt.Sprintf("%s%s", prefix, clusterName)
}

// PopulateConfigurationProfile inserts configurationProfileNo
// configuration_profile objects into the database
func (g DataGenerator) PopulateConfigurationProfile() []error {
	gofakeit.Seed(0)
	var err error
	var errs []error

	confStr := "{\"no_op\":\"X\", \"watch\":[\"a\",\"b\",\"c\"]}"
	for i := 0; i < g.config.ConfProfileNo; i++ {
		_, err = g.storage.StoreConfigurationProfile(gofakeit.Username(), gofakeit.Sentence(1), confStr)
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

// PopulateOperatorConfiguration inserts operatorConfigurationNo
// operator_configuration objects into the database
func (g DataGenerator) PopulateOperatorConfiguration() []error {
	var errs []error
	gofakeit.Seed(0)
	clusters, queryErr := g.storage.ListOfClusters()
	if queryErr != nil {
		errs = append(errs, queryErr)
		return errs
	}

	for i := 0; i < g.config.OperatorConfigurationNo; i++ {
		//creates configuration profiles and operator Configuration
		_, err := g.storage.CreateClusterConfiguration(
			string(clusters[i%len(clusters)].Name),
			gofakeit.Username(),
			gofakeit.Sentence(1),
			gofakeit.Sentence(3),
			gofakeit.Name())

		if err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}

// PopulateTrigger inserts triggerNo
// trigger objects into the database
func (g DataGenerator) PopulateTrigger(triggerType string) []error {
	var errs []error
	gofakeit.Seed(0)
	clusters, queryErr := g.storage.ListOfClusters()
	if queryErr != nil {
		errs = append(errs, queryErr)
		return errs
	}

	for i := 0; i < g.config.TriggerNo; i++ {
		err := g.storage.NewTrigger(
			string(clusters[i%len(clusters)].Name),
			triggerType,
			gofakeit.Username(),
			gofakeit.Sentence(2),
			gofakeit.URL())
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}

// InsertTriggerType inserts one trigger_type object with ttype type and
// description
func (g DataGenerator) InsertTriggerType(ttype string, description string) error {
	return g.storage.NewTriggerType(ttype, description)
}
