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

package utils_test

// Documentation in literate-programming-style is available at:
// https://redhatinsights.github.io/insights-operator-controller/packages/utils/maps_test.html

import (
	"net/url"
	"testing"

	"github.com/RedHatInsights/insights-operator-controller/utils"
)

// TestLowerCaseKeys test the utility function LowerCaseKeys
func TestLowerCaseKeys(t *testing.T) {
	inputMap := make(map[string]interface{})
	inputMap["FOO"] = "bar"
	outputMap := utils.LowerCaseKeys(inputMap)

	// check number of pairs stored in output map
	if len(outputMap) != 1 {
		t.Fatal("Invalid length of output map")
	}

	// check pairs in output map
	v, ok := outputMap["foo"]
	if !ok {
		t.Fatal("Key was not found")
	}
	if v != "bar" {
		t.Fatal("Improper value", v)
	}

	_, ok = outputMap["FOO"]
	if ok {
		t.Fatal("Old key is still there")
	}
}

// TestMergeMaps checks the utility function MergeMaps
func TestMergeMaps(t *testing.T) {
	// first map to be merged
	inputMap1 := make(map[string]interface{})
	inputMap1["foo"] = "foo"

	// second map to be merged
	inputMap2 := make(map[string]interface{})
	inputMap2["bar"] = "bar"

	// try to merge maps
	outputMap := utils.MergeMaps(inputMap1, inputMap2)

	// check number of pairs stored in output map
	if len(outputMap) != 2 {
		t.Fatal("Invalid length of output map")
	}

	// check for first key+value pair
	v1, ok := outputMap["foo"]
	if !ok {
		t.Fatal("Key was not found")
	}
	if v1 != "foo" {
		t.Fatal("Improper value", v1)
	}

	// check for second key+value pair
	v2, ok := outputMap["bar"]
	if !ok {
		t.Fatal("Key was not found")
	}
	if v2 != "bar" {
		t.Fatal("Improper value", v2)
	}
}

// TestMergeMapsTMapOfInterfaces checks the utility function MergeMapsT
func TestMergeMapsTMapOfInterfaces(t *testing.T) {
	// first map to be merged
	inputMap1 := make(map[string]interface{})
	inputMap1["foo"] = "foo"

	// second map to be merged
	inputMap2 := make(map[string]interface{})
	inputMap2["bar"] = "bar"

	// try to merge maps
	outputMap := utils.MergeMapsT(inputMap1, inputMap2)

	// check number of pairs stored in output map
	if len(outputMap) != 2 {
		t.Fatal("Invalid length of output map")
	}

	// check for first key+value pair
	v1, ok := outputMap["foo"]
	if !ok {
		t.Fatal("Key was not found")
	}
	if v1 != "foo" {
		t.Fatal("Improper value", v1)
	}

	// check for second key+value pair
	v2, ok := outputMap["bar"]
	if !ok {
		t.Fatal("Key was not found")
	}
	if v2 != "bar" {
		t.Fatal("Improper value", v2)
	}
}

// TestMergeMapsTMapOfStrings checks the utility function MergeMapsT
func TestMergeMapsTMapOfStrings(t *testing.T) {
	// first map to be merged
	inputMap1 := make(map[string][]string)
	inputMap1["foo"] = []string{"foo"}

	// second map to be merged
	inputMap2 := make(map[string][]string)
	inputMap2["bar"] = []string{"bar"}

	// try to merge maps
	outputMap := utils.MergeMapsT(inputMap1, inputMap2)

	// check number of pairs stored in output map
	if len(outputMap) != 2 {
		t.Fatal("Invalid length of output map")
	}

	// check for first key+value pair
	v1, ok := outputMap["foo"]
	if !ok {
		t.Fatal("Key was not found")
	}
	if v1 != "foo" {
		t.Fatal("Improper value", v1)
	}

	// check for second key+value pair
	v2, ok := outputMap["bar"]
	if !ok {
		t.Fatal("Key was not found")
	}
	if v2 != "bar" {
		t.Fatal("Improper value", v2)
	}
}

