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

package tests

import (
	"encoding/json"
	"fmt"
	"github.com/verdverm/frisby"
)

// ClusterConfiguration represents cluster configuration record in the controller service.
//     ID: unique key
//     Cluster: cluster ID (not name)
//     Configuration: a JSON structure stored in a string
//     ChangeAt: timestamp of the last configuration change
//     ChangeBy: username of admin that created or updated the configuration
//     Active: flag indicating whether the configuration is active or not
//     Reason: a string with any comment(s) about the cluster configuration
type ClusterConfiguration struct {
	ID            int    `json:"id"`
	Cluster       string `json:"cluster"`
	Configuration string `json:"configuration"`
	ChangedAt     string `json:"changed_at"`
	ChangedBy     string `json:"changed_by"`
	Active        string `json:"active"`
	Reason        string `json:"reason"`
}

func compareConfigurations(f *frisby.Frisby, configurations []ClusterConfiguration, expected []ClusterConfiguration) {
	if len(configurations) != len(expected) {
		f.AddError(fmt.Sprintf("%d configurations are expected, but got %d", len(expected), len(configurations)))
	}

	for i := 0; i < len(expected); i++ {
		// just ignore timestamp as we are going to test REST API w/o mocking the database
		configurations[i].ChangedAt = ""
		expected[i].ChangedAt = ""
		if configurations[i] != expected[i] {
			f.AddError(fmt.Sprintf("Different configuration info returned: %v != %v", configurations[i], expected[i]))
			fmt.Println(configurations[i].Configuration)
			fmt.Println(expected[i].Configuration)
		}
	}
}

func readConfigurations(f *frisby.Frisby) []ClusterConfiguration {
	f.Get(API_URL + "client/configuration")
	f.Send()
	f.ExpectStatus(200)
	f.ExpectHeader("Content-Type", "application/json; charset=utf-8")

	configurations := []ClusterConfiguration{}
	text, err := f.Resp.Content()
	if err != nil {
		f.AddError(err.Error())
	} else {
		json.Unmarshal(text, &configurations)
	}
	return configurations
}

func checkInitialListOfConfigurations() {
	f := frisby.Create("Check list of configurations")

	configurations := readConfigurations(f)
	f.ExpectJsonLength("", 6)

	expected := []ClusterConfiguration{
		{0, "00000000-0000-0000-0000-000000000000", "0", "2019-01-01T00:00:00Z", "tester", "0", "no reason"},
		{1, "00000000-0000-0000-0000-000000000000", "1", "2019-01-01T00:00:00Z", "tester", "0", "no reason"},
		{2, "00000000-0000-0000-0000-000000000000", "2", "2019-01-01T00:00:00Z", "tester", "1", "no reason"},
		{3, "00000000-0000-0000-0000-000000000001", "1", "2019-10-11T00:00:00Z", "tester", "1", "no reason"},
		{4, "00000000-0000-0000-0000-000000000002", "2", "2019-10-11T00:00:00Z", "tester", "1", "no reason so far"},
		{5, "00000000-0000-0000-0000-000000000003", "0", "2019-10-11T00:00:00Z", "tester", "0", "disabled one"},
	}
	compareConfigurations(f, configurations, expected)

	f.PrintReport()
}

func checkEnableExistingConfiguration() {
	f := frisby.Create("Check that configuration can be enabled")

	configurations := readConfigurations(f)
	f.ExpectJsonLength("", 6)

	expected := []ClusterConfiguration{
		{0, "00000000-0000-0000-0000-000000000000", "0", "2019-01-01T00:00:00Z", "tester", "0", "no reason"},
		{1, "00000000-0000-0000-0000-000000000000", "1", "2019-01-01T00:00:00Z", "tester", "0", "no reason"},
		{2, "00000000-0000-0000-0000-000000000000", "2", "2019-01-01T00:00:00Z", "tester", "1", "no reason"},
		{3, "00000000-0000-0000-0000-000000000001", "1", "2019-10-11T00:00:00Z", "tester", "1", "no reason"},
		{4, "00000000-0000-0000-0000-000000000002", "2", "2019-10-11T00:00:00Z", "tester", "1", "no reason so far"},
		{5, "00000000-0000-0000-0000-000000000003", "0", "2019-10-11T00:00:00Z", "tester", "0", "disabled one"},
	}
	compareConfigurations(f, configurations, expected)

	f.Put(API_URL + "client/configuration/0/enable")
	f.Send()
	f.ExpectStatus(200)

	configurations = readConfigurations(f)
	f.ExpectJsonLength("", 6)

	expected = []ClusterConfiguration{
		{0, "00000000-0000-0000-0000-000000000000", "0", "2019-01-01T00:00:00Z", "tester", "1", "no reason"},
		{1, "00000000-0000-0000-0000-000000000000", "1", "2019-01-01T00:00:00Z", "tester", "0", "no reason"},
		{2, "00000000-0000-0000-0000-000000000000", "2", "2019-01-01T00:00:00Z", "tester", "1", "no reason"},
		{3, "00000000-0000-0000-0000-000000000001", "1", "2019-10-11T00:00:00Z", "tester", "1", "no reason"},
		{4, "00000000-0000-0000-0000-000000000002", "2", "2019-10-11T00:00:00Z", "tester", "1", "no reason so far"},
		{5, "00000000-0000-0000-0000-000000000003", "0", "2019-10-11T00:00:00Z", "tester", "0", "disabled one"},
	}
	compareConfigurations(f, configurations, expected)

	f.PrintReport()
}

