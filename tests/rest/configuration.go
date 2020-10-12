/*
Copyright © 2019, 2020 Red Hat, Inc.

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

// Documentation in literate-programming-style is available at:
// https://redhatinsights.github.io/insights-operator-controller/packages/tests/rest/configuration.html

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

// ClusterConfigurationsResponse represents default response for cluster configuration request
type ClusterConfigurationsResponse struct {
	Status        string                 `json:"status"`
	Configuration []ClusterConfiguration `json:"configuration"`
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
		}
	}
}

// readConfigurations tries to read list of configurations from the HTTP server response
func readConfigurations(f *frisby.Frisby) []ClusterConfiguration {
	f.Get(API_URL + "client/configuration")
	f.Send()
	f.ExpectStatus(200)
	f.ExpectHeader("Content-Type", "application/json; charset=utf-8")

	// default return value from this function
	response := ClusterConfigurationsResponse{}

	// try to read payload from response
	text, err := f.Resp.Content()
	if err != nil {
		f.AddError(err.Error())
	} else {
		// try to unmarshall response body and check if it's correct
		err = json.Unmarshal(text, &response)
		if err != nil {
			// any error needs to be recorded
			f.AddError(err.Error())
		}
	}
	return response.Configuration
}

func checkNumberOfConfigurations(f *frisby.Frisby, configurations []ClusterConfiguration, expected int) {
	if len(configurations) != expected {
		f.AddError(fmt.Sprintf("Number of returned configurations %d differs from expected number %d", len(configurations), expected))
	}
}

func checkInitialListOfConfigurations() {
	f := frisby.Create("Check list of configurations")

	configurations := readConfigurations(f)
	checkNumberOfConfigurations(f, configurations, 6)

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
	checkNumberOfConfigurations(f, configurations, 6)

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
	checkNumberOfConfigurations(f, configurations, 6)

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
	checkNumberOfConfigurations(f, configurations, 6)

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
	checkNumberOfConfigurations(f, configurations, 6)

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
	checkNumberOfConfigurations(f, configurations, 6)

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
	f.ExpectStatus(404)

	configurations = readConfigurations(f)
	checkNumberOfConfigurations(f, configurations, 6)

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
	checkNumberOfConfigurations(f, configurations, 6)

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
	f.ExpectStatus(404)

	configurations = readConfigurations(f)
	checkNumberOfConfigurations(f, configurations, 6)

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
	checkNumberOfConfigurations(f, configurations, 6)

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
	checkNumberOfConfigurations(f, configurations, 5)

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
	checkNumberOfConfigurations(f, configurations, 5)

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
	f.ExpectStatus(404)

	configurations = readConfigurations(f)
	checkNumberOfConfigurations(f, configurations, 5)

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
	f.ExpectContent(`{"configuration":"{\"no_op\":\"X\", \"watch\":[\"a\",\"b\",\"c\"]}","status":"ok"}`)
	f.PrintReport()
}

func checkDescribeNonExistingConfiguration() {
	f := frisby.Create("Check describing (reading) non existing configuration")
	f.Get(API_URL + "client/configuration/42")
	f.Send()
	f.ExpectStatus(404)
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
