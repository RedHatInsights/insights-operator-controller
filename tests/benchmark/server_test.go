package testsbenchmark

import (
	"testing"

	tests "github.com/redhatinsighs/insights-operator-controller/tests/rest"
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
