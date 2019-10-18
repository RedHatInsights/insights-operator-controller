package server

import (
	"bytes"
	"encoding/json"
	"github.com/redhatinsighs/insights-operator-controller/storage"
	"io/ioutil"
	"net/http"
	"testing"
)

func readListOfConfigurationProfiles(t *testing.T) []storage.ConfigurationProfile {
	response, err := http.Get(API_URL + "client/profile")
	if err != nil {
		t.Errorf("Communication error with the server %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected HTTP status 200 OK, got %d", response.StatusCode)
	}

	body, readErr := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	if readErr != nil {
		t.Errorf("Unable to read response body")
	}

	profiles := []storage.ConfigurationProfile{}
	err = json.Unmarshal(body, &profiles)
	if err != nil {
		t.Error(err)
	}
	return profiles
}

func compareConfigurationProfiles(t *testing.T, profiles []storage.ConfigurationProfile, expected []storage.ConfigurationProfile) {
	if len(profiles) != len(expected) {
		t.Errorf("%d configuration profiles are expected, but got %d", len(expected), len(profiles))
	}

	for i := 0; i < len(expected); i++ {
		// just ignore timestamp as we are going to test REST API w/o mocking the database
		profiles[i].ChangedAt = ""
		expected[i].ChangedAt = ""
		if profiles[i] != expected[i] {
			t.Errorf("Different profile info returned: %v != %v", profiles[i], expected[i])
		}
	}
}

func TestGetListOfConfigurationProfiles(t *testing.T) {
	profiles := readListOfConfigurationProfiles(t)

	expected := []storage.ConfigurationProfile{
		{0, `{"no_op":"X", "watch":["a","b","c"]}`, "2019-01-01T00:00:00Z", "tester", "cfg1"},
		{1, `{"no_op":"X", "watch":["a","b","c"]}`, "2019-01-01T00:00:00Z", "tester", "cfg2"},
		{2, `{"no_op":"X", "watch":["a","b","c"]}`, "2019-10-11T00:00:00Z", "tester", "cfg3"},
		{3, `{"no_op":"Y", "watch":["d","e"]}`, "2019-10-11T00:00:00Z", "tester", "cfg3"},
	}
	compareConfigurationProfiles(t, profiles, expected)
}

func TestGetConfigurationProfile(t *testing.T) {
	response, err := http.Get(API_URL + "client/profile/0")
	if err != nil {
		t.Errorf("Communication error with the server %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected HTTP status 200 OK, got %d", response.StatusCode)
	}

	body, readErr := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	if readErr != nil {
		t.Errorf("Unable to read response body")
	}

	profile := storage.ConfigurationProfile{}
	err = json.Unmarshal(body, &profile)
	if err != nil {
		t.Error(err)
	}

	expected := storage.ConfigurationProfile{0, `{"no_op":"X", "watch":["a","b","c"]}`, "2019-01-01T00:00:00Z", "tester", "cfg1"}
	if profile != expected {
		t.Errorf("Different profile is returned for profile ID = 0: %v != %v", profile, expected)
	}
}

func TestChangeConfigurationProfile(t *testing.T) {
	var client http.Client

	var jsonStr = []byte(`{"no_op":"Z", "watch":[]}`)
	request, err := http.NewRequest("PUT", API_URL+"client/profile/3?username=foo&description=bar", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Errorf("Error creating request %v", err)
	}

	response, err := client.Do(request)
	if err != nil {
		t.Errorf("Communication error with the server %v", err)
	}
	if response.StatusCode != http.StatusAccepted {
		t.Errorf("Expected HTTP status 202 Accepted, got %d", response.StatusCode)
	}

	profiles := readListOfConfigurationProfiles(t)

	expected := []storage.ConfigurationProfile{
		{0, `{"no_op":"X", "watch":["a","b","c"]}`, "2019-01-01T00:00:00Z", "tester", "cfg1"},
		{1, `{"no_op":"X", "watch":["a","b","c"]}`, "2019-01-01T00:00:00Z", "tester", "cfg2"},
		{2, `{"no_op":"X", "watch":["a","b","c"]}`, "2019-10-11T00:00:00Z", "tester", "cfg3"},
		{3, `{"no_op":"Z", "watch":[]}`, "2019-10-11T00:00:00Z", "foo", "bar"},
	}
	compareConfigurationProfiles(t, profiles, expected)
}

func TestDeleteConfigurationProfile(t *testing.T) {
	var client http.Client

	request, err := http.NewRequest("DELETE", API_URL+"client/profile/3", nil)
	if err != nil {
		t.Errorf("Error creating request %v", err)
	}

	response, err := client.Do(request)
	if err != nil {
		t.Errorf("Communication error with the server %v", err)
	}
	if response.StatusCode != http.StatusAccepted {
		t.Errorf("Expected HTTP status 202 Accepted, got %d", response.StatusCode)
	}

	profiles := readListOfConfigurationProfiles(t)

	expected := []storage.ConfigurationProfile{
		{0, `{"no_op":"X", "watch":["a","b","c"]}`, "2019-01-01T00:00:00Z", "tester", "cfg1"},
		{1, `{"no_op":"X", "watch":["a","b","c"]}`, "2019-01-01T00:00:00Z", "tester", "cfg2"},
		{2, `{"no_op":"X", "watch":["a","b","c"]}`, "2019-10-11T00:00:00Z", "tester", "cfg3"},
	}
	compareConfigurationProfiles(t, profiles, expected)
}

func TestCreateNewConfigurationProfile(t *testing.T) {
	var client http.Client

	var jsonStr = []byte(`{"no_op":"W", "watch":[]}`)
	request, err := http.NewRequest("POST", API_URL+"client/profile?username=tester2&description=description", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Errorf("Error creating request %v", err)
	}

	response, err := client.Do(request)
	if err != nil {
		t.Errorf("Communication error with the server %v", err)
	}
	if response.StatusCode != http.StatusCreated {
		t.Errorf("Expected HTTP status 201 Accepted, got %d", response.StatusCode)
	}

	profiles := readListOfConfigurationProfiles(t)

	expected := []storage.ConfigurationProfile{
		{0, `{"no_op":"X", "watch":["a","b","c"]}`, "2019-01-01T00:00:00Z", "tester", "cfg1"},
		{1, `{"no_op":"X", "watch":["a","b","c"]}`, "2019-01-01T00:00:00Z", "tester", "cfg2"},
		{2, `{"no_op":"X", "watch":["a","b","c"]}`, "2019-10-11T00:00:00Z", "tester", "cfg3"},
		{3, `{"no_op":"W", "watch":[]}`, "2019-10-11T00:00:00Z", "tester2", "description"},
	}
	compareConfigurationProfiles(t, profiles, expected)
}