func checkDisableExistingConfiguration() {
	f := frisby.Create("Check that configuration can be disabled")

	configurations := readConfigurations(f)
	f.ExpectJsonLength("", 6)

	expected := []ClusterConfiguration{
		{0, "00000000-0000-0000-0000-000000000000", "0", "2019-01-01T00:00:00Z", "tester", "1", "no reason"},
		{1, "00000000-0000-0000-0000-000000000000", "1", "2019-01-01T00:00:00Z", "tester", "0", "no reason"},
		{2, "00000000-0000-0000-0000-000000000000", "2", "2019-01-01T00:00:00Z", "tester", "1", "no reason"},
		{3, "00000000-0000-0000-0000-000000000001", "1", "2019-10-11T00:00:00Z", "tester", "1", "no reason"},
		{4, "00000000-0000-0000-0000-000000000002", "2", "2019-10-11T00:00:00Z", "tester", "1", "no reason so far"},
		{5, "00000000-0000-0000-0000-000000000003", "0", "2019-10-11T00:00:00Z", "tester", "0", "disabled one"},
	}
	compareConfigurations(f, configurations, expected)

	f.Put(API_URL + "client/configuration/0/disable")
	f.Send()
	f.ExpectStatus(200)

	configurations = readConfigurations(f)
	f.ExpectJsonLength("", 6)

	expected = []ClusterConfiguration{
		{0, "00000000-0000-0000-0000-000000000000", "0", "2019-01-01T00:00:00Z", "tester", "0", "no reason"},
		{1, "00000000-0000-0000-0000-000000000000", "1", "2019-01-01T00:00:00Z", "tester", "0", "no reason"},
		{2, "00000000-0000-0000-0000-000000000000", "2", "2019-01-01T00:00:00Z", "tester", "1", "no reason"},
		{3, "00000000-0000-0000-0000-000000000001", "1", "2019-10-11T00:00:00Z", "tester", "1", "no reason"},
		{4, "00000000-0000-0000-0000-000000000002", "2", "2019-10-11T00:00:00Z", "tester", "1", "no reason so far"},
		{5, "00000000-0000-0000-0000-000000000003", "0", "2019-10-11T00:00:00Z", "tester", "0", "disabled one"},
	}
	compareConfigurations(f, configurations, expected)

	f.PrintReport()
}

