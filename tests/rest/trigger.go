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
)

// Trigger represents trigger record in the controller service
//     ID: unique key
//     Type: ID of trigger type
//     Cluster: cluster ID (not name)
//     Reason: a string with any comment(s) about the trigger
//     Link: link to any document with customer ACK with the trigger
//     TriggeredAt: timestamp of the last configuration change
//     TriggeredBy: username of admin that created or updated the trigger
//     AckedAt: timestamp where the insights operator acked the trigger
//     Parameters: parameters that needs to be pass to trigger code
//     Active: flag indicating whether the trigger is still active or not
type Trigger struct {
	ID          int    `json:"id"`
	Type        string `json:"type"`
	Cluster     string `json:"cluster"`
	Reason      string `json:"reason"`
	Link        string `json:"link"`
	TriggeredAt string `json:"triggered_at"`
	TriggeredBy string `json:"triggered_by"`
	AckedAt     string `json:"acked_at"`
	Parameters  string `json:"parameters"`
	Active      int    `json:"active"`
}

// TriggerResponse represents default response for trigger request
type TriggerResponse struct {
	Status   string    `json:"status"`
	Triggers []Trigger `json:"triggers"`
}

func compareTriggers(f *frisby.Frisby, triggers []Trigger, expected []Trigger) {
	if len(triggers) != len(expected) {
		f.AddError(fmt.Sprintf("%d triggers are expected, but got %d", len(expected), len(triggers)))
	}

	for i := 0; i < len(expected); i++ {
		// just ignore timestamp as we are going to test REST API w/o mocking the database
		triggers[i].TriggeredAt = ""
		expected[i].TriggeredAt = ""
		if triggers[i] != expected[i] {
			f.AddError(fmt.Sprintf("Different trigger info returned: %v != %v", triggers[i], expected[i]))
			fmt.Println(triggers[i])
			fmt.Println(expected[i])
		}
	}
}

func readTriggers(f *frisby.Frisby) []Trigger {
	f.Get(API_URL + "client/trigger")
	f.Send()
	f.ExpectStatus(200)
	f.ExpectHeader("Content-Type", "application/json; charset=utf-8")

	response := TriggerResponse{}
	text, err := f.Resp.Content()
	if err != nil {
		f.AddError(err.Error())
	} else {
		json.Unmarshal(text, &response)
		fmt.Println("Work")
		fmt.Println(response.Triggers)
	}
	return response.Triggers
}

func checkInitialListOfTriggers() {
	f := frisby.Create("Check list of triggers")

	triggers := readTriggers(f)
	f.ExpectJsonLength("", 4)

	expected := []Trigger{
		{1, "must-gather", "00000000-0000-0000-0000-000000000000", "reason", "link", "1970-01-01T00:00:00Z", "tester", "1970-01-01T00:00:00Z", "{}", 0},
		{2, "must-gather", "00000000-0000-0000-0000-000000000000", "reason", "link", "1970-01-01T00:00:00Z", "tester", "1970-01-01T00:00:00Z", "{}", 1},
		{3, "must-gather", "00000000-0000-0000-0000-000000000001", "reason", "link", "1970-01-01T00:00:00Z", "tester", "1970-01-01T00:00:00Z", "{}", 0},
		{4, "must-gather", "00000000-0000-0000-0000-000000000001", "reason", "link", "1970-01-01T00:00:00Z", "tester", "1970-01-01T00:00:00Z", "{}", 1},
	}
	compareTriggers(f, triggers, expected)

	f.PrintReport()
}

