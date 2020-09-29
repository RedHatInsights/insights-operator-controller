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

// Documentation in literate-programming-style is available at:
// https://redhatinsights.github.io/insights-operator-controller/packages/server/operator_test.html

import (
	"net/http"
	"testing"
)

// TestNonErrorsConfigurationWithoutData tests OK behaviour with empty DB (schema only)
func TestNonErrorsOperatorWithoutData(t *testing.T) {
	serv := MockedIOCServer(t, false)
	defer serv.Storage.Close()

	nonErrorTT := []testCase{
		{"ReadConfigurationForOperator Not Found", serv.ReadConfigurationForOperator, http.StatusNotFound, "GET", true, requestData{"cluster": "1"}, requestData{}, ""},
		{"GetActiveTriggersForCluster Not Found", serv.GetActiveTriggersForCluster, http.StatusNotFound, "GET", true, requestData{"cluster": "1"}, requestData{}, ""},
		{"AckTriggerForCluster Not Found", serv.AckTriggerForCluster, http.StatusNotFound, "GET", true, requestData{"cluster": "1", "trigger": "1"}, requestData{}, ""},
		{"RegisterCluster OK", serv.RegisterCluster, http.StatusCreated, "PUT", true, requestData{"cluster": "1"}, requestData{}, ""},
	}

	for _, tt := range nonErrorTT {
		testRequest(t, tt)
	}
}

// TestNonErrorsOperatorWithData tests OK behaviour with mock data
func TestNonErrorsOperatorWithData(t *testing.T) {
	serv := MockedIOCServer(t, true)
	defer serv.Storage.Close()

	nonErrorTT := []testCase{
		{"ReadConfigurationForOperator OK", serv.ReadConfigurationForOperator, http.StatusOK, "GET", true, requestData{"cluster": "00000000-0000-0000-0000-000000000000"}, requestData{}, ""},
		{"GetActiveTriggersForCluster OK", serv.GetActiveTriggersForCluster, http.StatusOK, "GET", true, requestData{"cluster": "00000000-0000-0000-0000-000000000000"}, requestData{}, ""},
		{"GetActiveTriggersForCluster No triggers OK", serv.GetActiveTriggersForCluster, http.StatusOK, "GET", true, requestData{"cluster": "00000000-0000-0000-0000-000000000004"}, requestData{}, ""},
		{"AckTriggerForCluster OK", serv.AckTriggerForCluster, http.StatusOK, "GET", true, requestData{"cluster": "00000000-0000-0000-0000-000000000000", "trigger": "1"}, requestData{}, ""},
		{"RegisterCluster OK", serv.RegisterCluster, http.StatusCreated, "PUT", true, requestData{"cluster": "1"}, requestData{}, ""},
	}

	for _, tt := range nonErrorTT {
		testRequest(t, tt)
	}
}

// TestDatabaseErrorOperator tests unexpected behaviour by closing DB connection (consistency check)
func TestDatabaseErrorOperator(t *testing.T) {
	serv := MockedIOCServer(t, true)

	dbErrorTT := []testCase{
		{"ReadConfigurationForOperator DB error", serv.ReadConfigurationForOperator, http.StatusInternalServerError, "GET", true, requestData{"cluster": "00000000-0000-0000-0000-000000000000"}, requestData{}, ""},
		{"GetActiveTriggersForCluster DB error", serv.GetActiveTriggersForCluster, http.StatusInternalServerError, "GET", true, requestData{"cluster": "00000000-0000-0000-0000-000000000000"}, requestData{}, ""},
		{"AckTriggerForCluster DB error", serv.AckTriggerForCluster, http.StatusInternalServerError, "GET", true, requestData{"cluster": "00000000-0000-0000-0000-000000000000", "trigger": "1"}, requestData{}, ""},
		{"RegisterCluster DB error", serv.RegisterCluster, http.StatusInternalServerError, "PUT", true, requestData{"cluster": "1"}, requestData{}, ""},
	}

	serv.Storage.Close()

	for _, tt := range dbErrorTT {
		testRequest(t, tt)
	}
}

// TestParameterErrorsOperator tests wrong request paramaters
func TestParameterErrorsOperator(t *testing.T) {
	serv := MockedIOCServer(t, true)
	defer serv.Storage.Close()

	paramErrorTT := []testCase{
		{"ReadConfigurationForOperator no cluster", serv.ReadConfigurationForOperator, http.StatusBadRequest, "GET", true, requestData{}, requestData{}, ""},
		{"GetActiveTriggersForCluster no cluster", serv.GetActiveTriggersForCluster, http.StatusBadRequest, "GET", true, requestData{}, requestData{}, ""},
		{"AckTriggerForCluster no trigger", serv.AckTriggerForCluster, http.StatusBadRequest, "GET", true, requestData{"cluster": "00000000-0000-0000-0000-000000000000"}, requestData{}, ""},
		{"AckTriggerForCluster no cluster", serv.AckTriggerForCluster, http.StatusBadRequest, "GET", true, requestData{"trigger": "1"}, requestData{}, ""},
		{"RegisterCluster no cluster", serv.RegisterCluster, http.StatusBadRequest, "PUT", true, requestData{}, requestData{}, ""},
	}

	for _, tt := range paramErrorTT {
		testRequest(t, tt)
	}
}
