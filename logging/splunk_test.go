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

package logging_test

// Documentation in literate-programming-style is available at:
// https://redhatinsights.github.io/insights-operator-controller/packages/logging/splunk_test.html

import (
	"testing"

	"github.com/RedHatInsights/insights-operator-controller/logging"
)

const (
	SplunkAddress = "address"
	SplunkToken   = "token"
	Source        = "source"
	SourceType    = "sourceType"
	Index         = "index"
)

// constructEnabledClient constructs enabled Splunk client
func constructEnabledClient() logging.Client {
	return logging.NewClient(true, SplunkAddress, SplunkToken, Source, SourceType, Index)
}

// constructDisabledClient constructs disabled Splunk client
func constructDisabledClient() logging.Client {
	return logging.NewClient(false, SplunkAddress, SplunkToken, Source, SourceType, Index)
}

// TestNewClientEnabled checks if it is possible to construct enabled Splunk client
func TestNewClientEnabled(t *testing.T) {
	c := constructEnabledClient()
	if c.ClientImpl == nil {
		t.Fatal("ClientImpl should not be nil for non enabled Splunk client")
	}
}

// TestNewClientDisabled checks if it is possible to construct disabled Splunk client
func TestNewClientDisabled(t *testing.T) {
	c := constructDisabledClient()
	if c.ClientImpl != nil {
		t.Fatal("ClientImpl should be nil for disabled Splunk client")
	}
}

// TestLogOperationForEnabledClient checks the Splunt.Log method
func TestLogOperationForEnabledClient(t *testing.T) {
	c := constructEnabledClient()
	err := c.Log("foo", "bar")
	if err == nil {
		t.Fatal("Error should be returned for enabled client with improper address")
	}
}

// TestLogOperationForDisabledClient checks the Splunt.Log method
func TestLogOperationForDisabledClient(t *testing.T) {
	c := constructDisabledClient()
	err := c.Log("foo", "bar")
	if err != nil {
		t.Fatal("Error should not be returned for disabled client")
	}
}

// TestLogActionOperationForEnabledClient checks the Splunt.LogAction method
func TestLogActionOperationForEnabledClient(t *testing.T) {
	c := constructEnabledClient()
	err := c.LogAction("foo", "bar", "description")
	if err == nil {
		t.Fatal("Error should be returned for enabled client with improper address")
	}
}

// TestLogActionOperationForDisabledClient checks the Splunt.LogAction method
func TestLogActionOperationForDisabledClient(t *testing.T) {
	c := constructDisabledClient()
	err := c.LogAction("foo", "bar", "description")
	if err != nil {
		t.Fatal("Error should not be returned for disabled client")
	}
}

// TestLogTriggerActionOperationForEnabledClient checks the Splunt.LogTriggerAction method
func TestLogTriggerActionOperationForEnabledClient(t *testing.T) {
	c := constructEnabledClient()
	err := c.LogTriggerAction("action", "user", "cluster", "trigger")
	if err == nil {
		t.Fatal("Error should be returned for enabled client with improper address")
	}
}

// TestLogTriggerActionOperationForDisabledClient checks the Splunt.LogTriggerAction method
func TestLogTriggerActionOperationForDisabledClient(t *testing.T) {
	c := constructDisabledClient()
	err := c.LogTriggerAction("action", "user", "cluster", "trigger")
	if err != nil {
		t.Fatal("Error should not be returned for disabled client")
	}
}

// TestLogWithTimeOperationForEnabledClient checks the Splunt.LogWithTime method
func TestLogWithTimeOperationForEnabledClient(t *testing.T) {
	c := constructEnabledClient()
	err := c.LogWithTime(123, "bar", "description")
	if err == nil {
		t.Fatal("Error should be returned for enabled client with improper address")
	}
}

// TestLogWithTimeOperationForDisabledClient checks the Splunt.LogWithTime method
func TestLogWithTimeOperationForDisabledClient(t *testing.T) {
	c := constructDisabledClient()
	err := c.LogWithTime(123, "bar", "description")
	if err != nil {
		t.Fatal("Error should not be returned for disabled client")
	}
}
