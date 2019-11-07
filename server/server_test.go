/*
Copyright © 2019 Red Hat, Inc.

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

package server

import (
	"net/http"
	"testing"
)

func TestMainEndpoint(t *testing.T) {
	if !RunServiceTests {
		return
	}
	response, err := http.Get(API_URL)
	if err != nil {
		t.Errorf("Communication error with the server %v", err)
		return
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected HTTP status 200 OK, got %d", response.StatusCode)
		return
	}
}
