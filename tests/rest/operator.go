/*
Copyright © 2019 Red Hat, Inc.

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

type OperatorConfiguration map[string]interface{}

func checkConfigurationForCluster0() {
	f := frisby.Create("Check /operator/configuration/cluster0")
	f.Get(API_URL + "/operator/configuration/cluster0")
	f.Send()
	f.ExpectStatus(200)

	// check the content JSON response
	f.ExpectJson("no_op", "X")
	f.ExpectJson("watch.0", "a")
	f.ExpectJson("watch.1", "b")
	f.ExpectJson("watch.2", "c")
	f.PrintReport()
}

func checkRegisterNewCluster() {
	f := frisby.Create("Check if new cluster can be registered")
	f.Put(API_URL + "/operator/register/cluster6")
	f.Send()
	f.ExpectStatus(201)
}

func OperatorTests() {
	checkConfigurationForCluster0()
	checkRegisterNewCluster()
}
