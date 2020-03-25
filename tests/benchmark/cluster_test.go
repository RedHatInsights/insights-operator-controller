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
	"fmt"
	"testing"

	tests "github.com/RedHatInsights/insights-operator-controller/tests/rest"
	testdata "github.com/RedHatInsights/insights-operator-controller/tests/setup"
	"github.com/verdverm/frisby"
)

func BenchmarkReadListOfClusters(b *testing.B) {
	f := frisby.Create("Benchmark Check the initial list of clusters")
	f.Get(tests.API_URL + "/client/cluster")
	for i := 0; i < b.N; i++ {
		f.Send()
	}
}

func BenchmarkAddClusters(b *testing.B) {
	f := frisby.Create("Benchmark Check adding new cluster")
	for i := 0; i < b.N; i++ {
		f.Post(tests.API_URL + "client/cluster/" + testdata.GetClusterName(b.N+i))
		f.Send()
	}
}

func BenchmarkDeleteCluster(b *testing.B) {
	f := frisby.Create("BenchMark Check deleting existing cluster")
	for i := 0; i < b.N; i++ {
		f.Delete(tests.API_URL + "client/cluster/" + fmt.Sprintf("%d", b.N+i))
		f.Send()
	}
}
