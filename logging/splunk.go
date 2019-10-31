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
package logging

import (
	"github.com/ZachtimusPrime/Go-Splunk-HTTP/splunk"
)

type Client struct {
	ClientImpl *splunk.Client
}

func NewClient(address string, token string, source string, source_type string, index string) Client {
	url := address + "/services/collector/raw"
	splunk := splunk.NewClient(nil, url, token, source, source_type, index)
	return Client{ClientImpl: splunk}
}

func (client Client) Log(key string, value string) error {
	err := client.ClientImpl.Log(
		map[string]string{key: value})
	return err
}

func (client Client) LogAction(action string, user string, description string) error {
	err := client.ClientImpl.Log(
		map[string]string{
			"action":      action,
			"user":        user,
			"description": description})
	return err
}

func (client Client) LogWithTime(time int64, key string, value string) error {
	err := client.ClientImpl.LogWithTime(
		time,
		map[string]string{key: value})
	return err
}
