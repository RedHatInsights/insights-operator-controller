// Utilities for testing to avoid code repetition. File not in ./utils because
// of cyclic imports
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
	"bytes"
	"encoding/json"
	"flag"
	"github.com/RedHatInsights/insights-operator-controller/logging"
	"github.com/RedHatInsights/insights-operator-controller/server"
	"github.com/RedHatInsights/insights-operator-controller/storage"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"testing"
)

const (
	expectedBody = `{"color_s":"blue","extra_data_m":{"param1":1,"param2":false}}`
	contentType  = "Content-Type"
	appJSON      = "application/json; charset=utf-8"
	emptyStr     = ""
	sqliteDB     = "test.db"
)

// MockedHTTPServer prepares new instance of testing HTTP server
func MockedHTTPServer(handler func(responseWriter http.ResponseWriter, request *http.Request)) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(handler))
}

// MockedIOCServer returns an insights-operator-controller Server with disabled Splunk
// and a SQLite db for testing purposes
func MockedIOCServer(t *testing.T) *server.Server {
	splunk := logging.NewClient(false, emptyStr, emptyStr, emptyStr, emptyStr, emptyStr)

	db := MockedSQLite(t)

	s := server.Server{
		Address:  emptyStr, // not necessary since handlers are called directly
		UseHTTPS: false,
		Storage:  db,
		Splunk:   splunk,
		TLSCert:  emptyStr,
		TLSKey:   emptyStr,
	}

	return &s
}

func MockedSQLite(t *testing.T) storage.Storage {
	dbDriver := flag.String("dbdriver", "sqlite3", "database driver specification")
	storageSpecification := flag.String("storage", sqliteDB, "storage specification")
	flag.Parse()

	rmsqlite := exec.Command("rm", "-f", sqliteDB)
	rmsqlite.Run()

	db := storage.New(*dbDriver, *storageSpecification)

	// run schema_sqlite.sql
	cmd := exec.Command("sqlite3", sqliteDB)

	schema, err := os.Open("../local_storage/schema_sqlite.sql")
	if err != nil {
		t.Fatalf("Unable to open schema_sqlite")
	}
	defer schema.Close()

	var out, stderr bytes.Buffer
	// stdin for the command `sqlite3 dbname` since we can't use <
	cmd.Stdin = schema
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err = cmd.Run()

	if err != nil {
		log.Fatalf("Error executing query. Command Output: %+v\n: %+v, %v", out.String(), stderr.String(), err)
	}

	return db
}

// CheckResponse ...
func CheckResponse(url string, expectedStatusCode int, checkPayload bool, t *testing.T) {
	res, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != expectedStatusCode {
		t.Errorf("Expected status code %v, got %v", expectedStatusCode, res.StatusCode)
	}

	if checkPayload {
		cType := res.Header.Get(contentType)
		if cType != appJSON {
			t.Errorf("Unexpected content type. Expected %v, got %v", appJSON, cType)
		}

		body, err := ioutil.ReadAll(res.Body)
		defer res.Body.Close()

		/*
			var expected map[string]interface{}
			err = json.NewDecoder(strings.NewReader(expectedBody)).Decode(&expected)
			if err != nil {
				t.Fatal(err)
			}
		*/

		var response map[string]interface{}
		err = json.Unmarshal(body, &response)
		if err != nil {
			t.Fatal(err)
		}
		/*
			if equal := reflect.DeepEqual(response, expected); !equal {
				t.Errorf("Expected response %v.", expected)
			}
		*/
	}
}
