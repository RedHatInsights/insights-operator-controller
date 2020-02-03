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

func TestNonErrorsCluster(t *testing.T) {
	serv := MockedIOCServer(t)

	nonErrorTT := []testCase{
		{"GetClusters OK", serv.GetClusters, http.StatusOK, "GET", true, requestData{}},
		{"NewCluster OK", serv.NewCluster, http.StatusCreated, "POST", false, requestData{"name": "test"}},
		{"GetClusterByID OK", serv.GetClusterByID, http.StatusOK, "GET", true, requestData{"id": "1"}},
		{"DeleteCluster OK", serv.DeleteCluster, http.StatusAccepted, "DELETE", false, requestData{"id": "1"}},
	}

	for _, tt := range nonErrorTT {
		testRequest(t, tt)
	}
}

func TestDatabaseErrorsCluster(t *testing.T) {
	serv := MockedIOCServer(t)
	dbErrorTT := []testCase{
		{"GetClusters DB error", serv.GetClusters, http.StatusInternalServerError, "GET", true, requestData{}},
		{"NewCluster DB error", serv.NewCluster, http.StatusInternalServerError, "POST", false, requestData{"name": "test"}},
		{"GetClusterByID DB error", serv.GetClusterByID, http.StatusInternalServerError, "GET", true, requestData{"id": "1"}},
		{"DeleteCluster DB error", serv.DeleteCluster, http.StatusInternalServerError, "DELETE", false, requestData{"id": "1"}},
		{"SearchCluster DB error", serv.SearchCluster, http.StatusInternalServerError, "GET", true, requestData{"name": "test"}},
	}

	// close DB
	serv.Storage.Close()

	for _, tt := range dbErrorTT {
		testRequest(t, tt)
	}
}

func TestParameterErrorsCluster(t *testing.T) {
	serv := MockedIOCServer(t)
	dbErrorTT := []testCase{
		{"NewCluster no param", serv.NewCluster, http.StatusBadRequest, "POST", false, requestData{}},
		{"NewCluster empty name key", serv.NewCluster, http.StatusBadRequest, "POST", false, requestData{"name": ""}},
		{"GetClusterByID no param", serv.GetClusterByID, http.StatusBadRequest, "GET", true, requestData{}},
		{"GetClusterByID non-str id", serv.GetClusterByID, http.StatusBadRequest, "GET", true, requestData{"id": "test"}},
		{"DeleteCluster no param", serv.DeleteCluster, http.StatusBadRequest, "DELETE", false, requestData{}},
		{"DeleteCluster by name", serv.DeleteCluster, http.StatusBadRequest, "DELETE", false, requestData{"name": "test"}},
		{"DeleteCluster by non-int id", serv.DeleteCluster, http.StatusBadRequest, "DELETE", false, requestData{"id": "test"}},
		{"SearchCluster wrong format", serv.SearchCluster, http.StatusBadRequest, "GET", true, requestData{"name": "test"}},
	}

	for _, tt := range dbErrorTT {
		testRequest(t, tt)
	}
}
