/*
Copyright © 2019, 2020, 2021, 2022 Red Hat, Inc.

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

// Documentation in literate-programming-style is available at:
// https://redhatinsights.github.io/insights-operator-controller/packages/tests/rest/profile.html

import (
	"encoding/json"
	"fmt"
	"github.com/verdverm/frisby"
	"strings"
)

// ConfigurationProfile represents configuration profile record in the controller service.
//     ID: unique key
//     Configuration: a JSON structure stored in a string
//     ChangeAt: username of admin that created or updated the configuration
//     ChangeBy: timestamp of the last configuration change
//     Description: a string with any comment(s) about the configuration
type ConfigurationProfile struct {
	ID            int    `json:"id"`
	Configuration string `json:"configuration"`
	ChangedAt     string `json:"changed_at"`
	ChangedBy     string `json:"changed_by"`
	Description   string `json:"description"`
}

// ConfigurationProfilesResponse represents default response for configuration profile request
type ConfigurationProfilesResponse struct {
	Status   string                 `json:"status"`
	Profiles []ConfigurationProfile `json:"profiles"`
}

// ConfigurationProfileResponse represents default response for single configuration profile request
type ConfigurationProfileResponse struct {
	Status  string               `json:"status"`
	Profile ConfigurationProfile `json:"profile"`
}

func compareConfigurationProfiles(f *frisby.Frisby, profiles, expected []ConfigurationProfile) {
	if len(profiles) != len(expected) {
		f.AddError(fmt.Sprintf("%d configuration profiles are expected, but got %d", len(expected), len(profiles)))
	}

	for i := 0; i < len(expected); i++ {
		// just ignore timestamp as we are going to test REST API w/o mocking the database
		profiles[i].ChangedAt = ""
		expected[i].ChangedAt = ""
		if profiles[i] != expected[i] {
			f.AddError(fmt.Sprintf("Different profile info returned: %v != %v", profiles[i], expected[i]))
			fmt.Println(profiles[i].Configuration)
			fmt.Println(expected[i].Configuration)
		}
	}
}

func compareConfigurationProfilesWithoutID(f *frisby.Frisby, profiles, expected []ConfigurationProfile) {
	if len(profiles) != len(expected) {
		f.AddError(fmt.Sprintf("%d configuration profiles are expected, but got %d", len(expected), len(profiles)))
	}

	for i := 0; i < len(expected); i++ {
		// just ignore IDs
		profiles[i].ID = 0
		expected[i].ID = 0
		// just ignore timestamp as we are going to test REST API w/o mocking the database
		profiles[i].ChangedAt = ""
		expected[i].ChangedAt = ""
		if profiles[i] != expected[i] {
			f.AddError(fmt.Sprintf("Different profile info returned: %v != %v", profiles[i], expected[i]))
			fmt.Println(profiles[i].Configuration)
			fmt.Println(expected[i].Configuration)
		}
	}
}

// readConfigurationProfileFromResponse tries to read configuration profile from the HTTP server response
func readConfigurationProfileFromResponse(f *frisby.Frisby) ConfigurationProfile {
	// default return value from this function
	response := ConfigurationProfileResponse{}

	// try to read payload from response
	text, err := f.Resp.Content()
	if err != nil {
		f.AddError(err.Error())
	} else {
		// try to unmarshall response body and check if it's correct
		err = json.Unmarshal(text, &response)
		if err != nil {
			// any error needs to be recorded
			f.AddError(err.Error())
		}
	}
	return response.Profile
}

// readConfigurationProfilesFromResponse tries to read list of configuration profiles from the HTTP server response
func readConfigurationProfilesFromResponse(f *frisby.Frisby) []ConfigurationProfile {
	// default return value from this function
	response := ConfigurationProfilesResponse{}

	// try to read payload from response
	text, err := f.Resp.Content()
	if err != nil {
		f.AddError(err.Error())
	} else {
		// try to unmarshall response body and check if it's correct
		err = json.Unmarshal(text, &response)
		if err != nil {
			// any error needs to be recorded
			f.AddError(err.Error())
		}
	}
	return response.Profiles
}

func checkNumberOfProfiles(f *frisby.Frisby, profiles []ConfigurationProfile, expected int) {
	if len(profiles) != expected {
		f.AddError(fmt.Sprintf("Number of returned configuration profiles %d differs from expected number %d", len(profiles), expected))
	}
}

func checkInitialListOfConfigurationProfiles() {
	f := frisby.Create("Check list of configuration profiles")
	f.Get(API_URL + "/client/profile")
	f.Send()
	f.ExpectStatus(200)
	f.ExpectHeader("Content-Type", "application/json; charset=utf-8")

	profiles := readConfigurationProfilesFromResponse(f)
	checkNumberOfProfiles(f, profiles, 4)

	expected := []ConfigurationProfile{
		{0, `{"no_op":"X", "watch":["a","b","c"]}`, "2019-01-01T00:00:00Z", "tester", "cfg1"},
		{1, `{"no_op":"X", "watch":["a","b","c"]}`, "2019-01-01T00:00:00Z", "tester", "cfg2"},
		{2, `{"no_op":"X", "watch":["a","b","c"]}`, "2019-10-11T00:00:00Z", "tester", "cfg3"},
		{3, `{"no_op":"Y", "watch":["d","e"]}`, "2019-10-11T00:00:00Z", "tester", "cfg3"},
	}
	compareConfigurationProfiles(f, profiles, expected)

	f.PrintReport()
}

func checkGetExistingConfigurationProfile() {
	f := frisby.Create("Check getting configuration profile that exists")
	f.Get(API_URL + "/client/profile/0")
	f.Send()
	f.ExpectStatus(200)
	f.ExpectHeader("Content-Type", "application/json; charset=utf-8")

	profile := readConfigurationProfileFromResponse(f)
	expected := ConfigurationProfile{0, `{"no_op":"X", "watch":["a","b","c"]}`, "2019-01-01T00:00:00Z", "tester", "cfg1"}
	if profile != expected {
		f.AddError(fmt.Sprintf("Different profile is returned for profile ID = 0: %v != %v", profile, expected))
	}

	f.PrintReport()
}

func checkGetNonexistentConfigurationProfile() {
	f := frisby.Create("Check getting configuration profile that does not exist")
	f.Get(API_URL + "/client/profile/1234")
	f.Send()
	f.ExpectStatus(404)

	f.PrintReport()
}

func checkChangeConfigurationProfile() {
	f := frisby.Create("Check changing configuration profile")
	f.Put(API_URL + "client/profile/3?username=foo&description=bar")
	f.Req.Body = strings.NewReader(`{"no_op":"Z", "watch":[]}`)

	f.Send()
	f.ExpectStatus(200)
	f.ExpectHeader("Content-Type", "application/json; charset=utf-8")
}

func checkListOfConfigurationProfilesWithUpdatedItem() {
	f := frisby.Create("Check list of configuration profiles with new item")
	f.Get(API_URL + "/client/profile")
	f.Send()
	f.ExpectStatus(200)
	f.ExpectHeader("Content-Type", "application/json; charset=utf-8")

	profiles := readConfigurationProfilesFromResponse(f)
	checkNumberOfProfiles(f, profiles, 4)

	expected := []ConfigurationProfile{
		{0, `{"no_op":"X", "watch":["a","b","c"]}`, "2019-01-01T00:00:00Z", "tester", "cfg1"},
		{1, `{"no_op":"X", "watch":["a","b","c"]}`, "2019-01-01T00:00:00Z", "tester", "cfg2"},
		{2, `{"no_op":"X", "watch":["a","b","c"]}`, "2019-10-11T00:00:00Z", "tester", "cfg3"},
		{3, `{"no_op":"Z", "watch":[]}`, "2019-10-11T00:00:00Z", "foo", "bar"},
	}
	compareConfigurationProfiles(f, profiles, expected)

	f.PrintReport()
}

func checkChangeNonExistingConfigurationProfile() {
	f := frisby.Create("Check changing non existing configuration profile")
	f.Put(API_URL + "client/profile/35?username=foo&description=bar")
	f.Req.Body = strings.NewReader(`{"no_op":"Z", "watch":[]}`)

	f.Send()
	f.ExpectStatus(404)
	f.ExpectHeader("Content-Type", "application/json; charset=utf-8")
}

func checkDeleteConfigurationProfile() {
	f := frisby.Create("Check deletion of configuration profile")
	f.Delete(API_URL + "client/profile/3?")

	f.Send()
	f.ExpectStatus(200)
	f.ExpectHeader("Content-Type", "application/json; charset=utf-8")
}

func checkDeleteNonexistingConfigurationProfile() {
	f := frisby.Create("Check deletion of non existing configuration profile")
	f.Delete(API_URL + "client/profile/35?")

	f.Send()
	f.ExpectStatus(404)
	f.ExpectHeader("Content-Type", "application/json; charset=utf-8")
}

func checkCreateNewConfigurationProfile() {
	f := frisby.Create("Check creating new configuration profile")
	f.Post(API_URL + "client/profile?username=tester2&description=description")
	f.Req.Body = strings.NewReader(`{"no_op":"W", "watch":[]}`)

	f.Send()
	f.ExpectStatus(201)
	f.ExpectHeader("Content-Type", "application/json; charset=utf-8")
}

func checkListOfConfigurationProfilesWithDeletedItem() {
	f := frisby.Create("Check list of configuration profiles with deleted item")
	f.Get(API_URL + "/client/profile")
	f.Send()
	f.ExpectStatus(200)
	f.ExpectHeader("Content-Type", "application/json; charset=utf-8")

	profiles := readConfigurationProfilesFromResponse(f)
	checkNumberOfProfiles(f, profiles, 3)

	expected := []ConfigurationProfile{
		{0, `{"no_op":"X", "watch":["a","b","c"]}`, "2019-01-01T00:00:00Z", "tester", "cfg1"},
		{1, `{"no_op":"X", "watch":["a","b","c"]}`, "2019-01-01T00:00:00Z", "tester", "cfg2"},
		{2, `{"no_op":"X", "watch":["a","b","c"]}`, "2019-10-11T00:00:00Z", "tester", "cfg3"},
	}
	compareConfigurationProfiles(f, profiles, expected)

	f.PrintReport()
}

func checkListOfConfigurationProfilesWithAddedItem() {
	f := frisby.Create("Check list of configuration profiles with added item")
	f.Get(API_URL + "/client/profile")
	f.Send()
	f.ExpectStatus(200)
	f.ExpectHeader("Content-Type", "application/json; charset=utf-8")

	profiles := readConfigurationProfilesFromResponse(f)
	checkNumberOfProfiles(f, profiles, 4)

	expected := []ConfigurationProfile{
		{0, `{"no_op":"X", "watch":["a","b","c"]}`, "2019-01-01T00:00:00Z", "tester", "cfg1"},
		{1, `{"no_op":"X", "watch":["a","b","c"]}`, "2019-01-01T00:00:00Z", "tester", "cfg2"},
		{2, `{"no_op":"X", "watch":["a","b","c"]}`, "2019-10-11T00:00:00Z", "tester", "cfg3"},
		{3, `{"no_op":"W", "watch":[]}`, "2019-10-11T00:00:00Z", "tester2", "description"},
	}
	compareConfigurationProfilesWithoutID(f, profiles, expected)

	f.PrintReport()
}

// ProfileTests run all configuration profile-related REST API tests.
func ProfileTests() {
	checkInitialListOfConfigurationProfiles()
	checkGetExistingConfigurationProfile()
	checkGetNonexistentConfigurationProfile()

	checkChangeConfigurationProfile()
	checkListOfConfigurationProfilesWithUpdatedItem()

	checkChangeNonExistingConfigurationProfile()
	checkListOfConfigurationProfilesWithUpdatedItem()

	checkDeleteConfigurationProfile()
	checkListOfConfigurationProfilesWithDeletedItem()

	checkDeleteNonexistingConfigurationProfile()
	checkListOfConfigurationProfilesWithDeletedItem()

	checkCreateNewConfigurationProfile()
	checkListOfConfigurationProfilesWithAddedItem()
}