func checkActivateExistingTrigger() {
	f := frisby.Create("Check activate existing trigger")

	triggers := readTriggers(f)
	f.ExpectJsonLength("", 4)

	expected := []Trigger{
		{1, "must-gather", "00000000-0000-0000-0000-000000000000", "reason", "link", "1970-01-01T00:00:00Z", "tester", "1970-01-01T00:00:00Z", "{}", 0},
		{2, "must-gather", "00000000-0000-0000-0000-000000000000", "reason", "link", "1970-01-01T00:00:00Z", "tester", "1970-01-01T00:00:00Z", "{}", 1},
		{3, "must-gather", "00000000-0000-0000-0000-000000000001", "reason", "link", "1970-01-01T00:00:00Z", "tester", "1970-01-01T00:00:00Z", "{}", 0},
		{4, "must-gather", "00000000-0000-0000-0000-000000000001", "reason", "link", "1970-01-01T00:00:00Z", "tester", "1970-01-01T00:00:00Z", "{}", 1},
	}
	compareTriggers(f, triggers, expected)

	f.Post(API_URL + "client/trigger/1/activate")
	f.Send()
	f.ExpectStatus(200)
	f.ExpectJson("status", "ok")

	triggers = readTriggers(f)
	f.ExpectJsonLength("", 4)

	expected = []Trigger{
		{1, "must-gather", "00000000-0000-0000-0000-000000000000", "reason", "link", "1970-01-01T00:00:00Z", "tester", "1970-01-01T00:00:00Z", "{}", 1},
		{2, "must-gather", "00000000-0000-0000-0000-000000000000", "reason", "link", "1970-01-01T00:00:00Z", "tester", "1970-01-01T00:00:00Z", "{}", 1},
		{3, "must-gather", "00000000-0000-0000-0000-000000000001", "reason", "link", "1970-01-01T00:00:00Z", "tester", "1970-01-01T00:00:00Z", "{}", 0},
		{4, "must-gather", "00000000-0000-0000-0000-000000000001", "reason", "link", "1970-01-01T00:00:00Z", "tester", "1970-01-01T00:00:00Z", "{}", 1},
	}
	compareTriggers(f, triggers, expected)

	f.PrintReport()
}

func checkDeactivateExistingTrigger() {
	f := frisby.Create("Check activate existing trigger")

	triggers := readTriggers(f)
	f.ExpectJsonLength("", 4)

	expected := []Trigger{
		{1, "must-gather", "00000000-0000-0000-0000-000000000000", "reason", "link", "1970-01-01T00:00:00Z", "tester", "1970-01-01T00:00:00Z", "{}", 1},
		{2, "must-gather", "00000000-0000-0000-0000-000000000000", "reason", "link", "1970-01-01T00:00:00Z", "tester", "1970-01-01T00:00:00Z", "{}", 1},
		{3, "must-gather", "00000000-0000-0000-0000-000000000001", "reason", "link", "1970-01-01T00:00:00Z", "tester", "1970-01-01T00:00:00Z", "{}", 0},
		{4, "must-gather", "00000000-0000-0000-0000-000000000001", "reason", "link", "1970-01-01T00:00:00Z", "tester", "1970-01-01T00:00:00Z", "{}", 1},
	}
	compareTriggers(f, triggers, expected)

	f.Post(API_URL + "client/trigger/1/deactivate")
	f.Send()
	f.ExpectStatus(200)
	f.ExpectJson("status", "ok")

	triggers = readTriggers(f)
	f.ExpectJsonLength("", 4)

	expected = []Trigger{
		{1, "must-gather", "00000000-0000-0000-0000-000000000000", "reason", "link", "1970-01-01T00:00:00Z", "tester", "1970-01-01T00:00:00Z", "{}", 0},
		{2, "must-gather", "00000000-0000-0000-0000-000000000000", "reason", "link", "1970-01-01T00:00:00Z", "tester", "1970-01-01T00:00:00Z", "{}", 1},
		{3, "must-gather", "00000000-0000-0000-0000-000000000001", "reason", "link", "1970-01-01T00:00:00Z", "tester", "1970-01-01T00:00:00Z", "{}", 0},
		{4, "must-gather", "00000000-0000-0000-0000-000000000001", "reason", "link", "1970-01-01T00:00:00Z", "tester", "1970-01-01T00:00:00Z", "{}", 1},
	}
	compareTriggers(f, triggers, expected)

	f.PrintReport()
}

