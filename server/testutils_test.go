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
	//"encoding/json"
	"github.com/RedHatInsights/insights-operator-controller/logging"
	"github.com/RedHatInsights/insights-operator-controller/server"
	"github.com/RedHatInsights/insights-operator-controller/storage"
	//"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"testing"

	"github.com/gorilla/mux"
)

const (
	contentType = "Content-Type"
	appJSON     = "application/json; charset=utf-8"
	emptyStr    = ""
	sqliteDB    = "test.db"
)

type handlerFunction func(writer http.ResponseWriter, request *http.Request)

type requestData map[string]string

type testCase struct {
	testName         string
	fName            handlerFunction
	expectedHeader   int
	requestMethod    string
	checkContentType bool
	requestData      requestData
}

func testRequest(t *testing.T, test testCase) {
	t.Run(test.testName, func(t *testing.T) {

		req, _ := http.NewRequest(test.requestMethod, "", nil)

		req = mux.SetURLVars(req, test.requestData)

		rr := httptest.NewRecorder()

		test.fName(rr, req) // call the handlerFunction

		CheckResponse(t, rr, test.expectedHeader, test.checkContentType)
	})
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
	dbDriver := "sqlite3"
	storageSpecification := sqliteDB

	rmsqlite := exec.Command("rm", "-f", sqliteDB)
	rmsqlite.Run()

	db := storage.New(dbDriver, storageSpecification)

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
func CheckResponse(t *testing.T, rr *httptest.ResponseRecorder, expectedStatusCode int, checkContentType bool) {
	if statusCode := rr.Code; statusCode != expectedStatusCode {
		t.Errorf("Expected status code %v, got %v", expectedStatusCode, statusCode)
	}

	if checkContentType {
		cType := rr.Header().Get(contentType)
		if cType != appJSON {
			t.Errorf("Unexpected content type. Expected %v, got %v", appJSON, cType)
		}
	}
	/*
		body, err := ioutil.ReadAll(res.Body)
		defer res.Body.Close()


			var expected map[string]interface{}
			err = json.NewDecoder(strings.NewReader(expectedBody)).Decode(&expected)
			if err != nil {
				t.Fatal(err)
			}


		var response map[string]interface{}
		err = json.Unmarshal(body, &response)
		if err != nil {
			t.Fatal(err)
		}

			if equal := reflect.DeepEqual(response, expected); !equal {
				t.Errorf("Expected response %v.", expected)
			}
		}
	*/
}
