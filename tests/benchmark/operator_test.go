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
	//"fmt"
	"testing"

	tests "github.com/RedHatInsights/insights-operator-controller/tests/rest"
	testdata "github.com/RedHatInsights/insights-operator-controller/tests/setup"
	"github.com/verdverm/frisby"
)

func BenchmarkConfigurationForCluster(b *testing.B) {
	f := frisby.Create("benchmark Check reading the configuration for clusters")
	for i := 0; i < b.N; i++ {
		f.Get(tests.API_URL + "/operator/configuration/" + testdata.GetClusterName(i))
		f.Send()
	}
}

func BenchmarkRegisterNewCluster(b *testing.B) {
	f := frisby.Create("Benchmark Check cluster registration")
	for i := 0; i < b.N; i++ {
		f.Put(tests.API_URL + "/operator/register/" + testdata.GetClusterName(i))
		f.Send()
	}
}
