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

import (
	"encoding/json"
	"fmt"
	"github.com/verdverm/frisby"
	"strings"
)

type ConfigurationProfile struct {
	Id            int    `json:"id"`
	Configuration string `json:"configuration"`
	ChangedAt     string `json:"changed_at"`
	ChangedBy     string `json:"changed_by"`
	Description   string `json:"description"`
}

func compareConfigurationProfiles(f *frisby.Frisby, profiles []ConfigurationProfile, expected []ConfigurationProfile) {
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

func readConfigurationProfileFromResponse(f *frisby.Frisby) ConfigurationProfile {
	profile := ConfigurationProfile{}
	text, err := f.Resp.Content()
	if err != nil {
		f.AddError(err.Error())
	} else {
		json.Unmarshal(text, &profile)
	}
	return profile
}

func readConfigurationProfilesFromResponse(f *frisby.Frisby) []ConfigurationProfile {
	profiles := []ConfigurationProfile{}
	text, err := f.Resp.Content()
	if err != nil {
		f.AddError(err.Error())
	} else {
		json.Unmarshal(text, &profiles)
	}
	return profiles
}

func checkInitialListOfConfigurationProfiles() {
	f := frisby.Create("Check list of configuration profiles")
	f.Get(API_URL + "/client/profile")
	f.Send()
	f.ExpectStatus(200)
	f.ExpectHeader("Content-Type", "application/json; charset=utf-8")
	f.ExpectJsonLength("", 4)

	profiles := readConfigurationProfilesFromResponse(f)
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
	f.ExpectStatus(400)

	f.PrintReport()
}

func checkChangeConfigurationProfile() {
	f := frisby.Create("Check changing configuration profile")
	f.Put(API_URL + "client/profile/3?username=foo&description=bar")
	f.Req.Body = strings.NewReader(`{"no_op":"Z", "watch":[]}`)

	f.Send()
	f.ExpectStatus(202)
	f.ExpectHeader("Content-Type", "application/json; charset=utf-8")
}

func checkListOfConfigurationProfilesWithUpdatedItem() {
	f := frisby.Create("Check list of configuration profiles with new item")
	f.Get(API_URL + "/client/profile")
	f.Send()
	f.ExpectStatus(200)
	f.ExpectHeader("Content-Type", "application/json; charset=utf-8")
	f.ExpectJsonLength("", 4)

	profiles := readConfigurationProfilesFromResponse(f)
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
	f := frisby.Create("Check changing configuration profile")
	f.Put(API_URL + "client/profile/35?username=foo&description=bar")
	f.Req.Body = strings.NewReader(`{"no_op":"Z", "watch":[]}`)

	f.Send()
	f.ExpectStatus(202)
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
	f := frisby.Create("Check deletion of configuration profile")
	f.Delete(API_URL + "client/profile/35?")

	f.Send()
	f.ExpectStatus(200)
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
	f.ExpectJsonLength("", 3)

	profiles := readConfigurationProfilesFromResponse(f)
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
	f.ExpectJsonLength("", 4)

	profiles := readConfigurationProfilesFromResponse(f)
	expected := []ConfigurationProfile{
		{0, `{"no_op":"X", "watch":["a","b","c"]}`, "2019-01-01T00:00:00Z", "tester", "cfg1"},
		{1, `{"no_op":"X", "watch":["a","b","c"]}`, "2019-01-01T00:00:00Z", "tester", "cfg2"},
		{2, `{"no_op":"X", "watch":["a","b","c"]}`, "2019-10-11T00:00:00Z", "tester", "cfg3"},
		{3, `{"no_op":"W", "watch":[]}`, "2019-10-11T00:00:00Z", "tester2", "description"},
	}
	compareConfigurationProfiles(f, profiles, expected)

	f.PrintReport()
}

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
