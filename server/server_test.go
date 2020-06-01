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
		{"Main endpoint", serv.MainEndpoint, http.StatusOK, "GET", true, requestData{}, requestData{}, ""},
	}

	for _, tt := range nonErrorTT {
		testRequest(t, tt)
	}
}