func checkActivateAlreadyActivatedTrigger() {
	f := frisby.Create("Check activate already activated trigger")

	triggers := readTriggers(f)
	f.ExpectJsonLength("", 4)

	expected := []Trigger{
		{1, "must-gather", "00000000-0000-0000-0000-000000000000", "reason", "link", "1970-01-01T00:00:00Z", "tester", "1970-01-01T00:00:00Z", "{}", 0},
		{2, "must-gather", "00000000-0000-0000-0000-000000000000", "reason", "link", "1970-01-01T00:00:00Z", "tester", "1970-01-01T00:00:00Z", "{}", 1},
		{3, "must-gather", "00000000-0000-0000-0000-000000000001", "reason", "link", "1970-01-01T00:00:00Z", "tester", "1970-01-01T00:00:00Z", "{}", 0},
		{4, "must-gather", "00000000-0000-0000-0000-000000000001", "reason", "link", "1970-01-01T00:00:00Z", "tester", "1970-01-01T00:00:00Z", "{}", 1},
	}
	compareTriggers(f, triggers, expected)

	f.Post(API_URL + "client/trigger/2/activate")
	f.Send()
	f.ExpectStatus(200)
	f.ExpectJson("status", "ok")

	triggers = readTriggers(f)
	f.ExpectJsonLength("", 4)

	expected = []Trigger{
		{1, "must-gather", "00000000-0000-0000-0000-000000000000", "reason", "link", "1970-01-01T00:00:00Z", "tester", "1970-01-01T00:00:00Z", "{}", 0},
		{2, "must-gather", "00000000-0000-0000-0000-000000000000", "reason", "link", "1970-01-01T00:00:00Z", "tester", "1970-01-01T00:00:00Z", "{}", 1},
		{3, "must-gather", "00000000-0000-0000-0000-000000000001", "reason", "link", "1970-01-01T00:00:00Z", "tester", "1970-01-01T00:00:00Z", "{}", 0},
		{4, "must-gather", "00000000-0000-0000-0000-000000000001", "reason", "link", "1970-01-01T00:00:00Z", "tester", "1970-01-01T00:00:00Z", "{}", 1},
	}
	compareTriggers(f, triggers, expected)

	f.PrintReport()
}

func checkDeactivateAlreadyDeactivatedTrigger() {
	f := frisby.Create("Check deactivate already deactivated trigger")

	triggers := readTriggers(f)
	f.ExpectJsonLength("", 4)

	expected := []Trigger{
		{1, "must-gather", "00000000-0000-0000-0000-000000000000", "reason", "link", "1970-01-01T00:00:00Z", "tester", "1970-01-01T00:00:00Z", "{}", 0},
		{2, "must-gather", "00000000-0000-0000-0000-000000000000", "reason", "link", "1970-01-01T00:00:00Z", "tester", "1970-01-01T00:00:00Z", "{}", 1},
		{3, "must-gather", "00000000-0000-0000-0000-000000000001", "reason", "link", "1970-01-01T00:00:00Z", "tester", "1970-01-01T00:00:00Z", "{}", 0},
		{4, "must-gather", "00000000-0000-0000-0000-000000000001", "reason", "link", "1970-01-01T00:00:00Z", "tester", "1970-01-01T00:00:00Z", "{}", 1},
	}
	compareTriggers(f, triggers, expected)

	f.Post(API_URL + "client/trigger/3/deactivate")
	f.Send()
	f.ExpectStatus(200)
	f.ExpectJson("status", "ok")

	triggers = readTriggers(f)
	f.ExpectJsonLength("", 4)

	expected = []Trigger{
		{1, "must-gather", "00000000-0000-0000-0000-000000000000", "reason", "link", "1970-01-01T00:00:00Z", "tester", "1970-01-01T00:00:00Z", "{}", 0},
		{2, "must-gather", "00000000-0000-0000-0000-000000000000", "reason", "link", "1970-01-01T00:00:00Z", "tester", "1970-01-01T00:00:00Z", "{}", 1},
		{3, "must-gather", "00000000-0000-0000-0000-000000000001", "reason", "link", "1970-01-01T00:00:00Z", "tester", "1970-01-01T00:00:00Z", "{}", 0},
		{4, "must-gather", "00000000-0000-0000-0000-000000000001", "reason", "link", "1970-01-01T00:00:00Z", "tester", "1970-01-01T00:00:00Z", "{}", 1},
	}
	compareTriggers(f, triggers, expected)

	f.PrintReport()
}

