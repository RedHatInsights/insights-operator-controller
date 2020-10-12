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

package main

// Generated documentation is available at:
// https://godoc.org/github.com/RedHatInsights/insights-operator-controller/tests
//
// Documentation in literate-programming-style is available at:
// https://redhatinsights.github.io/insights-operator-controller/packages/tests/rest_api_tests.html

import (
	"github.com/RedHatInsights/insights-operator-controller/server"
	tests "github.com/RedHatInsights/insights-operator-controller/tests/rest"
	"github.com/verdverm/frisby"
)

func main() {
	if server.Environment != "production" {
		tests.ServerTests()
		tests.ClusterTests()
		tests.ConfigurationTests()
		tests.OperatorTests()
		tests.ProfileTests()
		tests.TriggerTests()
	} else {
		tests.AuthTests()
	}
	frisby.Global.PrintReport()
}
