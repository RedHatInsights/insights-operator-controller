package testsbenchmark

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
	gen.InsertTriggerType("must-gather", "Triggers must-gather operation on selected cluster")
	gen.PopulateTrigger("must-gather")
}
