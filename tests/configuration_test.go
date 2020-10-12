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

// Documentation in literate-programming-style is available at:
// https://redhatinsights.github.io/insights-operator-controller/packages/tests/benchmark/configuration_test.html

import (
	"fmt"
	"testing"

	tests "github.com/RedHatInsights/insights-operator-controller/tests/rest"
	"github.com/verdverm/frisby"
)

func BenchmarkReadListOfConfigurations(b *testing.B) {
	f := frisby.Create("Benchmark Check list of configurations")
	f.Get(tests.API_URL + "/client/configuration")
	for i := 0; i < b.N; i++ {
		f.Send()
	}
}

func BenchmarkEnableConfigurations(b *testing.B) {
	f := frisby.Create("Benchmark Check that configuration can be enabled")
	for i := 0; i < b.N; i++ {
		f.Put(tests.API_URL + "client/configuration/" + fmt.Sprintf("%d", i) + "/enable")
		f.Send()
	}
}

func BenchmarkDisableConfigurations(b *testing.B) {
	f := frisby.Create("Benchmark Check that configuration can be disabled")
	for i := 0; i < b.N; i++ {
		f.Put(tests.API_URL + "client/configuration/" + fmt.Sprintf("%d", i) + "/disable")
		f.Send()
	}
}

func BenchmarkDeleteConfiguration(b *testing.B) {
	f := frisby.Create("Benchmark Check configuration delete")
	for i := 0; i < b.N; i++ {
		f.Delete(tests.API_URL + "client/configuration/" + fmt.Sprintf("%d", i+b.N))
		f.Send()
	}
}

func BenchmarkGetConfiguration(b *testing.B) {
	f := frisby.Create("Benchmark Check describing (reading) existing configuration")
	for i := 0; i < b.N; i++ {
		f.Get(tests.API_URL + "client/configuration/" + fmt.Sprintf("%d", i+b.N))
		f.Send()
	}
}
