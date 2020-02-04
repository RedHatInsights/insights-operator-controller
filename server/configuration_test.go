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

package server_test

import (
	"net/http"
	"testing"
)

// TestNonErrorsConfigurationWithoutData tests OK behaviour with empty DB (schema only)
func TestNonErrorsConfigurationWithoutData(t *testing.T) {
	serv := MockedIOCServer(t, false)

	nonErrorTT := []testCase{
		{"GetConfiguration Not Found", serv.GetConfiguration, http.StatusNotFound, "GET", true, requestData{"id": "1"}, requestData{}},
		{"DeleteConfiguration Not Found", serv.DeleteConfiguration, http.StatusNotFound, "DELETE", false, requestData{"id": "1"}, requestData{}},
		{"GetAllConfigurations Empty OK", serv.GetAllConfigurations, http.StatusOK, "GET", true, requestData{}, requestData{}},
		{"GetClusterConfiguration OK", serv.GetClusterConfiguration, http.StatusNotFound, "GET", true, requestData{"cluster": "1"}, requestData{}},
		{"EnableConfiguration OK", serv.EnableConfiguration, http.StatusNotFound, "PUT", false, requestData{"id": "1"}, requestData{}},
		{"DisableConfiguration OK", serv.DisableConfiguration, http.StatusNotFound, "PUT", false, requestData{"id": "1"}, requestData{}},
	}

	for _, tt := range nonErrorTT {
		testRequest(t, tt)
	}
}

// TestNonErrorsConfigurationWithData tests OK behaviour with mock data
func TestNonErrorsConfigurationWithData(t *testing.T) {
	serv := MockedIOCServer(t, true)

	nonErrorTT := []testCase{
		{"GetConfiguration OK", serv.GetConfiguration, http.StatusOK, "GET", true, requestData{"id": "1"}, requestData{}},
		{"GetAllConfigurations OK", serv.GetAllConfigurations, http.StatusOK, "GET", true, requestData{}, requestData{}},
		{"GetClusterConfiguration OK", serv.GetClusterConfiguration, http.StatusOK, "GET", true, requestData{"cluster": "1"}, requestData{}},
		{"EnableConfiguration OK", serv.EnableConfiguration, http.StatusAccepted, "PUT", false, requestData{"id": "1"}, requestData{}},
		{"DisableConfiguration OK", serv.DisableConfiguration, http.StatusAccepted, "PUT", false, requestData{"id": "1"}, requestData{}},
		{"DeleteConfiguration OK", serv.DeleteConfiguration, http.StatusAccepted, "DELETE", false, requestData{"id": "1"}, requestData{}},
	}

	for _, tt := range nonErrorTT {
		testRequest(t, tt)
	}
}

/*
	{"NewClusterConfiguration OK", serv.NewClusterConfiguration, http.StatusAccepted, "POST", false, requestData{"cluster": "1"}, requestData{"username": "test", "reason": "unknown", "description": "testing"}},
	{"EnableClusterConfiguration OK", serv.EnableClusterConfiguration, http.StatusAccepted, "PUT", false, requestData{}, requestData{}},
	{"DisableClusterConfiguration OK", serv.DisableClusterConfiguration, http.StatusAccepted, "PUT", false, requestData{}, requestData{}},
*/
