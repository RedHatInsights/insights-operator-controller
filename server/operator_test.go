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
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
)

type OperatorConfiguration map[string]interface{}

func TestRegisterCluster(t *testing.T) {
	if !RunServiceTests {
		return
	}
	var client http.Client

	request, err := http.NewRequest("PUT", API_URL+"operator/register/cluster6", nil)
	if err != nil {
		t.Errorf("Error creating request %v", err)
	}

	response, err := client.Do(request)
	if err != nil {
		t.Errorf("Communication error with the server %v", err)
	}
	if response.StatusCode != http.StatusCreated {
		t.Errorf("Expected HTTP status 201 Created, got %d", response.StatusCode)
	}
}

func TestReadClusterConfiguration(t *testing.T) {
	if !RunServiceTests {
		return
	}
	var client http.Client

	request, err := http.NewRequest("GET", API_URL+"operator/configuration/cluster0", nil)
	if err != nil {
		t.Errorf("Error creating request %v", err)
	}

	response, err := client.Do(request)
	if err != nil {
		t.Errorf("Communication error with the server %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected HTTP status 200 OK, got %d", response.StatusCode)
	}

	body, readErr := ioutil.ReadAll(response.Body)

	defer response.Body.Close()

	if readErr != nil {
		t.Errorf("Unable to read response body")
	}

	var configuration OperatorConfiguration = make(map[string]interface{})
	err = json.Unmarshal(body, &configuration)
	if err != nil {
		t.Error(err)
	}

	var expected OperatorConfiguration = make(map[string]interface{})

	payload := []byte(`{
"no_op" : "X",
"watch" : ["a", "b", "c"]
}`)
	json.Unmarshal(payload, &expected)

	for key, _ := range expected {
		_, found := configuration[key]
		if !found {
			t.Error("Wrong configuration returned by the service: ", configuration)
			break
		}
	}
}
