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

import (
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/RedHatInsights/insights-operator-controller/storage"
	"github.com/RedHatInsights/insights-operator-utils/responses"
)

// ListConfigurationProfiles - read list of configuration profiles.
func (s Server) ListConfigurationProfiles(writer http.ResponseWriter, request *http.Request) {
	profiles, err := s.Storage.ListConfigurationProfiles()
	if err == nil {
		responses.SendResponse(writer, responses.BuildOkResponseWithData("profiles", profiles))
	} else {
		responses.SendInternalServerError(writer, err.Error())
	}
}

// GetConfigurationProfile - read profile specified by its ID
func (s Server) GetConfigurationProfile(writer http.ResponseWriter, request *http.Request) {
	id, err := retrieveIDRequestParameter(request)
	if err != nil {
		responses.SendError(writer, "Error reading profile ID from request\n")
		return
	}

	profile, err := s.Storage.GetConfigurationProfile(int(id))
	if _, ok := err.(*storage.ItemNotFoundError); ok {
		responses.Send(http.StatusNotFound, writer, err.Error())
	} else if err != nil {
		responses.SendInternalServerError(writer, err.Error())
	} else {
		responses.SendResponse(writer, responses.BuildOkResponseWithData("profile", profile))
	}
}

// NewConfigurationProfile - create new configuration profile
func (s Server) NewConfigurationProfile(writer http.ResponseWriter, request *http.Request) {
	username, foundUsername := request.URL.Query()["username"]
	description, foundDescription := request.URL.Query()["description"]

	if !foundUsername {
		responses.SendError(writer, "User name needs to be specified\n")
		return
	}

	if !foundDescription {
		responses.SendError(writer, "Description needs to be specified\n")
		return
	}

	configuration, err := ioutil.ReadAll(request.Body)

	if err != nil || len(configuration) == 0 {
		responses.SendError(writer, "Configuration needs to be provided in the request body")
		return
	}

	s.Splunk.LogAction("NewConfigurationProfile", username[0], string(configuration))
	profiles, err := s.Storage.StoreConfigurationProfile(username[0], description[0], string(configuration))
	if err != nil {
		responses.SendInternalServerError(writer, err.Error())
	} else {
		responses.SendCreated(writer, responses.BuildOkResponseWithData("profiles", profiles))
	}
}

// DeleteConfigurationProfile - delete configuration profile
func (s Server) DeleteConfigurationProfile(writer http.ResponseWriter, request *http.Request) {
	id, err := retrieveIDRequestParameter(request)
	if err != nil {
		responses.SendError(writer, "Error reading profile ID from request\n")
		return
	}

	s.Splunk.LogAction("DeleteConfigurationProfile", "tester", strconv.Itoa(int(id)))
	profiles, err := s.Storage.DeleteConfigurationProfile(int(id))
	if _, ok := err.(*storage.ItemNotFoundError); ok {
		responses.Send(http.StatusNotFound, writer, err.Error())
	} else if err != nil {
		responses.SendInternalServerError(writer, err.Error())
	} else {
		responses.SendResponse(writer, responses.BuildOkResponseWithData("profiles", profiles))
	}
}

// ChangeConfigurationProfile - change configuration profile
func (s Server) ChangeConfigurationProfile(writer http.ResponseWriter, request *http.Request) {
	id, err := retrieveIDRequestParameter(request)
	username, foundUsername := request.URL.Query()["username"]
	description, foundDescription := request.URL.Query()["description"]

	if err != nil {
		responses.SendError(writer, "Error reading profile ID from request\n")
		return
	}

	if !foundUsername {
		responses.SendError(writer, "User name needs to be specified\n")
		return
	}

	if !foundDescription {
		responses.SendError(writer, "Description needs to be specified\n")
		return
	}

	configuration, err := ioutil.ReadAll(request.Body)
	if err != nil || len(configuration) == 0 {
		responses.SendError(writer, "Configuration needs to be provided in the request body")
		return
	}

	s.Splunk.LogAction("ChangeConfigurationProfile", username[0], string(configuration))
	profiles, err := s.Storage.ChangeConfigurationProfile(int(id), username[0], description[0], string(configuration))
	if _, ok := err.(*storage.ItemNotFoundError); ok {
		responses.Send(http.StatusNotFound, writer, err.Error())
	} else if err != nil {
		responses.SendInternalServerError(writer, err.Error())
	} else {
		responses.SendResponse(writer, responses.BuildOkResponseWithData("profiles", profiles))
	}
}