func checkActivateNonExistingTrigger() {
	f := frisby.Create("Check activate non existing trigger")

	triggers := readTriggers(f)
	f.ExpectJsonLength("", 4)

	expected := []Trigger{
		{1, "must-gather", "00000000-0000-0000-0000-000000000000", "reason", "link", "1970-01-01T00:00:00Z", "tester", "1970-01-01T00:00:00Z", "{}", 0},
		{2, "must-gather", "00000000-0000-0000-0000-000000000000", "reason", "link", "1970-01-01T00:00:00Z", "tester", "1970-01-01T00:00:00Z", "{}", 1},
		{3, "must-gather", "00000000-0000-0000-0000-000000000001", "reason", "link", "1970-01-01T00:00:00Z", "tester", "1970-01-01T00:00:00Z", "{}", 0},
		{4, "must-gather", "00000000-0000-0000-0000-000000000001", "reason", "link", "1970-01-01T00:00:00Z", "tester", "1970-01-01T00:00:00Z", "{}", 1},
	}
	compareTriggers(f, triggers, expected)

	f.Post(API_URL + "client/trigger/42/activate")
	f.Send()
	f.ExpectStatus(200)

	triggers = readTriggers(f)
	f.ExpectJsonLength("", 4)

	expected = []Trigger{
		{1, "must-gather", "00000000-0000-0000-0000-000000000000", "reason", "link", "1970-01-01T00:00:00Z", "tester", "1970-01-01T00:00:00Z", "{}", 0},
		{2, "must-gather", "00000000-0000-0000-0000-000000000000", "reason", "link", "1970-01-01T00:00:00Z", "tester", "1970-01-01T00:00:00Z", "{}", 1},
		{3, "must-gather", "00000000-0000-0000-0000-000000000001", "reason", "link", "1970-01-01T00:00:00Z", "tester", "1970-01-01T00:00:00Z", "{}", 0},
		{4, "must-gather", "00000000-0000-0000-0000-000000000001", "reason", "link", "1970-01-01T00:00:00Z", "tester", "1970-01-01T00:00:00Z", "{}", 1},
	}
	compareTriggers(f, triggers, expected)

	f.PrintReport()
}

func checkDeactivateNonExistingTrigger() {
	f := frisby.Create("Check deactivate non existing trigger")

	triggers := readTriggers(f)
	f.ExpectJsonLength("", 4)

	expected := []Trigger{
		{1, "must-gather", "00000000-0000-0000-0000-000000000000", "reason", "link", "1970-01-01T00:00:00Z", "tester", "1970-01-01T00:00:00Z", "{}", 0},
		{2, "must-gather", "00000000-0000-0000-0000-000000000000", "reason", "link", "1970-01-01T00:00:00Z", "tester", "1970-01-01T00:00:00Z", "{}", 1},
		{3, "must-gather", "00000000-0000-0000-0000-000000000001", "reason", "link", "1970-01-01T00:00:00Z", "tester", "1970-01-01T00:00:00Z", "{}", 0},
		{4, "must-gather", "00000000-0000-0000-0000-000000000001", "reason", "link", "1970-01-01T00:00:00Z", "tester", "1970-01-01T00:00:00Z", "{}", 1},
	}
	compareTriggers(f, triggers, expected)

	f.Post(API_URL + "client/trigger/42/deactivate")
	f.Send()
	f.ExpectStatus(200)

	triggers = readTriggers(f)
	f.ExpectJsonLength("", 4)

	expected = []Trigger{
		{1, "must-gather", "00000000-0000-0000-0000-000000000000", "reason", "link", "1970-01-01T00:00:00Z", "tester", "1970-01-01T00:00:00Z", "{}", 0},
		{2, "must-gather", "00000000-0000-0000-0000-000000000000", "reason", "link", "1970-01-01T00:00:00Z", "tester", "1970-01-01T00:00:00Z", "{}", 1},
		{3, "must-gather", "00000000-0000-0000-0000-000000000001", "reason", "link", "1970-01-01T00:00:00Z", "tester", "1970-01-01T00:00:00Z", "{}", 0},
		{4, "must-gather", "00000000-0000-0000-0000-000000000001", "reason", "link", "1970-01-01T00:00:00Z", "tester", "1970-01-01T00:00:00Z", "{}", 1},
	}
	compareTriggers(f, triggers, expected)

	f.PrintReport()
}

