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

// TestNonErrorsClusterWithoutData tests OK behaviour with empty DB (schema only)
func TestNonErrorsClusterWithoutData(t *testing.T) {
	serv := MockedIOCServer(t, false)
	defer serv.Storage.Close()

	nonErrorTT := []testCase{
		{"GetClusters OK", serv.GetClusters, http.StatusOK, "GET", true, requestData{}, requestData{}, ""},
		{"GetClusterByID OK", serv.GetClusterByID, http.StatusNotFound, "GET", true, requestData{"id": "1"}, requestData{}, ""},
		{"SearchCluster by id OK", serv.SearchCluster, http.StatusNotFound, "GET", true, requestData{}, requestData{"id": "1"}, ""},
		{"SearchCluster by name OK", serv.SearchCluster, http.StatusNotFound, "GET", true, requestData{}, requestData{"name": "test"}, ""},
		{"DeleteCluster OK", serv.DeleteCluster, http.StatusNotFound, "DELETE", false, requestData{"id": "2"}, requestData{}, ""},
		{"NewCluster OK", serv.NewCluster, http.StatusCreated, "POST", false, requestData{"name": "test"}, requestData{}, ""},
	}

	for _, tt := range nonErrorTT {
		testRequest(t, tt)
	}
}

// TestNonErrorsClusterWithData tests OK behaviour with mock data
func TestNonErrorsClusterWithData(t *testing.T) {
	serv := MockedIOCServer(t, true)
	defer serv.Storage.Close()

	nonErrorTT := []testCase{
		{"GetClusters OK", serv.GetClusters, http.StatusOK, "GET", true, requestData{}, requestData{}, ""},
		{"NewCluster OK", serv.NewCluster, http.StatusCreated, "POST", false, requestData{"name": "test"}, requestData{}, ""},
		{"GetClusterByID OK", serv.GetClusterByID, http.StatusOK, "GET", true, requestData{"id": "1"}, requestData{}, ""},
		{"SearchCluster OK", serv.SearchCluster, http.StatusOK, "GET", true, requestData{}, requestData{"name": "test"}, ""},
		{"DeleteCluster OK", serv.DeleteCluster, http.StatusOK, "DELETE", false, requestData{"id": "1"}, requestData{}, ""},
	}

	for _, tt := range nonErrorTT {
		testRequest(t, tt)
	}
}

// TestDatabaseErrorsCluster tests unexpected behaviour by closing DB connection (consistency check)
func TestDatabaseErrorCluster(t *testing.T) {
	serv := MockedIOCServer(t, false)

	dbErrorTT := []testCase{
		{"GetClusters DB error", serv.GetClusters, http.StatusInternalServerError, "GET", true, requestData{}, requestData{}, ""},
		{"NewCluster DB error", serv.NewCluster, http.StatusInternalServerError, "POST", false, requestData{"name": "test"}, requestData{}, ""},
		{"GetClusterByID DB error", serv.GetClusterByID, http.StatusInternalServerError, "GET", true, requestData{"id": "1"}, requestData{}, ""},
		{"DeleteCluster DB error", serv.DeleteCluster, http.StatusInternalServerError, "DELETE", false, requestData{"id": "1"}, requestData{}, ""},
		{"SearchCluster DB error", serv.SearchCluster, http.StatusInternalServerError, "GET", true, requestData{}, requestData{"name": "test"}, ""},
	}

	// close DB
	serv.Storage.Close()

	for _, tt := range dbErrorTT {
		testRequest(t, tt)
	}
}

// TestParameterErrorsCluster tests wrong request parameters
func TestParameterErrorsCluster(t *testing.T) {
	serv := MockedIOCServer(t, true)
	defer serv.Storage.Close()

	paramErrorTT := []testCase{
		{"NewCluster no param", serv.NewCluster, http.StatusBadRequest, "POST", false, requestData{}, requestData{}, ""},
		{"NewCluster empty name key", serv.NewCluster, http.StatusBadRequest, "POST", false, requestData{"name": ""}, requestData{}, ""},
		{"GetClusterByID no param", serv.GetClusterByID, http.StatusBadRequest, "GET", true, requestData{}, requestData{}, ""},
		{"GetClusterByID non-str id", serv.GetClusterByID, http.StatusBadRequest, "GET", true, requestData{"id": "test"}, requestData{}, ""},
		{"DeleteCluster no param", serv.DeleteCluster, http.StatusBadRequest, "DELETE", false, requestData{}, requestData{}, ""},
		{"DeleteCluster by name", serv.DeleteCluster, http.StatusBadRequest, "DELETE", false, requestData{"name": "test"}, requestData{}, ""},
		{"DeleteCluster by non-int id", serv.DeleteCluster, http.StatusBadRequest, "DELETE", false, requestData{"id": "test"}, requestData{}, ""},
		{"SearchCluster no params", serv.SearchCluster, http.StatusBadRequest, "GET", true, requestData{}, requestData{}, ""},
		{"SearchCluster wrong data type", serv.SearchCluster, http.StatusBadRequest, "GET", true, requestData{"name": ""}, requestData{}, ""},
	}

	for _, tt := range paramErrorTT {
		testRequest(t, tt)
	}
}
