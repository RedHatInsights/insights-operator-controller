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

package server

// Generated documentation is available at:
// https://godoc.org/github.com/RedHatInsights/insights-operator-controller/server
//
// Documentation in literate-programming-style is available at:
// https://redhatinsights.github.io/insights-operator-controller/packages/server/profile.html

import (
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/RedHatInsights/insights-operator-controller/storage"
	"github.com/RedHatInsights/insights-operator-utils/responses"
)

// ListConfigurationProfiles method reads list of configuration profiles.
func (s Server) ListConfigurationProfiles(writer http.ResponseWriter, request *http.Request) {
	// try to read list of configuration profiles from storage
	profiles, err := s.Storage.ListConfigurationProfiles()

	// check if the storage operation was successful
	if err == nil {
		responses.SendOK(writer, responses.BuildOkResponseWithData("profiles", profiles))
	} else {
		responses.SendInternalServerError(writer, err.Error())
	}
}

// GetConfigurationProfile method reads profile specified by its ID
func (s Server) GetConfigurationProfile(writer http.ResponseWriter, request *http.Request) {
	// profile ID needs to be specified in request
	id, err := retrieveIDRequestParameter(request)
	if err != nil {
		responses.SendBadRequest(writer, "Error reading profile ID from request\n")
		return
	}

	// try to read configuration for profile specified by its ID
	profile, err := s.Storage.GetConfigurationProfile(int(id))

	// check if the storage operation was successful
	if _, ok := err.(*storage.ItemNotFoundError); ok {
		responses.Send(http.StatusNotFound, writer, err.Error())
	} else if err != nil {
		responses.SendInternalServerError(writer, err.Error())
	} else {
		responses.SendOK(writer, responses.BuildOkResponseWithData("profile", profile))
	}
}

// NewConfigurationProfile method creates new configuration profile
func (s Server) NewConfigurationProfile(writer http.ResponseWriter, request *http.Request) {
	// username needs to be specified in request
	username, foundUsername := request.URL.Query()["username"]

	// description needs to be specified in request
	description, foundDescription := request.URL.Query()["description"]

	if !foundUsername {
		responses.SendBadRequest(writer, "User name needs to be specified\n")
		return
	}

	if !foundDescription {
		responses.SendBadRequest(writer, "Description needs to be specified\n")
		return
	}

	// read configuration from request body
	configuration, err := ioutil.ReadAll(request.Body)

	if err != nil || len(configuration) == 0 {
		responses.SendBadRequest(writer, "Configuration needs to be provided in the request body")
		return
	}

	// try to record the action NewConfigurationProfile into Splunk
	err = s.Splunk.LogAction("NewConfigurationProfile", username[0], string(configuration))
	// and check whether the Splunk operation was successful
	checkSplunkOperation(err)

	// try to store configuration profile into storage
	profiles, err := s.Storage.StoreConfigurationProfile(username[0], description[0], string(configuration))

	// check if the storage operation was successful
	if err != nil {
		responses.SendInternalServerError(writer, err.Error())
	} else {
		responses.SendCreated(writer, responses.BuildOkResponseWithData("profiles", profiles))
	}
}

// DeleteConfigurationProfile method deletes configuration profile
func (s Server) DeleteConfigurationProfile(writer http.ResponseWriter, request *http.Request) {
	// profile ID needs to be specified in request
	id, err := retrieveIDRequestParameter(request)
	if err != nil {
		responses.SendBadRequest(writer, "Error reading profile ID from request\n")
		return
	}

	// try to record the action DeleteConfigurationProfile into Splunk
	err = s.Splunk.LogAction("DeleteConfigurationProfile", "tester", strconv.Itoa(int(id)))
	// and check whether the Splunk operation was successful
	checkSplunkOperation(err)

	// try to delete configuration profile from storage
	profiles, err := s.Storage.DeleteConfigurationProfile(int(id))

	// check if the storage operation was successful
	if _, ok := err.(*storage.ItemNotFoundError); ok {
		responses.Send(http.StatusNotFound, writer, err.Error())
	} else if err != nil {
		responses.SendInternalServerError(writer, err.Error())
	} else {
		responses.SendOK(writer, responses.BuildOkResponseWithData("profiles", profiles))
	}
}

// ChangeConfigurationProfile method changes configuration profile
func (s Server) ChangeConfigurationProfile(writer http.ResponseWriter, request *http.Request) {
	// profile ID needs to be specified in request
	id, err := retrieveIDRequestParameter(request)

	// username needs to be specified in request
	username, foundUsername := request.URL.Query()["username"]

	// description needs to be specified in request
	description, foundDescription := request.URL.Query()["description"]

	if err != nil {
		responses.SendBadRequest(writer, "Error reading profile ID from request\n")
		return
	}

	if !foundUsername {
		responses.SendBadRequest(writer, "User name needs to be specified\n")
		return
	}

	if !foundDescription {
		responses.SendBadRequest(writer, "Description needs to be specified\n")
		return
	}

	configuration, err := ioutil.ReadAll(request.Body)
	if err != nil || len(configuration) == 0 {
		responses.SendBadRequest(writer, "Configuration needs to be provided in the request body")
		return
	}

	// try to record the action ChangeConfigurationProfile into Splunk
	err = s.Splunk.LogAction("ChangeConfigurationProfile", username[0], string(configuration))
	// and check whether the Splunk operation was successful
	checkSplunkOperation(err)

	// try to change configuration profile configuration in storage
	profiles, err := s.Storage.ChangeConfigurationProfile(int(id), username[0], description[0], string(configuration))

	// check if the storage operation was successful
	if _, ok := err.(*storage.ItemNotFoundError); ok {
		responses.Send(http.StatusNotFound, writer, err.Error())
	} else if err != nil {
		responses.SendInternalServerError(writer, err.Error())
	} else {
		responses.SendOK(writer, responses.BuildOkResponseWithData("profiles", profiles))
	}
}