func checkDeleteExistingTrigger() {
	f := frisby.Create("Check delete existing trigger")

	triggers := readTriggers(f)
	f.ExpectJsonLength("", 4)

	expected := []Trigger{
		{1, "must-gather", "00000000-0000-0000-0000-000000000000", "reason", "link", "1970-01-01T00:00:00Z", "tester", "1970-01-01T00:00:00Z", "{}", 0},
		{2, "must-gather", "00000000-0000-0000-0000-000000000000", "reason", "link", "1970-01-01T00:00:00Z", "tester", "1970-01-01T00:00:00Z", "{}", 1},
		{3, "must-gather", "00000000-0000-0000-0000-000000000001", "reason", "link", "1970-01-01T00:00:00Z", "tester", "1970-01-01T00:00:00Z", "{}", 0},
		{4, "must-gather", "00000000-0000-0000-0000-000000000001", "reason", "link", "1970-01-01T00:00:00Z", "tester", "1970-01-01T00:00:00Z", "{}", 1},
	}
	compareTriggers(f, triggers, expected)

	f.Delete(API_URL + "client/trigger/1")
	f.Send()
	f.ExpectStatus(200)

	triggers = readTriggers(f)
	f.ExpectJsonLength("", 3)

	expected = []Trigger{
		{2, "must-gather", "00000000-0000-0000-0000-000000000000", "reason", "link", "1970-01-01T00:00:00Z", "tester", "1970-01-01T00:00:00Z", "{}", 1},
		{3, "must-gather", "00000000-0000-0000-0000-000000000001", "reason", "link", "1970-01-01T00:00:00Z", "tester", "1970-01-01T00:00:00Z", "{}", 0},
		{4, "must-gather", "00000000-0000-0000-0000-000000000001", "reason", "link", "1970-01-01T00:00:00Z", "tester", "1970-01-01T00:00:00Z", "{}", 1},
	}
	compareTriggers(f, triggers, expected)

	f.PrintReport()
}

func checkDeleteNonExistingTrigger() {
	f := frisby.Create("Check delete non-existing trigger")

	triggers := readTriggers(f)
	f.ExpectJsonLength("", 3)

	expected := []Trigger{
		{2, "must-gather", "00000000-0000-0000-0000-000000000000", "reason", "link", "1970-01-01T00:00:00Z", "tester", "1970-01-01T00:00:00Z", "{}", 1},
		{3, "must-gather", "00000000-0000-0000-0000-000000000001", "reason", "link", "1970-01-01T00:00:00Z", "tester", "1970-01-01T00:00:00Z", "{}", 0},
		{4, "must-gather", "00000000-0000-0000-0000-000000000001", "reason", "link", "1970-01-01T00:00:00Z", "tester", "1970-01-01T00:00:00Z", "{}", 1},
	}
	compareTriggers(f, triggers, expected)

	f.Delete(API_URL + "client/trigger/42")
	f.Send()
	f.ExpectStatus(200)

	triggers = readTriggers(f)
	f.ExpectJsonLength("", 3)

	expected = []Trigger{
		{2, "must-gather", "00000000-0000-0000-0000-000000000000", "reason", "link", "1970-01-01T00:00:00Z", "tester", "1970-01-01T00:00:00Z", "{}", 1},
		{3, "must-gather", "00000000-0000-0000-0000-000000000001", "reason", "link", "1970-01-01T00:00:00Z", "tester", "1970-01-01T00:00:00Z", "{}", 0},
		{4, "must-gather", "00000000-0000-0000-0000-000000000001", "reason", "link", "1970-01-01T00:00:00Z", "tester", "1970-01-01T00:00:00Z", "{}", 1},
	}
	compareTriggers(f, triggers, expected)

	f.PrintReport()
}

func checkGetTriggersForCluster0() {
	f := frisby.Create("Check get trigger for a cluster")
	f.Get(API_URL + "client/cluster/00000000-0000-0000-0000-000000000000/trigger")
	f.Send()
	f.ExpectStatus(200)
	f.ExpectHeader("Content-Type", "application/json; charset=utf-8")

	triggers := TriggerResponse{}
	text, err := f.Resp.Content()
	if err != nil {
		f.AddError(err.Error())
	} else {
		json.Unmarshal(text, &triggers)
	}

	expected := []Trigger{
		{2, "must-gather", "00000000-0000-0000-0000-000000000000", "reason", "link", "1970-01-01T00:00:00Z", "tester", "1970-01-01T00:00:00Z", "{}", 1},
	}
	compareTriggers(f, triggers.Triggers, expected)

	f.PrintReport()
}

func checkGetTriggersForCluster1() {
	f := frisby.Create("Check get trigger for a cluster")
	f.Get(API_URL + "client/cluster/00000000-0000-0000-0000-000000000001/trigger")
	f.Send()
	f.ExpectStatus(200)
	f.ExpectHeader("Content-Type", "application/json; charset=utf-8")

	triggers := TriggerResponse{}
	text, err := f.Resp.Content()
	if err != nil {
		f.AddError(err.Error())
	} else {
		json.Unmarshal(text, &triggers)
	}

	expected := []Trigger{
		{3, "must-gather", "00000000-0000-0000-0000-000000000001", "reason", "link", "1970-01-01T00:00:00Z", "tester", "1970-01-01T00:00:00Z", "{}", 0},
		{4, "must-gather", "00000000-0000-0000-0000-000000000001", "reason", "link", "1970-01-01T00:00:00Z", "tester", "1970-01-01T00:00:00Z", "{}", 1},
	}
	compareTriggers(f, triggers.Triggers, expected)

	f.PrintReport()
}

