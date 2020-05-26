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

import (
	"testing"

	"github.com/asaskevich/govalidator"

	"github.com/RedHatInsights/insights-operator-controller/utils"
)

// Template is a template used by unit tests
var Template = utils.MergeMaps(map[string]interface{}{
	"id":   "int~Error reading and decoding cluster ID from query",
	"name": "",
	"":     "oneOfIdOrName~Either cluster ID or name needs to be specified",
}, utils.PaginationTemplate)

// Structure to hold deserialized input
type Destination struct {
	utils.Pagination
	ID   int    `schema:"id"`
	Name string `schema:"name"`
}

// TestDecodeValidRequestForImproperInput test the function DecodeValidRequest for invalid input
func TestDecodeValidRequestForImproperInput(t *testing.T) {
	var dst interface{}

	src := make(map[string]interface{})

	err := utils.DecodeValidRequest(&dst, Template, src)
	if err == nil {
		t.Fatal("error is expected for improper input")
	}

}

// oneOfIDOrNameValidation validates that id or name is filled
func oneOfIDOrNameValidation(i interface{}, context interface{}) bool {
	// Tag oneOfIdOrName
	v, ok := context.(map[string]interface{})
	if !ok {
		return false
	}
	// the int validation is done next by validator, we are just checking if its filled
	if id, ok := v["id"].(string); ok && len(id) != 0 {
		return true
	}
	if name, ok := v["name"].(string); ok && len(name) != 0 {
		return true
	}
	return false
}

// TestDecodeValidRequestForProperInput1 test the function DecodeValidRequest for valid input
func TestDecodeValidRequestForProperInput1(t *testing.T) {
	govalidator.CustomTypeTagMap.Set("oneOfIdOrName", govalidator.CustomTypeValidator(oneOfIDOrNameValidation))
	var dst Destination

	src := make(map[string]interface{})
	src["id"] = "42"
	src["limit"] = 100
	src["offset"] = 0

	err := utils.DecodeValidRequest(&dst, Template, src)
	if err != nil {
		t.Fatal("error is not expected for proper input", err)
	}
}

// TestDecodeValidRequestForProperInput2 test the function DecodeValidRequest for valid input
func TestDecodeValidRequestForProperInput2(t *testing.T) {
	govalidator.CustomTypeTagMap.Set("oneOfIdOrName", govalidator.CustomTypeValidator(oneOfIDOrNameValidation))
	var dst Destination

	src := make(map[string]interface{})
	src["name"] = "cluster name"
	src["limit"] = 100
	src["offset"] = 0

	err := utils.DecodeValidRequest(&dst, Template, src)
	if err != nil {
		t.Fatal("error is not expected for proper input", err)
	}
}
