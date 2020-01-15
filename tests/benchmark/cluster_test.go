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
