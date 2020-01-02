/*
Copyright © 2019, 2020 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package tests

import "github.com/verdverm/frisby"

func checkConfigurationForCluster0() {
	f := frisby.Create("Check reading the configuration for cluster 00000000-0000-0000-0000-000000000000")
	f.Get(API_URL + "/operator/configuration/00000000-0000-0000-0000-000000000000")
	f.Send()
	f.ExpectStatus(200)
	f.ExpectHeader("Content-Type", "application/json; charset=utf-8")

	// check the content JSON response
	f.ExpectJson("configuration", "{\"no_op\":\"X\", \"watch\":[\"a\",\"b\",\"c\"]}")
	f.ExpectJson("status", "ok")
	f.PrintReport()
}

func checkRegisterNewCluster() {
	f := frisby.Create("Check if new cluster can be registered")
	f.Put(API_URL + "/operator/register/00000000-0000-0000-0000-000000000006")
	f.Send()
	f.ExpectStatus(201)
	f.ExpectHeader("Content-Type", "application/json; charset=utf-8")
}

func checkNonExistingConfiguration() {
	f := frisby.Create("Try to read configuration that does not exist")
	// configuration can't exists
	f.Get(API_URL + "/operator/configuration/00000000-0000-0000-0000-000000000006")
	f.Send()
	f.ExpectStatus(400)
	f.PrintReport()
}

// OperatorTests run all operator-related REST API tests.
func OperatorTests() {
	checkConfigurationForCluster0()
	// configuration for non-existing cluster
	checkNonExistingConfiguration()
	checkRegisterNewCluster()
	// configuration for newly created cluster
	checkNonExistingConfiguration()
	// TODO:
	// working with clusters with improper names
}
