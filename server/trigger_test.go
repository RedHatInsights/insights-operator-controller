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
func TestNonErrorsTriggerWithoutData(t *testing.T) {
	serv := MockedIOCServer(t, false)
	defer serv.Storage.Close()

	nonErrorTT := []testCase{
		{"GetAllTriggers OK", serv.GetAllTriggers, http.StatusOK, "GET", true, requestData{}, requestData{}, ""},
		{"GetTrigger Not Found", serv.GetTrigger, http.StatusNotFound, "GET", true, requestData{"id": "1"}, requestData{}, ""},
		{"DeleteTrigger Not Found", serv.DeleteTrigger, http.StatusNotFound, "DELETE", true, requestData{"id": "1"}, requestData{}, ""},
		{"ActivateTrigger Not Found", serv.ActivateTrigger, http.StatusNotFound, "PUT", true, requestData{"id": "1"}, requestData{}, ""},
		{"DeactivateTrigger Not Found", serv.DeactivateTrigger, http.StatusNotFound, "PUT", true, requestData{"id": "1"}, requestData{}, ""},
		{"GetClusterTriggers Not Found", serv.GetClusterTriggers, http.StatusNotFound, "GET", true, requestData{"cluster": "1"}, requestData{}, ""},
		{"RegisterClusterTrigger Not Found", serv.RegisterClusterTrigger, http.StatusNotFound, "POST", true, requestData{"cluster": "00000000-0000-0000-0000-000000000000", "trigger": "must-gather"}, requestData{"username": "tester", "reason": "test", "link": "link"}, ""},
	}

	for _, tt := range nonErrorTT {
		testRequest(t, tt)
	}
}

// TestNonErrorsTriggerWithData tests OK behaviour with mock data
func TestNonErrorsTriggerWithData(t *testing.T) {
	serv := MockedIOCServer(t, true)
	defer serv.Storage.Close()

	nonErrorTT := []testCase{
		{"GetAllTriggers OK", serv.GetAllTriggers, http.StatusOK, "GET", true, requestData{}, requestData{}, ""},
		{"GetTrigger OK", serv.GetTrigger, http.StatusOK, "GET", true, requestData{"id": "1"}, requestData{}, ""},
		{"ActivateTrigger OK", serv.ActivateTrigger, http.StatusOK, "PUT", true, requestData{"id": "1"}, requestData{}, ""},
		{"DeactivateTrigger OK", serv.DeactivateTrigger, http.StatusOK, "PUT", true, requestData{"id": "1"}, requestData{}, ""},
		{"GetClusterTriggers OK", serv.GetClusterTriggers, http.StatusOK, "GET", true, requestData{"cluster": "1"}, requestData{}, ""},
		{"RegisterClusterTrigger OK", serv.RegisterClusterTrigger, http.StatusOK, "POST", true, requestData{"cluster": "00000000-0000-0000-0000-000000000000", "trigger": "must-gather"}, requestData{"username": "tester", "reason": "test", "link": "link"}, ""},
		{"DeleteTrigger OK", serv.DeleteTrigger, http.StatusOK, "DELETE", true, requestData{"id": "1"}, requestData{}, ""},
	}

	for _, tt := range nonErrorTT {
		testRequest(t, tt)
	}
}

