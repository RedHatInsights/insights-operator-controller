// Utils for working with maps

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
// https://redhatinsights.github.io/insights-operator-controller/packages/utils/maps.html

import "net/url"

import "strings"

// mergeMapsT merges multiple maps of various types, trying to convert them
// srcs can be;
// map[string]interface{} (for example body dedoded from json)
// map[string][]string  (for example r.URL.Query() or r.Form (as url.Values))
func mergeMapsT(srcs ...interface{}) map[string]interface{} {
	input := make(map[string]interface{})
	for _, srcm := range srcs {
		var toMerge map[string]interface{}

		switch v := srcm.(interface{}).(type) {
		case map[string]interface{}:
			toMerge = v

		case map[string][]string:
			toMerge = interMap(v)

		case url.Values:
			toMerge = interMap((map[string][]string)(v))

		default:
			panic("unknown type")
		}

		input = MergeMaps(input, toMerge)
	}
	return input
}

// stringsMap converts map[string]interface{} to flat map[string][]string
// Main use case is convert already flat map[string][]string or map[string]string
func stringsMap(m map[string]interface{}) map[string][]string {
	sm := make(map[string][]string)
	// add more types if needed
	for k, v := range m {
		if s, ok := v.(string); ok {
			sm[k] = []string{s}
		}
		if s, ok := v.([]string); ok {
			sm[k] = s
		}
		if s, ok := v.(map[string]interface{}); ok {
			fm := stringsMap(s)
			for kk, vv := range fm {
				sm[kk] = vv
			}
		}
	}
	return sm
}

// LowerCaseKeys returns new map with keys changed to lowercase
func LowerCaseKeys(m map[string]interface{}) map[string]interface{} {
	nm := map[string]interface{}{}

	for k, v := range m {
		nm[strings.ToLower(k)] = v
	}
	return nm
}

// interMap converts map from [string][]string to [string]interface{}
// taking last value from many in []string if more exists
func interMap(m map[string][]string) map[string]interface{} {
	r := make(map[string]interface{})
	for k, v := range m {
		// use last provided value
		if len(v) > 0 {
			r[k] = v[len(v)-1]
		}
	}
	return r
}

// MergeMaps Merges provided maps into one
func MergeMaps(maps ...map[string]interface{}) map[string]interface{} {
	r := make(map[string]interface{})

	for _, m := range maps {
		for k, v := range m {
			r[k] = v
		}
	}
	return r
}
