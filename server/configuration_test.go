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
	defer serv.Storage.Close()

	nonErrorTT := []testCase{
		{"GetConfiguration Not Found", serv.GetConfiguration, http.StatusNotFound, "GET", true, requestData{"id": "1"}, requestData{}, ""},
		{"DeleteConfiguration Not Found", serv.DeleteConfiguration, http.StatusNotFound, "DELETE", false, requestData{"id": "1"}, requestData{}, ""},
		{"GetAllConfigurations OK", serv.GetAllConfigurations, http.StatusOK, "GET", true, requestData{}, requestData{}, ""},
		{"GetClusterConfiguration Not Found", serv.GetClusterConfiguration, http.StatusNotFound, "GET", true, requestData{"cluster": "1"}, requestData{}, ""},
		{"EnableConfiguration Not Found", serv.EnableConfiguration, http.StatusNotFound, "PUT", false, requestData{"id": "1"}, requestData{}, ""},
		{"DisableConfiguration Not Found", serv.DisableConfiguration, http.StatusNotFound, "PUT", false, requestData{"id": "1"}, requestData{}, ""},
	}

	for _, tt := range nonErrorTT {
		testRequest(t, tt)
	}
}

// TestNonErrorsConfigurationWithData tests OK behaviour with mock data
func TestNonErrorsConfigurationWithData(t *testing.T) {
	serv := MockedIOCServer(t, true)
	defer serv.Storage.Close()

	nonErrorTT := []testCase{
		{"GetConfiguration OK", serv.GetConfiguration, http.StatusOK, "GET", true, requestData{"id": "1"}, requestData{}, ""},
		{"GetAllConfigurations OK", serv.GetAllConfigurations, http.StatusOK, "GET", true, requestData{}, requestData{}, ""},
		{"GetClusterConfiguration OK", serv.GetClusterConfiguration, http.StatusOK, "GET", true, requestData{"cluster": "00000000-0000-0000-0000-000000000001"}, requestData{}, ""},
		{"EnableConfiguration OK", serv.EnableConfiguration, http.StatusOK, "PUT", true, requestData{"id": "1"}, requestData{}, ""},
		{"DisableConfiguration OK", serv.DisableConfiguration, http.StatusOK, "PUT", true, requestData{"id": "1"}, requestData{}, ""},
		{"DeleteConfiguration OK", serv.DeleteConfiguration, http.StatusOK, "DELETE", true, requestData{"id": "1"}, requestData{}, ""},
		{"EnableClusterConfiguration OK", serv.EnableClusterConfiguration, http.StatusOK, "PUT", true, requestData{"cluster": "00000000-0000-0000-0000-000000000000"}, requestData{"username": "tester", "reason": "test"}, ""},
		{"DisableClusterConfiguration OK", serv.DisableClusterConfiguration, http.StatusOK, "PUT", true, requestData{"cluster": "00000000-0000-0000-0000-000000000000"}, requestData{"username": "tester", "reason": "test"}, ""},
		{"NewClusterConfiguration OK", serv.NewClusterConfiguration, http.StatusOK, "POST", false, requestData{"cluster": "00000000-0000-0000-0000-000000000000"}, requestData{"username": "test", "reason": "unknown", "description": "testing"}, "Test config"},
	}

	for _, tt := range nonErrorTT {
		testRequest(t, tt)
	}
}

// TestDatabaseErrorConfiguration tests unexpected behaviour by closing DB connection (consistency check)
func TestDatabaseErrorConfiguration(t *testing.T) {
	serv := MockedIOCServer(t, true)

	dbErrorTT := []testCase{
		{"GetConfiguration DB error", serv.GetConfiguration, http.StatusInternalServerError, "GET", false, requestData{"id": "1"}, requestData{}, ""},
		{"GetAllConfigurations DB error", serv.GetAllConfigurations, http.StatusInternalServerError, "GET", false, requestData{}, requestData{}, ""},
		{"GetClusterConfiguration DB error", serv.GetClusterConfiguration, http.StatusInternalServerError, "GET", false, requestData{"cluster": "1"}, requestData{}, ""},
		{"EnableConfiguration DB error", serv.EnableConfiguration, http.StatusInternalServerError, "PUT", false, requestData{"id": "1"}, requestData{}, ""},
		{"DisableConfiguration DB error", serv.DisableConfiguration, http.StatusInternalServerError, "PUT", false, requestData{"id": "1"}, requestData{}, ""},
		{"DeleteConfiguration DB error", serv.DeleteConfiguration, http.StatusInternalServerError, "DELETE", false, requestData{"id": "1"}, requestData{}, ""},
		{"EnableClusterConfiguration DB error", serv.EnableClusterConfiguration, http.StatusInternalServerError, "PUT", false, requestData{"cluster": "00000000-0000-0000-0000-000000000000"}, requestData{"username": "tester", "reason": "test"}, ""},
		{"DisableClusterConfiguration DB error", serv.DisableClusterConfiguration, http.StatusInternalServerError, "PUT", false, requestData{"cluster": "00000000-0000-0000-0000-000000000000"}, requestData{"username": "tester", "reason": "test"}, ""},
		{"NewClusterConfiguration DB error", serv.NewClusterConfiguration, http.StatusInternalServerError, "POST", false, requestData{"cluster": "00000000-0000-0000-0000-000000000000"}, requestData{"username": "test", "reason": "unknown", "description": "testing"}, "Test config"},
	}

	// close storage
	serv.Storage.Close()

	for _, tt := range dbErrorTT {
		testRequest(t, tt)
	}
}