// TestDatabaseErrorTrigger tests unexpected behaviour by closing DB connection (consistency check)
func TestDatabaseErrorTrigger(t *testing.T) {
	serv := MockedIOCServer(t, true)

	dbErrorTT := []testCase{
		{"GetAllTriggers DB error", serv.GetAllTriggers, http.StatusInternalServerError, "GET", true, requestData{}, requestData{}, ""},
		{"GetTrigger DB error", serv.GetTrigger, http.StatusInternalServerError, "GET", true, requestData{"id": "1"}, requestData{}, ""},
		{"DeleteTrigger DB error", serv.DeleteTrigger, http.StatusInternalServerError, "DELETE", true, requestData{"id": "1"}, requestData{}, ""},
		{"ActivateTrigger DB error", serv.ActivateTrigger, http.StatusInternalServerError, "PUT", true, requestData{"id": "1"}, requestData{}, ""},
		{"DeactivateTrigger DB error", serv.DeactivateTrigger, http.StatusInternalServerError, "PUT", true, requestData{"id": "1"}, requestData{}, ""},
		{"GetClusterTriggers DB error", serv.GetClusterTriggers, http.StatusInternalServerError, "GET", true, requestData{"cluster": "1"}, requestData{}, ""},
		{"RegisterClusterTrigger DB error", serv.RegisterClusterTrigger, http.StatusInternalServerError, "POST", true, requestData{"cluster": "00000000-0000-0000-0000-000000000000", "trigger": "must-gather"}, requestData{"username": "tester", "reason": "test", "link": "link"}, ""},
	}

	serv.Storage.Close()

	for _, tt := range dbErrorTT {
		testRequest(t, tt)
	}
}

// TestParameterErrorsTrigger tests wrong request parameters
func TestParameterErrorsTrigger(t *testing.T) {
	serv := MockedIOCServer(t, true)
	defer serv.Storage.Close()

	paramErrorTT := []testCase{
		{"GetTrigger no id", serv.GetTrigger, http.StatusBadRequest, "GET", true, requestData{}, requestData{}, ""},
		{"GetTrigger non-int id", serv.GetTrigger, http.StatusBadRequest, "GET", true, requestData{"id": "non-int"}, requestData{}, ""},
		{"DeleteTrigger no id", serv.DeleteTrigger, http.StatusBadRequest, "DELETE", true, requestData{}, requestData{}, ""},
		{"DeleteTrigger non-int id", serv.DeleteTrigger, http.StatusBadRequest, "DELETE", true, requestData{"id": "test"}, requestData{}, ""},
		{"ActivateTrigger no id", serv.ActivateTrigger, http.StatusBadRequest, "PUT", true, requestData{}, requestData{}, ""},
		{"ActivateTrigger non-int id", serv.ActivateTrigger, http.StatusBadRequest, "PUT", true, requestData{"id": "non-int"}, requestData{}, ""},
		{"DeactivateTrigger no id", serv.DeactivateTrigger, http.StatusBadRequest, "PUT", true, requestData{}, requestData{}, ""},
		{"DeactivateTrigger non-int id", serv.DeactivateTrigger, http.StatusBadRequest, "PUT", true, requestData{"id": "non-int"}, requestData{}, ""},
		{"GetClusterTriggers no id", serv.GetClusterTriggers, http.StatusBadRequest, "GET", true, requestData{}, requestData{}, ""},
		{"RegisterClusterTrigger no trigger", serv.RegisterClusterTrigger, http.StatusBadRequest, "POST", true, requestData{"cluster": "1"}, requestData{"username": "tester", "reason": "test", "link": "link"}, ""},
		{"RegisterClusterTrigger no cluster", serv.RegisterClusterTrigger, http.StatusBadRequest, "POST", true, requestData{"trigger": "must-gather"}, requestData{"username": "tester", "reason": "test", "link": "link"}, ""},
		{"RegisterClusterTrigger no link", serv.RegisterClusterTrigger, http.StatusBadRequest, "POST", true, requestData{"cluster": "1", "trigger": "must-gather"}, requestData{"username": "tester", "reason": "test"}, ""},
		{"RegisterClusterTrigger no reason", serv.RegisterClusterTrigger, http.StatusBadRequest, "POST", true, requestData{"cluster": "1", "trigger": "must-gather"}, requestData{"username": "tester", "link": "link"}, ""},
		{"RegisterClusterTrigger no username", serv.RegisterClusterTrigger, http.StatusBadRequest, "POST", true, requestData{"cluster": "1", "trigger": "must-gather"}, requestData{"reason": "test", "link": "link"}, ""},
	}

	for _, tt := range paramErrorTT {
		testRequest(t, tt)
	}
}
