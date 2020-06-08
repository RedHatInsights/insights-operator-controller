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
	"net/http/httptest"
	"testing"
	"time"

	"github.com/RedHatInsights/insights-operator-controller/logging"
	"github.com/RedHatInsights/insights-operator-controller/server"
	"github.com/RedHatInsights/insights-operator-controller/storage"

	"github.com/RedHatInsights/insights-operator-controller/tests/helpers"
)

// TestAddDefaultHeaders tests middleware adding headers
func TestAddDefaultHeaders(t *testing.T) {
	expectedHeaders := map[string]string{
		"Access-Control-Allow-Methods":     "POST, GET, OPTIONS, PUT, DELETE",
		"Access-Control-Allow-Headers":     "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token",
		"Access-Control-Allow-Credentials": "true",
		"Access-Control-Allow-Origin":      "local",
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headers := w.Header()
		for k, v := range expectedHeaders {
			if header := headers.Get(k); header != v {
				t.Errorf("Unexpected header value of %v. Expected %v, got %v", k, v, header)
			}
		}
	})

	serv := MockedIOCServer(t, false)
	defer serv.Storage.Close()

	// Add header to test addition of Access-Control-Allow-Origin
	req, _ := http.NewRequest("GET", "/health-check", nil)
	req.Header.Set("Origin", "local")

	rr := httptest.NewRecorder()

	// call the handler with the middleware
	// test call LogRequest too
	handl := serv.LogRequest(serv.AddDefaultHeaders(handler))
	handl.ServeHTTP(rr, req)
}

// TestMainEndpoint tests OK behaviour with empty DB (schema only)
func TestMainEndpoint(t *testing.T) {
	serv := MockedIOCServer(t, false)
	defer serv.Storage.Close()

	nonErrorTT := []testCase{
		{"Main endpoint", serv.MainEndpoint, http.StatusOK, "GET", false, requestData{}, requestData{}, ""},
	}

	for _, tt := range nonErrorTT {
		testRequest(t, tt)
	}
}

// TestServerInitialize check the initialization method
func TestServerInitialize(t *testing.T) {
	helpers.RunTestWithTimeout(t, func(t *testing.T) {
		splunk := logging.NewClient(false, "", "", "", "", "")

		storageInstance, err := storage.New("sqlite3", ":memory:")
		if err != nil {
			t.Fatal(err)
		}
		defer storageInstance.Close()

		serv := server.Server{
			Address:  "localhost:9999",
			UseHTTPS: false,
			Storage:  storageInstance,
			Splunk:   splunk,
			TLSCert:  "",
			TLSKey:   "",
		}

		serv.Initialize()
	}, 5*time.Second, false)
}

// TestServerInitializeOnProduction check the initialization method
func TestServerInitializeOnProduction(t *testing.T) {
	helpers.RunTestWithTimeout(t, func(t *testing.T) {
		environment := server.Environment
		defer func() {
			server.Environment = environment
		}()

		server.Environment = "production"

		splunk := logging.NewClient(false, "", "", "", "", "")

		storageInstance, err := storage.New("sqlite3", ":memory:")
		if err != nil {
			t.Fatal(err)
		}
		defer storageInstance.Close()

		serv := server.Server{
			Address:  "localhost:10000",
			UseHTTPS: false,
			Storage:  storageInstance,
			Splunk:   splunk,
			TLSCert:  "",
			TLSKey:   "",
		}

		serv.Initialize()
	}, 5*time.Second, false)
}
