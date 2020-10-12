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

// Documentation in literate-programming-style is available at:
// https://redhatinsights.github.io/insights-operator-controller/packages/tests/benchmark/init_test.html

import (
	"os"
	"testing"

	testdata "github.com/RedHatInsights/insights-operator-controller/tests/setup"
)

const dbDriverEnv string = "DBDRIVER"
const storageSpecificationEnv = "STORAGE"

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	gen := testdata.NewDataGenerator(
		os.Getenv(dbDriverEnv), os.Getenv(storageSpecificationEnv))
	defer gen.Close()

	gen.PopulateCluster()
	gen.PopulateOperatorConfiguration()
	err := gen.InsertTriggerType("must-gather", "Triggers must-gather operation on selected cluster")
	if err != nil {
		panic(err)
	}
	gen.PopulateTrigger("must-gather")
}
