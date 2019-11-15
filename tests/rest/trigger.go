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

type Trigger struct {
	Id          int    `json:"id"`
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

	triggers := []Trigger{}
	text, err := f.Resp.Content()
	if err != nil {
		f.AddError(err.Error())
	} else {
		json.Unmarshal(text, &triggers)
	}
	return triggers
}

func TriggerTests() {
}
