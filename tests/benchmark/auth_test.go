// Copyright 2020, 2021, 2022 Red Hat, Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package testsbenchmark

// Documentation in literate-programming-style is available at:
// https://redhatinsights.github.io/insights-operator-controller/packages/tests/benchmark/auth_test.html

import (
	"os"
	"testing"

	tests "github.com/RedHatInsights/insights-operator-controller/tests/rest"
	"github.com/verdverm/frisby"
)

func BenchmarkMissingToken(b *testing.B) {
	f := frisby.Create("Check performance of missing authorization token").Get(tests.API_URL)
	for i := 0; i < b.N; i++ {
		f.Send()
	}
}

func BenchmarkMalformedToken(b *testing.B) {
	f := frisby.Create("Check performance of malformed authorization token").Get(tests.API_URL)
	f.SetHeader("Authorization", "Bearer abcdef1234")
	for i := 0; i < b.N; i++ {
		f.Send()
	}
}

func BenchmarkInvalidToken(b *testing.B) {
	f := frisby.Create("Check performance of invalid authorization token").Get(tests.API_URL)
	f.SetHeader("Authorization", "nonsense")
	for i := 0; i < b.N; i++ {
		f.Send()
	}
}

func BenchmarkSuccessfulAuth(b *testing.B) {
	const ldapToken string = "LDAP_TOKEN"
	f := frisby.Create("Check performance of valid authorization token").Get(tests.API_URL)
	if os.Getenv(ldapToken) == "" {
		f.AddError("Please provide LDAP_TOKEN env variable!")
	} else {
		f.SetHeader("Authorization", "Bearer "+os.Getenv(ldapToken))
		for i := 0; i < b.N; i++ {
			f.Send()
		}
	}
}
