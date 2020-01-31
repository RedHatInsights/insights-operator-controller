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
	"fmt"
	"github.com/RedHatInsights/insights-operator-controller/server"
	"net/http"
	"net/http/httptest"
	"testing"
	//"github.com/RedHatInsights/insights-operator-controller/utils"
)

type handlerFunction func(writer http.ResponseWriter, request *http.Request)

var nonErrorTT = []struct {
	testName         string
	fName            interface{}
	expectedHeader   int
	requestMethod    string
	checkContentType bool
}{
	{"GetClusters", server.GetClusters, http.StatusOK, "GET", true},
	{"NewCluster", server.NewCluster, http.StatusCreated, "POST", false},
	{"GetClusterByID", server.GetClusterByID, http.StatusAccepted, "GET", true},
	{"DeleteCluster", server.DeleteCluster, http.StatusAccepted, "DELETE", false},
	{"SearchCluster", server.SearchCluster, http.StatusOK, "GET", true},
}

func TestNonErrors(t *testing.T) {
	for _, tt := range nonErrorTT {
		t.Run(tt.testName, func(t *testing.T) {
			serv := MockedIOCServer(t)
			req, err := http.NewRequest(tt.requestMethod, "", nil)
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()

			fName := tt.fName
			serv.fName(rr, req) // eh, obviously doesn't work...
		})
	}
}

func TestGetClusters(t *testing.T) {
	/*
		server := MockedIOCServer(t)

		req, err := http.NewRequest("GET", "", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		server.GetClusters(rr, req)

	*/
	fmt.Println("hm")
}

/*
func TestNewCluster(t *testing.T) {

}
func TestGetClusterByID(t *testing.T) {

}
func TestDeleteCluster(t *testing.T) {

}
func TestSearchCluster(t *testing.T) {

}
*/