func checkEnableNonExistingConfiguration() {
	f := frisby.Create("Check what happens when non existing configuration is enabled")

	configurations := readConfigurations(f)
	f.ExpectJsonLength("", 6)

	expected := []ClusterConfiguration{
		{0, "00000000-0000-0000-0000-000000000000", "0", "2019-01-01T00:00:00Z", "tester", "0", "no reason"},
		{1, "00000000-0000-0000-0000-000000000000", "1", "2019-01-01T00:00:00Z", "tester", "0", "no reason"},
		{2, "00000000-0000-0000-0000-000000000000", "2", "2019-01-01T00:00:00Z", "tester", "1", "no reason"},
		{3, "00000000-0000-0000-0000-000000000001", "1", "2019-10-11T00:00:00Z", "tester", "1", "no reason"},
		{4, "00000000-0000-0000-0000-000000000002", "2", "2019-10-11T00:00:00Z", "tester", "1", "no reason so far"},
		{5, "00000000-0000-0000-0000-000000000003", "0", "2019-10-11T00:00:00Z", "tester", "0", "disabled one"},
	}
	compareConfigurations(f, configurations, expected)

	f.Put(API_URL + "client/configuration/42/enable")
	f.Send()
	f.ExpectStatus(200)

	configurations = readConfigurations(f)
	f.ExpectJsonLength("", 6)

	expected = []ClusterConfiguration{
		{0, "00000000-0000-0000-0000-000000000000", "0", "2019-01-01T00:00:00Z", "tester", "0", "no reason"},
		{1, "00000000-0000-0000-0000-000000000000", "1", "2019-01-01T00:00:00Z", "tester", "0", "no reason"},
		{2, "00000000-0000-0000-0000-000000000000", "2", "2019-01-01T00:00:00Z", "tester", "1", "no reason"},
		{3, "00000000-0000-0000-0000-000000000001", "1", "2019-10-11T00:00:00Z", "tester", "1", "no reason"},
		{4, "00000000-0000-0000-0000-000000000002", "2", "2019-10-11T00:00:00Z", "tester", "1", "no reason so far"},
		{5, "00000000-0000-0000-0000-000000000003", "0", "2019-10-11T00:00:00Z", "tester", "0", "disabled one"},
	}
	compareConfigurations(f, configurations, expected)

	f.PrintReport()
}

func checkDisableNonExistingConfiguration() {
	f := frisby.Create("Check what happens when non existing configuration is disabled")

	configurations := readConfigurations(f)
	f.ExpectJsonLength("", 6)

	expected := []ClusterConfiguration{
		{0, "00000000-0000-0000-0000-000000000000", "0", "2019-01-01T00:00:00Z", "tester", "0", "no reason"},
		{1, "00000000-0000-0000-0000-000000000000", "1", "2019-01-01T00:00:00Z", "tester", "0", "no reason"},
		{2, "00000000-0000-0000-0000-000000000000", "2", "2019-01-01T00:00:00Z", "tester", "1", "no reason"},
		{3, "00000000-0000-0000-0000-000000000001", "1", "2019-10-11T00:00:00Z", "tester", "1", "no reason"},
		{4, "00000000-0000-0000-0000-000000000002", "2", "2019-10-11T00:00:00Z", "tester", "1", "no reason so far"},
		{5, "00000000-0000-0000-0000-000000000003", "0", "2019-10-11T00:00:00Z", "tester", "0", "disabled one"},
	}
	compareConfigurations(f, configurations, expected)

	f.Put(API_URL + "client/configuration/42/enable")
	f.Send()
	f.ExpectStatus(200)

	configurations = readConfigurations(f)
	f.ExpectJsonLength("", 6)

	expected = []ClusterConfiguration{
		{0, "00000000-0000-0000-0000-000000000000", "0", "2019-01-01T00:00:00Z", "tester", "0", "no reason"},
		{1, "00000000-0000-0000-0000-000000000000", "1", "2019-01-01T00:00:00Z", "tester", "0", "no reason"},
		{2, "00000000-0000-0000-0000-000000000000", "2", "2019-01-01T00:00:00Z", "tester", "1", "no reason"},
		{3, "00000000-0000-0000-0000-000000000001", "1", "2019-10-11T00:00:00Z", "tester", "1", "no reason"},
		{4, "00000000-0000-0000-0000-000000000002", "2", "2019-10-11T00:00:00Z", "tester", "1", "no reason so far"},
		{5, "00000000-0000-0000-0000-000000000003", "0", "2019-10-11T00:00:00Z", "tester", "0", "disabled one"},
	}
	compareConfigurations(f, configurations, expected)

	f.PrintReport()
}

