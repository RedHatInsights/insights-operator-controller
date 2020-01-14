package testsbenchmark

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
