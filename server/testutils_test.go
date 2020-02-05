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
	"io/ioutil"
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
	dbDriver    = "sqlite3"
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
	reqData          requestData
	urlData          requestData
}

func testRequest(t *testing.T, test testCase) {
	t.Run(test.testName, func(t *testing.T) {

		req, _ := http.NewRequest(test.requestMethod, "", bytes.NewBufferString("Test Body"))

		// set URL vars
		q := req.URL.Query()
		for key, value := range test.urlData {
			q.Add(key, value)
		}
		// encode the parameters to be URL-safe
		req.URL.RawQuery = q.Encode()

		// set mux vars
		req = mux.SetURLVars(req, test.reqData)

		// new RequestRecorder which satisfies ResponseWriter interface
		rr := httptest.NewRecorder()

		// call the handlerFunction
		test.fName(rr, req)

		CheckResponse(t, rr, test.expectedHeader, test.checkContentType)
	})
}

func runSQLiteScript(t *testing.T, path string) {
	script, err := os.Open(path)
	if err != nil {
		t.Fatalf("Unable to open %v", path)
	}
	defer script.Close()

	// sqlite3 test.db
	cmd := exec.Command(dbDriver, sqliteDB)

	var out, stderr bytes.Buffer
	// stdin for the command `sqlite3 dbname` since we can't use <
	cmd.Stdin = script
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err = cmd.Run()

	if err != nil {
		log.Fatalf("Error executing query. Command Output: %+v\n: %+v, %v", out.String(), stderr.String(), err)
	}
}

// MockedIOCServer returns an insights-operator-controller Server with disabled Splunk
// and a SQLite db for testing purposes
func MockedIOCServer(t *testing.T, mockData bool) *server.Server {
	splunk := logging.NewClient(false, emptyStr, emptyStr, emptyStr, emptyStr, emptyStr)

	db := MockedSQLite(t, mockData)

	s := server.Server{
		Address:  emptyStr, // not necessary since handlers are called directly
		UseHTTPS: false,
		Storage:  db,
		Splunk:   splunk,
		TLSCert:  emptyStr,
		TLSKey:   emptyStr,
	}

	s.ClusterQuery = storage.NewClusterQuery(s.Storage)

	return &s
}

func MockedSQLite(t *testing.T, mockData bool) storage.Storage {
	dbDriver := dbDriver
	storageSpecification := sqliteDB

	rmsqlite := exec.Command("rm", "-f", sqliteDB)
	rmsqlite.Run()

	db := storage.New(dbDriver, storageSpecification)

	runSQLiteScript(t, "../local_storage/schema_sqlite.sql")

	if mockData {
		runSQLiteScript(t, "../local_storage/test_data_sqlite.sql")
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

	result := rr.Result()
	body, _ := ioutil.ReadAll(result.Body)
	defer result.Body.Close()

	t.Log(string(body))
}
