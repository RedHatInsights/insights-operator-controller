// Utils for Http Rest Request validation and decoding

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

package utils

// Generated documentation is available at:
// https://godoc.org/github.com/RedHatInsights/insights-operator-controller/utils
//
// Documentation in literate-programming-style is available at:
// https://redhatinsights.github.io/insights-operator-controller/packages/utils/validation.html

import (
	"github.com/asaskevich/govalidator"
	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()

// DecodeValidRequest validates input maps (From Query.URL, or decoded Json Body) against template and returns typed structure
// srcs can be list of either map[string]interface{} or map[string][]string
func DecodeValidRequest(dst interface{}, temp map[string]interface{}, srcs ...interface{}) error {
	input := LowerCaseKeys(mergeMapsT(srcs...))
	sm := stringsMap(input)
	// add invalid key to force non keyed validator to run
	input[""] = ""
	// validate generic map to report correct type errors
	_, err := govalidator.ValidateMap(input, temp)
	if err != nil {
		return err
	}
	// convert valid map to type safe struct
	err = decoder.Decode(dst, sm)
	return err
}

// Pagination defines type safe Pagination request components
type Pagination struct {
	Limit  int `schema:"limit" valid:"type(int)~Limit has to be a number"`
	Offset int `schema:"offset" valid:"type(int)~Offset has to be a number"`
}

// PaginationTemplate contains validation for Pagination components
var PaginationTemplate = map[string]interface{}{
	"limit":  "int~Limit has to be a number",
	"offset": "int~Offset has to be a number",
}
