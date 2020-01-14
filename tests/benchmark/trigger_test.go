package testsbenchmark

import (
	"fmt"
	"testing"

	tests "github.com/redhatinsighs/insights-operator-controller/tests/rest"
	testdata "github.com/redhatinsighs/insights-operator-controller/tests/setup"
	"github.com/verdverm/frisby"
)

func BenchmarkReadTriggers(b *testing.B) {
	f := frisby.Create("benchmark read triggers")
	f.Get(tests.API_URL + "client/trigger")
	for i := 0; i < b.N; i++ {
		f.Send()
	}
}

func BenchmarkActivateTrigger(b *testing.B) {
	f := frisby.Create("benchmark activate triggers")
	for i := 0; i < b.N; i++ {
		f.Post(tests.API_URL + "client/trigger/" + fmt.Sprintf("%d", i) + "/activate")
		f.Send()
	}
}

func BenchmarkDeactivateTrigger(b *testing.B) {
	f := frisby.Create("benchmark deactivate triggers")
	for i := 0; i < b.N; i++ {
		f.Post(tests.API_URL + "client/trigger/" + fmt.Sprintf("%d", i) + "/deactivate")
		f.Send()
	}
}

func BenchmarkDeleteTrigger(b *testing.B) {
	f := frisby.Create("benchmark delete triggers")
	for i := 0; i < b.N; i++ {
		f.Delete(tests.API_URL + "client/trigger/" + fmt.Sprintf("%d", i+b.N))
		f.Send()
	}
}

func BenchmarkGetTriggerForCluster(b *testing.B) {
	f := frisby.Create("Benchmark get trigger for a cluster")
	for i := 0; i < b.N; i++ {
		f.Get(tests.API_URL + "client/cluster/" + testdata.GetClusterName(i) + "/trigger")
		f.Send()
	}
}

func BenchmarkCreateNewTrigger(b *testing.B) {
	f := frisby.Create("Check create new trigger")
	for i := 0; i < b.N; i++ {
		f.Post(tests.API_URL + "client/cluster/" + testdata.GetClusterName(i) + "/trigger/must-gather?username=tester&reason=r&link=l")
		f.Send()
	}

}
