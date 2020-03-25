// Copyright 2020 Red Hat, Inc
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

import (
	"testing"

	tests "github.com/RedHatInsights/insights-operator-controller/tests/rest"
	"github.com/verdverm/frisby"
)

func BenchmarkRestAPIEntryPoint(b *testing.B) {
	f := frisby.Create("Check performance of entry point to REST API").Get(tests.API_URL)
	for i := 0; i < b.N; i++ {
		f.Send()
	}
}

func BenchmarkNonExistentEntryPoint(b *testing.B) {
	f := frisby.Create("Check performance of entry point to REST API").Get(tests.API_URL + "foobar")
	for i := 0; i < b.N; i++ {
		f.Send()
	}
}

func BenchmarkWrongEntryPoint(b *testing.B) {
	f := frisby.Create("Check performance of wrong entry point to REST API").Get(tests.API_URL + "../")
	for i := 0; i < b.N; i++ {
		f.Send()
	}
}

func BenchmarkWrongPostForEntryPoint(b *testing.B) {
	f := frisby.Create("Check performance of entry point to REST API with wrong method: POST").Post(tests.API_URL)
	for i := 0; i < b.N; i++ {
		f.Send()
	}
}

func BenchmarkWrongPutForEntryPoint(b *testing.B) {
	f := frisby.Create("Check performance of entry point to REST API with wrong method: PUT").Put(tests.API_URL)
	for i := 0; i < b.N; i++ {
		f.Send()
	}
}

func BenchmarkWrongDeleteForEntryPoint(b *testing.B) {
	f := frisby.Create("Check performance of entry point to REST API with wrong method: DELETE").Delete(tests.API_URL)
	for i := 0; i < b.N; i++ {
		f.Send()
	}
}