func checkGetTriggersForClusterX() {
	f := frisby.Create("Check get trigger for a cluster")
	f.Get(API_URL + "client/cluster/00000000-ffff-0000-0000-000000000001/trigger")
	f.Send()
	f.ExpectStatus(200)
	f.ExpectHeader("Content-Type", "application/json; charset=utf-8")

	triggers := TriggerResponse{}
	text, err := f.Resp.Content()
	if err != nil {
		f.AddError(err.Error())
	} else {
		json.Unmarshal(text, &triggers)
	}

	expected := []Trigger{}
	compareTriggers(f, triggers.Triggers, expected)

	f.PrintReport()
}

func checkCreateNewTrigger() {
	f := frisby.Create("Check create new trigger")
	f.Post(API_URL + "client/cluster/00000000-0000-0000-0000-000000000001/trigger/must-gather?username=tester&reason=r&link=l")
	f.Send()
	f.ExpectStatus(200)

	triggers := readTriggers(f)
	f.ExpectJsonLength("", 4)

	expected := []Trigger{
		{2, "must-gather", "00000000-0000-0000-0000-000000000000", "reason", "link", "1970-01-01T00:00:00Z", "tester", "1970-01-01T00:00:00Z", "{}", 1},
		{3, "must-gather", "00000000-0000-0000-0000-000000000001", "reason", "link", "1970-01-01T00:00:00Z", "tester", "1970-01-01T00:00:00Z", "{}", 0},
		{4, "must-gather", "00000000-0000-0000-0000-000000000001", "reason", "link", "1970-01-01T00:00:00Z", "tester", "1970-01-01T00:00:00Z", "{}", 1},
		{5, "must-gather", "00000000-0000-0000-0000-000000000001", "r", "l", "1970-01-01T00:00:00Z", "tester", "1970-01-01T00:00:00Z", "", 1},
	}

	compareTriggers(f, triggers, expected)
	f.PrintReport()
}

func checkCreateNewTriggerForWrongCluster() {
	f := frisby.Create("Check create new trigger for wrong cluster")
	f.Post(API_URL + "client/cluster/00000000-0000-ffff-0000-000000000001/trigger/must-gather?username=tester&reason=r&link=l")
	f.Send()
	f.ExpectStatus(400)

	triggers := readTriggers(f)
	f.ExpectJsonLength("", 4)

	expected := []Trigger{
		{2, "must-gather", "00000000-0000-0000-0000-000000000000", "reason", "link", "1970-01-01T00:00:00Z", "tester", "1970-01-01T00:00:00Z", "{}", 1},
		{3, "must-gather", "00000000-0000-0000-0000-000000000001", "reason", "link", "1970-01-01T00:00:00Z", "tester", "1970-01-01T00:00:00Z", "{}", 0},
		{4, "must-gather", "00000000-0000-0000-0000-000000000001", "reason", "link", "1970-01-01T00:00:00Z", "tester", "1970-01-01T00:00:00Z", "{}", 1},
		{5, "must-gather", "00000000-0000-0000-0000-000000000001", "r", "l", "1970-01-01T00:00:00Z", "tester", "1970-01-01T00:00:00Z", "", 1},
	}
	compareTriggers(f, triggers, expected)
	f.PrintReport()
}

// TriggerTests run all trigger-related REST API tests.
func TriggerTests() {
	checkInitialListOfTriggers()

	// trigger activation and deactivation - basic cases
	checkActivateExistingTrigger()
	checkDeactivateExistingTrigger()

	// trigger activation and deactivation - already activated and deactivated triggers
	checkActivateAlreadyActivatedTrigger()
	checkDeactivateAlreadyDeactivatedTrigger()

	// trigger activation and deactivation - non-existing triggers
	checkActivateNonExistingTrigger()
	checkDeactivateNonExistingTrigger()

	checkDeleteExistingTrigger()
	checkDeleteNonExistingTrigger()

	checkGetTriggersForCluster0()
	checkGetTriggersForCluster1()
	checkGetTriggersForClusterX()

	checkCreateNewTrigger()
	checkCreateNewTriggerForWrongCluster()
}
