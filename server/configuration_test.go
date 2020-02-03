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

func TestNonErrorsConfiguration(t *testing.T) {
	serv := MockedIOCServer(t)

	nonErrorTT := []testCase{
		{"GetConfiguration OK", serv.GetConfiguration, http.StatusOK, "GET", true, requestData{"id": "1"}},
		{"DeleteConfiguration OK", serv.DeleteConfiguration, http.StatusCreated, "DELETE", false, requestData{"name": "test"}},
		{"GetAllConfigurations OK", serv.GetAllConfigurations, http.StatusOK, "GET", true, requestData{}},
		{"GetClusterConfiguration OK", serv.GetClusterConfiguration, http.StatusAccepted, "GET", false, requestData{"id": "1"}},
		{"EnableConfiguration OK", serv.EnableConfiguration, http.StatusAccepted, "PUT", false, requestData{"id": "1"}},
		{"DisableConfiguration OK", serv.DisableConfiguration, http.StatusAccepted, "PUT", false, requestData{"id": "1"}},
	}

	for _, tt := range nonErrorTT {
		testRequest(t, tt)
	}
}

/*
func TestParameterErrorConfiguration(t *testing.T) {
	serv := MockedIOCServer(t)

	nonErrorTT := []testCase{
		{"GetConfiguration OK", serv.GetConfiguration, http.StatusOK, "GET", true, requestData{}},
		{"DeleteConfiguration OK", serv.DeleteConfiguration, http.StatusCreated, "POST", false, requestData{"name": "test"}},
		{"GetAllConfigurations OK", serv.GetAllConfigurations, http.StatusOK, "GET", true, requestData{"id": "1"}},
		{"GetClusterConfiguration OK", serv.GetClusterConfiguration, http.StatusAccepted, "DELETE", false, requestData{"id": "1"}},
		{"EnableConfiguration OK", serv.EnableConfiguration, http.StatusAccepted, "DELETE", false, requestData{"id": "1"}},
		{"DisableConfiguration OK", serv.DisableConfiguration, http.StatusAccepted, "DELETE", false, requestData{"id": "1"}},
	}

	for _, tt := range nonErrorTT {
		testRequest(t, tt)
	}
}
*/

/*
func TestSendConfiguration(t *testing.T) {
	payload := "enabled"
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sendConfiguration(w, payload)
	}))
	defer testServer.Close()
	res, err := http.Get(testServer.URL)
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyStr := string(bodyBytes)
	if bodyStr != payload {
		t.Fatal("SendConfiguration doesn't send correct data. Got %v, expected %v", bodyString, payload)
	}
}
*/
