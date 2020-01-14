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
