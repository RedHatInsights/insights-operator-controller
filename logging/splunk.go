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

package logging

// Generated documentation is available at:
// https://godoc.org/github.com/RedHatInsights/insights-operator-controller/logging
//
// Documentation in literate-programming-style is available at:
// https://redhatinsights.github.io/insights-operator-controller/packages/logging/splunk.html

import (
	"github.com/ZachtimusPrime/Go-Splunk-HTTP/splunk"
)

// Client represents a Splunk client instance.
type Client struct {
	ClientImpl *splunk.Client
}

// NewClient creates a new instance of Splunk client.
func NewClient(enabled bool, address, token, source, sourceType, index string) Client {
	if enabled {
		url := address + "/services/collector/raw"
		splunkClient := splunk.NewClient(nil, url, token, source, sourceType, index)
		return Client{ClientImpl: splunkClient}
	}
	return Client{ClientImpl: nil}
}

// Log add a new message into the Splunk log.
func (client Client) Log(key, value string) error {
	if client.ClientImpl != nil {
		err := client.ClientImpl.Log(
			map[string]string{key: value})
		return err
	}
	return nil
}

// LogAction add a new message about performed action into the Splunk log.
func (client Client) LogAction(action, user, description string) error {
	if client.ClientImpl != nil {
		err := client.ClientImpl.Log(
			map[string]string{
				"action":      action,
				"user":        user,
				"description": description})
		return err
	}
	return nil
}

// LogTriggerAction add a new message about performed trigger-related action into the Splunk log.
func (client Client) LogTriggerAction(action, user, cluster, trigger string) error {
	if client.ClientImpl != nil {
		err := client.ClientImpl.Log(
			map[string]string{
				"action":  action,
				"user":    user,
				"cluster": cluster,
				"trigger": trigger})
		return err
	}
	return nil
}

// LogWithTime add a new message with timestamp into the Splunk log.
func (client Client) LogWithTime(time int64, key, value string) error {
	if client.ClientImpl != nil {
		err := client.ClientImpl.LogWithTime(
			time,
			map[string]string{key: value})
		return err
	}
	return nil
}