func checkDeleteExistingConfiguration() {
	f := frisby.Create("Check what happens when existing configuration is deleted")

	configurations := readConfigurations(f)
	f.ExpectJsonLength("", 6)

	expected := []ClusterConfiguration{
		{0, "00000000-0000-0000-0000-000000000000", "0", "2019-01-01T00:00:00Z", "tester", "0", "no reason"},
		{1, "00000000-0000-0000-0000-000000000000", "1", "2019-01-01T00:00:00Z", "tester", "0", "no reason"},
		{2, "00000000-0000-0000-0000-000000000000", "2", "2019-01-01T00:00:00Z", "tester", "1", "no reason"},
		{3, "00000000-0000-0000-0000-000000000001", "1", "2019-10-11T00:00:00Z", "tester", "1", "no reason"},
		{4, "00000000-0000-0000-0000-000000000002", "2", "2019-10-11T00:00:00Z", "tester", "1", "no reason so far"},
		{5, "00000000-0000-0000-0000-000000000003", "0", "2019-10-11T00:00:00Z", "tester", "0", "disabled one"},
	}
	compareConfigurations(f, configurations, expected)

	f.Delete(API_URL + "client/configuration/0")
	f.Send()
	f.ExpectStatus(200)

	configurations = readConfigurations(f)
	f.ExpectJsonLength("", 5)

	expected = []ClusterConfiguration{
		{1, "00000000-0000-0000-0000-000000000000", "1", "2019-01-01T00:00:00Z", "tester", "0", "no reason"},
		{2, "00000000-0000-0000-0000-000000000000", "2", "2019-01-01T00:00:00Z", "tester", "1", "no reason"},
		{3, "00000000-0000-0000-0000-000000000001", "1", "2019-10-11T00:00:00Z", "tester", "1", "no reason"},
		{4, "00000000-0000-0000-0000-000000000002", "2", "2019-10-11T00:00:00Z", "tester", "1", "no reason so far"},
		{5, "00000000-0000-0000-0000-000000000003", "0", "2019-10-11T00:00:00Z", "tester", "0", "disabled one"},
	}
	compareConfigurations(f, configurations, expected)

	f.PrintReport()
}

func checkDeleteNonExistingConfiguration() {
	f := frisby.Create("Check what happens when non existing configuration is deleted")

	configurations := readConfigurations(f)
	f.ExpectJsonLength("", 5)

	expected := []ClusterConfiguration{
		{1, "00000000-0000-0000-0000-000000000000", "1", "2019-01-01T00:00:00Z", "tester", "0", "no reason"},
		{2, "00000000-0000-0000-0000-000000000000", "2", "2019-01-01T00:00:00Z", "tester", "1", "no reason"},
		{3, "00000000-0000-0000-0000-000000000001", "1", "2019-10-11T00:00:00Z", "tester", "1", "no reason"},
		{4, "00000000-0000-0000-0000-000000000002", "2", "2019-10-11T00:00:00Z", "tester", "1", "no reason so far"},
		{5, "00000000-0000-0000-0000-000000000003", "0", "2019-10-11T00:00:00Z", "tester", "0", "disabled one"},
	}
	compareConfigurations(f, configurations, expected)

	f.Delete(API_URL + "client/configuration/42")
	f.Send()
	f.ExpectStatus(200)

	configurations = readConfigurations(f)
	f.ExpectJsonLength("", 5)

	expected = []ClusterConfiguration{
		{1, "00000000-0000-0000-0000-000000000000", "1", "2019-01-01T00:00:00Z", "tester", "0", "no reason"},
		{2, "00000000-0000-0000-0000-000000000000", "2", "2019-01-01T00:00:00Z", "tester", "1", "no reason"},
		{3, "00000000-0000-0000-0000-000000000001", "1", "2019-10-11T00:00:00Z", "tester", "1", "no reason"},
		{4, "00000000-0000-0000-0000-000000000002", "2", "2019-10-11T00:00:00Z", "tester", "1", "no reason so far"},
		{5, "00000000-0000-0000-0000-000000000003", "0", "2019-10-11T00:00:00Z", "tester", "0", "disabled one"},
	}
	compareConfigurations(f, configurations, expected)

	f.PrintReport()
}

func checkDescribeExistingConfiguration() {
	f := frisby.Create("Check describing (reading) existing configuration")
	f.Get(API_URL + "client/configuration/1")
	f.Send()
	f.ExpectStatus(200)
	f.ExpectContent(`{"no_op":"X", "watch":["a","b","c"]}`)
	f.PrintReport()
}

func checkDescribeNonExistingConfiguration() {
	f := frisby.Create("Check describing (reading) non existing configuration")
	f.Get(API_URL + "client/configuration/42")
	f.Send()
	f.ExpectStatus(400)
	f.PrintReport()
}

// ConfigurationTests run all configuration-related REST API tests.
func ConfigurationTests() {
	checkInitialListOfConfigurations()
	checkEnableExistingConfiguration()
	checkDisableExistingConfiguration()
	checkEnableNonExistingConfiguration()
	checkDisableNonExistingConfiguration()
	checkDeleteExistingConfiguration()
	checkDeleteNonExistingConfiguration()
	checkDescribeExistingConfiguration()
	checkDescribeNonExistingConfiguration()
}
