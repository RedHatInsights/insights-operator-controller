package testsbenchmark

import (
	"os"
	"testing"

	tests "github.com/redhatinsighs/insights-operator-controller/tests/rest"
	"github.com/verdverm/frisby"
)

const (
	contentType = "Content-Type"
	appJSON     = "application/json; charset=utf-8"
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