// TestMergeMapsTURLValues checks the utility function MergeMapsT
func TestMergeMapsTURLValues(t *testing.T) {
	// first map to be merged
	inputMap1 := url.Values{}
	inputMap1.Set("foo", "foo")

	// second map to be merged
	inputMap2 := url.Values{}
	inputMap2.Set("bar", "bar")

	// try to merge maps
	outputMap := utils.MergeMapsT(inputMap1, inputMap2)

	// check number of pairs stored in output map
	if len(outputMap) != 2 {
		t.Fatal("Invalid length of output map")
	}

	// check for first key+value pair
	v1, ok := outputMap["foo"]
	if !ok {
		t.Fatal("Key was not found")
	}
	if v1 != "foo" {
		t.Fatal("Improper value", v1)
	}

	// check for second key+value pair
	v2, ok := outputMap["bar"]
	if !ok {
		t.Fatal("Key was not found")
	}
	if v2 != "bar" {
		t.Fatal("Improper value", v2)
	}
}

// TestMergeMapsWrongType check if wrong types are handled properly by the utility function mergeMapsT
func TestMergeMapsWrongType(t *testing.T) {
	// check for panic()
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic as expected")
		}
	}()

	// try to merge values with incorrect data types
	utils.MergeMapsT("foo", "bar")
}

// TestStringsMapMapOfStrings check the utility function stringsMap
func TestStringsMapMapOfStrings(t *testing.T) {
	inputMap := make(map[string]interface{})
	inputMap["foo"] = "bar"

	// construct output map from input map
	outputMap := utils.StringsMap(inputMap)

	// check number of pairs stored in output map
	if len(outputMap) != 1 {
		t.Fatal("Invalid length of output map")
	}

	// check pairs in output map
	v, ok := outputMap["foo"]
	if !ok {
		t.Fatal("Key was not found")
	}
	if v[0] != "bar" {
		t.Fatal("Improper value", v)
	}
}

// TestStringsMapMapOfSlices check the utility function stringsMap
func TestStringsMapMapOfSlices(t *testing.T) {
	inputMap := make(map[string]interface{})
	inputMap["foo"] = []string{"bar"}

	// construct output map from input map
	outputMap := utils.StringsMap(inputMap)

	// check number of pairs stored in output map
	if len(outputMap) != 1 {
		t.Fatal("Invalid length of output map")
	}

	// check pairs in output map
	v, ok := outputMap["foo"]
	if !ok {
		t.Fatal("Key was not found")
	}
	if v[0] != "bar" {
		t.Fatal("Improper value", v)
	}
}

// TestStringsMapEmptyInterface check the utility function stringsMap
func TestStringsMapEmptyInterface(t *testing.T) {
	inputMap := make(map[string]interface{})
	inputMap["foo"] = make(map[string]interface{})

	// construct output map from input map
	outputMap := utils.StringsMap(inputMap)

	// check number of pairs stored in output map
	if len(outputMap) != 0 {
		t.Fatal("Invalid length of output map")
	}

	// check pairs in output map
	_, ok := outputMap["foo"]
	if ok {
		t.Fatal("Key was found, which is not expected")
	}
}

// TestStringsMapFilledInInterface check the utility function stringsMap
func TestStringsMapFilledInInterface(t *testing.T) {
	var inputMap = make(map[string]interface{})

	entry := make(map[string]interface{})
	entry["bar"] = "baz"
	inputMap["foo"] = entry

	// construct output map from input map
	outputMap := utils.StringsMap(inputMap)

	// check number of pairs stored in output map
	if len(outputMap) != 1 {
		t.Fatal("Invalid length of output map")
	}

	// check pairs in output map
	_, ok := outputMap["foo"]
	if ok {
		t.Fatal("Key was found, which is not expected")
	}

	// check pairs in output map
	_, ok = outputMap["bar"]
	if !ok {
		t.Fatal("Key was not found")
	}
}
