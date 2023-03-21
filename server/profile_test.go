/*
Copyright © 2019, 2020, 2021, 2022, 2023 Red Hat, Inc.

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

// Documentation in literate-programming-style is available at:
// https://redhatinsights.github.io/insights-operator-controller/packages/server/profile_test.html

import (
	"net/http"
	"testing"
)

// TestNonErrorsConfigurationWithoutData tests OK behaviour with empty DB (schema only)
func TestNonErrorsProfileWithoutData(t *testing.T) {
	serv := MockedIOCServer(t, false)
	defer serv.Storage.Close()

	nonErrorTT := []testCase{
		{"GetConfigurationProfile Not Found", serv.GetConfigurationProfile, http.StatusNotFound, "GET", true, requestData{"id": "1"}, requestData{}, ""},
		{"ListConfigurationProfiles OK", serv.ListConfigurationProfiles, http.StatusOK, "GET", true, requestData{}, requestData{}, ""},
		{"DeleteConfigurationProfile Not Found", serv.DeleteConfigurationProfile, http.StatusNotFound, "DELETE", true, requestData{"id": "1"}, requestData{}, ""},
		{"ChangeConfigurationProfile Not Found", serv.ChangeConfigurationProfile, http.StatusNotFound, "PUT", true, requestData{"id": "1"}, requestData{"username": "tester", "description": "test"}, "Test config"},
		{"NewConfigurationProfile OK", serv.NewConfigurationProfile, http.StatusCreated, "POST", true, requestData{}, requestData{"username": "tester", "description": "test"}, "Test config"},
	}

	for _, tt := range nonErrorTT {
		testRequest(t, &tt)
	}
}

// TestNonErrorsProfileWithData tests OK behaviour with mock data
func TestNonErrorsProfileWithData(t *testing.T) {
	serv := MockedIOCServer(t, true)
	defer serv.Storage.Close()

	nonErrorTT := []testCase{
		{"GetConfigurationProfile OK", serv.GetConfigurationProfile, http.StatusOK, "GET", true, requestData{"id": "1"}, requestData{}, ""},
		{"ListConfigurationProfiles OK", serv.ListConfigurationProfiles, http.StatusOK, "GET", true, requestData{}, requestData{}, ""},
		{"DeleteConfigurationProfile OK", serv.DeleteConfigurationProfile, http.StatusOK, "DELETE", true, requestData{"id": "1"}, requestData{}, ""},
		{"NewConfigurationProfile OK", serv.NewConfigurationProfile, http.StatusCreated, "POST", true, requestData{}, requestData{"username": "tester", "description": "test"}, "Test config"},
		{"ChangeConfigurationProfile OK", serv.ChangeConfigurationProfile, http.StatusOK, "PUT", true, requestData{"id": "4"}, requestData{"username": "tester", "description": "test"}, "Test config"},
	}

	for _, tt := range nonErrorTT {
		testRequest(t, &tt)
	}
}

// TestDatabaseErrorProfile tests unexpected behaviour by closing DB connection (consistency check)
func TestDatabaseErrorProfile(t *testing.T) {
	serv := MockedIOCServer(t, true)

	dbErrorTT := []testCase{
		{"GetConfigurationProfile Not Found", serv.GetConfigurationProfile, http.StatusInternalServerError, "GET", true, requestData{"id": "1"}, requestData{}, ""},
		{"ListConfigurationProfiles OK", serv.ListConfigurationProfiles, http.StatusInternalServerError, "GET", true, requestData{}, requestData{}, ""},
		{"DeleteConfigurationProfile Not Found", serv.DeleteConfigurationProfile, http.StatusInternalServerError, "DELETE", true, requestData{"id": "1"}, requestData{}, ""},
		{"ChangeConfigurationProfile Not Found", serv.ChangeConfigurationProfile, http.StatusInternalServerError, "PUT", true, requestData{"id": "1"}, requestData{"username": "tester", "description": "test"}, "Test config"},
		{"NewConfigurationProfile OK", serv.NewConfigurationProfile, http.StatusInternalServerError, "POST", true, requestData{}, requestData{"username": "tester", "description": "test"}, "Test config"},
	}

	serv.Storage.Close()

	for _, tt := range dbErrorTT {
		testRequest(t, &tt)
	}
}

// TestParameterErrorsProfile tests wrong request paramaters
func TestParameterErrorsProfile(t *testing.T) {
	serv := MockedIOCServer(t, true)
	defer serv.Storage.Close()

	paramErrorTT := []testCase{
		{"GetConfigurationProfile no id", serv.GetConfigurationProfile, http.StatusBadRequest, "GET", true, requestData{}, requestData{}, ""},
		{"GetConfigurationProfile non-int id", serv.GetConfigurationProfile, http.StatusBadRequest, "GET", true, requestData{"id": "non-int"}, requestData{}, ""},
		{"DeleteConfigurationProfile no id", serv.DeleteConfigurationProfile, http.StatusBadRequest, "DELETE", true, requestData{}, requestData{}, ""},
		{"DeleteConfigurationProfile non-int id", serv.DeleteConfigurationProfile, http.StatusBadRequest, "DELETE", true, requestData{"id": "non-int"}, requestData{}, ""},
		{"ChangeConfigurationProfile no id", serv.ChangeConfigurationProfile, http.StatusBadRequest, "PUT", true, requestData{}, requestData{"username": "tester", "description": "test"}, "Test config"},
		{"ChangeConfigurationProfile non-int id", serv.ChangeConfigurationProfile, http.StatusBadRequest, "PUT", true, requestData{"id": "non-int"}, requestData{"username": "tester", "description": "test"}, "Test config"},
		{"ChangeConfigurationProfile no description", serv.ChangeConfigurationProfile, http.StatusBadRequest, "PUT", true, requestData{"id": "1"}, requestData{"username": "tester"}, "Test config"},
		{"ChangeConfigurationProfile no username", serv.ChangeConfigurationProfile, http.StatusBadRequest, "PUT", true, requestData{"id": "1"}, requestData{"description": "test"}, "Test config"},
		{"ChangeConfigurationProfile no config in body", serv.ChangeConfigurationProfile, http.StatusBadRequest, "PUT", true, requestData{"id": "1"}, requestData{"username": "tester", "description": "test"}, ""},
		{"NewConfigurationProfile no description", serv.NewConfigurationProfile, http.StatusBadRequest, "POST", true, requestData{}, requestData{"username": "tester"}, "Test config"},
		{"NewConfigurationProfile no username", serv.NewConfigurationProfile, http.StatusBadRequest, "POST", true, requestData{}, requestData{"description": "test"}, "Test config"},
		{"NewConfigurationProfile no config in body", serv.NewConfigurationProfile, http.StatusBadRequest, "POST", true, requestData{}, requestData{"username": "tester", "description": "test"}, ""},
	}

	for _, tt := range paramErrorTT {
		testRequest(t, &tt)
	}
}
