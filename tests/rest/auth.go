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

package tests

// Documentation in literate-programming-style is available at:
// https://redhatinsights.github.io/insights-operator-controller/packages/tests/rest/auth.html

import (
	"github.com/verdverm/frisby"
	"os"
)

const (
	contentType = "Content-Type"
	appJSON     = "application/json; charset=utf-8"
)

func checkMissingToken() {
	f := frisby.Create("Check missing authorization token").Get(API_URL)
	f.Send()
	f.ExpectStatus(403)
	f.ExpectHeader(contentType, appJSON)
	f.ExpectContent("Missing auth token")
	f.PrintReport()
}

func checkMalformedToken() {
	f := frisby.Create("Check malformed authorization token").Get(API_URL)
	f.SetHeader("Authorization", "Bearer abcdef1234")
	f.Send()
	f.ExpectStatus(403)
	f.ExpectHeader(contentType, appJSON)
	f.ExpectContent("Malformed authentication token")
	f.PrintReport()
}

func checkInvalidToken() {
	f := frisby.Create("Check invalid authorization token").Get(API_URL)
	f.SetHeader("Authorization", "nonsense")
	f.Send()
	f.ExpectStatus(403)
	f.ExpectHeader(contentType, appJSON)
	f.ExpectContent("Invalid/Malformed auth token")
	f.PrintReport()
}

func checkSuccessfulAuth() {
	const ldapToken string = "LDAP_TOKEN"
	f := frisby.Create("Check valid authorization token").Get(API_URL)
	if os.Getenv(ldapToken) == "" {
		f.AddError("Please provide LDAP_TOKEN env variable!")
	} else {
		f.SetHeader("Authorization", "Bearer "+os.Getenv(ldapToken))
		f.Send()
		f.ExpectHeader(contentType, appJSON)
		f.ExpectStatus(200)
		f.ExpectContent("Hello world!")
		f.PrintReport()
	}
}

// AuthTests - authorization related tests
func AuthTests() {
	checkMissingToken()
	checkMalformedToken()
	checkInvalidToken()
	checkSuccessfulAuth()
}