// TestParameterErrorsConfiguration tests wrong request parameters
func TestParameterErrorsConfiguration(t *testing.T) {
	serv := MockedIOCServer(t, true)
	defer serv.Storage.Close()

	paramErrorTT := []testCase{
		{"GetConfiguration no id", serv.GetConfiguration, http.StatusBadRequest, "GET", true, requestData{}, requestData{}, ""},
		{"GetConfiguration non-int id", serv.GetConfiguration, http.StatusBadRequest, "GET", true, requestData{"id": "non-int"}, requestData{}, ""},
		{"DeleteConfiguration no id", serv.DeleteConfiguration, http.StatusBadRequest, "DELETE", true, requestData{}, requestData{}, ""},
		{"DeleteConfiguration non-int id", serv.DeleteConfiguration, http.StatusBadRequest, "DELETE", true, requestData{"id": "non-int"}, requestData{}, ""},
		{"GetClusterConfiguration no id", serv.GetClusterConfiguration, http.StatusBadRequest, "GET", true, requestData{}, requestData{}, ""},
		{"EnableConfiguration no id", serv.EnableConfiguration, http.StatusBadRequest, "PUT", true, requestData{}, requestData{}, ""},
		{"DisableConfiguration no id", serv.DisableConfiguration, http.StatusBadRequest, "PUT", true, requestData{}, requestData{}, ""},
		{"EnableConfiguration non-int id", serv.EnableConfiguration, http.StatusBadRequest, "PUT", true, requestData{"id": "non-int"}, requestData{}, ""},
		{"DisableConfiguration non-int id", serv.DisableConfiguration, http.StatusBadRequest, "PUT", true, requestData{"id": "non-int"}, requestData{}, ""},
		{"EnableClusterConfiguration no cluster", serv.EnableClusterConfiguration, http.StatusBadRequest, "PUT", false, requestData{}, requestData{"username": "tester", "reason": "test"}, ""},
		{"EnableClusterConfiguration no reason", serv.EnableClusterConfiguration, http.StatusBadRequest, "PUT", false, requestData{"cluster": "00000000-0000-0000-0000-000000000000"}, requestData{"username": "tester"}, ""},
		{"EnableClusterConfiguration no username", serv.EnableClusterConfiguration, http.StatusBadRequest, "PUT", false, requestData{"cluster": "00000000-0000-0000-0000-000000000000"}, requestData{"reason": "test"}, ""},
		{"DisableClusterConfiguration no cluster", serv.DisableClusterConfiguration, http.StatusBadRequest, "PUT", false, requestData{}, requestData{"username": "tester", "reason": "test"}, ""},
		{"DisableClusterConfiguration no reason", serv.DisableClusterConfiguration, http.StatusBadRequest, "PUT", false, requestData{"cluster": "00000000-0000-0000-0000-000000000000"}, requestData{"username": "tester"}, ""},
		{"DisableClusterConfiguration no username", serv.DisableClusterConfiguration, http.StatusBadRequest, "PUT", false, requestData{"cluster": "00000000-0000-0000-0000-000000000000"}, requestData{"reason": "test"}, ""},
		{"NewClusterConfiguration no cluster", serv.NewClusterConfiguration, http.StatusBadRequest, "POST", false, requestData{}, requestData{"username": "test", "reason": "unknown", "description": "testing"}, "Test config"},
		{"NewClusterConfiguration no username", serv.NewClusterConfiguration, http.StatusBadRequest, "POST", false, requestData{"cluster": "1"}, requestData{"reason": "unknown", "description": "testing"}, "Test config"},
		{"NewClusterConfiguration no reason", serv.NewClusterConfiguration, http.StatusBadRequest, "POST", false, requestData{"cluster": "1"}, requestData{"username": "test", "description": "testing"}, "Test config"},
		{"NewClusterConfiguration no description", serv.NewClusterConfiguration, http.StatusBadRequest, "POST", false, requestData{"cluster": "1"}, requestData{"username": "test", "reason": "unknown"}, "Test config"},
		{"NewClusterConfiguration no config in body", serv.NewClusterConfiguration, http.StatusBadRequest, "POST", false, requestData{"cluster": "00000000-0000-0000-0000-000000000000"}, requestData{"username": "test", "reason": "unknown", "description": "testing"}, ""},
	}

	for _, tt := range paramErrorTT {
		testRequest(t, tt)
	}
}
