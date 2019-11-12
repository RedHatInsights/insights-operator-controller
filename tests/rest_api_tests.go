package main

import (
	"github.com/redhatinsighs/insights-operator-controller/tests/rest"
	"github.com/verdverm/frisby"
)

func main() {
	tests.ServerTests()
	tests.ClusterTests()
	tests.ConfigurationTests()
	tests.OperatorTests()
	tests.ProfileTests()
	tests.TriggerTests()
	frisby.Global.PrintReport()
}
