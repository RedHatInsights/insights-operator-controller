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
	"fmt"
	"log"

	"github.com/RedHatInsights/insights-operator-controller/storage"
	"github.com/RedHatInsights/insights-operator-utils/responses"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

// sendConfiguration is helper function to send cluster configuration to client
func sendConfiguration(writer http.ResponseWriter, configuration string) {
	resp := responses.BuildOkResponse()
	resp["configuration"] = configuration
	responses.SendResponse(writer, resp)
}

// GetConfiguration method returns single configuration by id
func (s Server) GetConfiguration(writer http.ResponseWriter, request *http.Request) {
	// configuration ID needs to be specified in request
	id, err := retrieveIDRequestParameter(request)
	if err != nil {
		responses.Send(http.StatusBadRequest, writer, err.Error())
		return
	}

	// try to read cluster configuration specified by ID from storage
	configuration, err := s.Storage.GetClusterConfigurationByID(id)

	// check if storage operation has been successful
	if _, ok := err.(*storage.ItemNotFoundError); ok {
		responses.Send(http.StatusNotFound, writer, err.Error())
	} else if err != nil {
		responses.SendInternalServerError(writer, err.Error())
	} else {
		sendConfiguration(writer, configuration)
	}
}

// DeleteConfiguration method removes single configuration specified by its ID
func (s Server) DeleteConfiguration(writer http.ResponseWriter, request *http.Request) {
	// configuration ID needs to be specified in request
	id, err := retrieveIDRequestParameter(request)
	if err != nil {
		responses.Send(http.StatusBadRequest, writer, err.Error())
		return
	}

	// try to record the action DeleteConfigurationById into Splunk
	err = s.Splunk.LogAction("DeleteClusterConfigurationById", "tester", fmt.Sprint(id))
	// and check whether the Splunk operation was successful
	checkSplunkOperation(err)

	// try to delete cluster configuration specified by its ID from storage
	err = s.Storage.DeleteClusterConfigurationByID(id)

	// check if storage operation has been successful
	if _, ok := err.(*storage.ItemNotFoundError); ok {
		responses.Send(http.StatusNotFound, writer, err.Error())
	} else if err != nil {
		responses.SendInternalServerError(writer, err.Error())
	} else {
		responses.SendResponse(writer, responses.BuildOkResponse())
	}
}

// GetAllConfigurations method reads and returns list of all configurations
func (s Server) GetAllConfigurations(writer http.ResponseWriter, request *http.Request) {
	// try to read list of all configurations from storage
	configuration, err := s.Storage.ListAllClusterConfigurations()

	// check if storage operation has been successful
	if err != nil {
		responses.SendInternalServerError(writer, err.Error())
		return
	}
	responses.SendResponse(writer, responses.BuildOkResponseWithData("configuration", configuration))
}

// GetClusterConfiguration method returns list of configuration for single cluster
func (s Server) GetClusterConfiguration(writer http.ResponseWriter, request *http.Request) {
	// cluster name needs to be specified it request
	cluster, found := mux.Vars(request)["cluster"]
	if !found {
		responses.Send(http.StatusBadRequest, writer, "Cluster ID needs to be specified")
		return
	}

	// try to read list of cluster configurations from storage
	configuration, err := s.Storage.ListClusterConfiguration(cluster)

	// check if storage operation has been successful
	if _, ok := err.(*storage.ItemNotFoundError); ok {
		responses.Send(http.StatusNotFound, writer, err.Error())
	} else if err != nil {
		responses.SendInternalServerError(writer, err.Error())
	} else {
		responses.SendResponse(writer, responses.BuildOkResponseWithData("configuration", configuration))
	}
}

// EnableOrDisableConfiguration method enables or disables single cluster configuration
func (s Server) EnableOrDisableConfiguration(writer http.ResponseWriter, request *http.Request, active string) {
	// configuration ID needs to be specified in request
	id, err := retrieveIDRequestParameter(request)
	if err != nil {
		responses.Send(http.StatusBadRequest, writer, err.Error())
		return
	}

	// "0" - disable
	// "1" (or other value) - enable
	if active == "0" {
		// try to write information about DisableClusterConfiguration operation into Splunk
		err = s.Splunk.LogAction("DisableClusterConfiguration", "tester", fmt.Sprint(id))
		// and check whether the Splunk operation was successful
		checkSplunkOperation(err)
	} else {
		// try to write information about EnableClusterConfiguration operation into Splunk
		err = s.Splunk.LogAction("EnableClusterConfiguration", "tester", fmt.Sprint(id))
		// and check whether the Splunk operation was successful
		checkSplunkOperation(err)
	}

	// try to enable or disable cluster configuration specified by its ID in storage
	err = s.Storage.EnableOrDisableClusterConfigurationByID(id, active)

	// check if storage operation has been successful
	if _, ok := err.(*storage.ItemNotFoundError); ok {
		responses.Send(http.StatusNotFound, writer, err.Error())
	} else if err != nil {
		responses.SendInternalServerError(writer, err.Error())
	} else {
		if active == "0" {
			sendConfiguration(writer, "disabled")
		} else {
			sendConfiguration(writer, "enabled")
		}
	}
}

// EnableConfiguration method enables single configuration
func (s Server) EnableConfiguration(writer http.ResponseWriter, request *http.Request) {
	s.EnableOrDisableConfiguration(writer, request, "1")
}

// DisableConfiguration method disables single configuration
func (s Server) DisableConfiguration(writer http.ResponseWriter, request *http.Request) {
	s.EnableOrDisableConfiguration(writer, request, "0")
}

// NewClusterConfiguration method creates configuration for single cluster
func (s Server) NewClusterConfiguration(writer http.ResponseWriter, request *http.Request) {
	// cluster name needs to be specified in request
	cluster, found := mux.Vars(request)["cluster"]
	if !found {
		responses.Send(http.StatusBadRequest, writer, "Cluster ID needs to be specified")
		return
	}

	// username needs to be specified in request
	username, foundUsername := request.URL.Query()["username"]

	// reason needs to be specified in request
	reason, foundReason := request.URL.Query()["reason"]

	// description needs to be specified in request
	description, foundDescription := request.URL.Query()["description"]

	if !foundUsername {
		responses.SendError(writer, "User name needs to be specified\n")
		return
	}

	if !foundReason {
		responses.SendError(writer, "Reason needs to be specified\n")
		return
	}

	if !foundDescription {
		responses.SendError(writer, "Description needs to be specified\n")
		return
	}

	// try to read configuration from request body
	configuration, err := ioutil.ReadAll(request.Body)
	if err != nil || len(configuration) == 0 {
		responses.SendError(writer, "Configuration needs to be provided in the request body")
		return
	}

	// try to create cluster configuration in storage
	configurations, err := s.Storage.CreateClusterConfiguration(cluster, username[0], reason[0], description[0], string(configuration))
	if err != nil {
		responses.SendInternalServerError(writer, err.Error())
		return
	}

	// try to write information about NewClusterConfiguration operation into Splunk
	err = s.Splunk.LogAction("NewClusterConfiguration", "tester", string(configuration))
	// and check whether the Splunk operation was successful
	checkSplunkOperation(err)

	responses.SendResponse(writer, responses.BuildOkResponseWithData("configurations", configurations))
}

// EnableClusterConfiguration method enables cluster configuration
func (s Server) EnableClusterConfiguration(writer http.ResponseWriter, request *http.Request) {
	// cluster name needs to be specified in request
	cluster, found := mux.Vars(request)["cluster"]
	if !found {
		responses.Send(http.StatusBadRequest, writer, "Cluster ID needs to be specified")
		return
	}

	// username needs to be specified in request
	username, foundUsername := request.URL.Query()["username"]

	// reason needs to be specified in request
	reason, foundReason := request.URL.Query()["reason"]

	if !foundUsername {
		responses.SendError(writer, "User name needs to be specified\n")
		return
	}

	if !foundReason {
		responses.SendError(writer, "Reason needs to be specified\n")
		return
	}

	// try to write information about EnableClusterConfiguration operation into Splunk
	err := s.Splunk.LogAction("EnableClusterConfiguration", username[0], cluster)
	// and check whether the Splunk operation was successful
	checkSplunkOperation(err)

	// perform storage operation, read the configuration
	configurations, err := s.Storage.EnableClusterConfiguration(cluster, username[0], reason[0])

	// check if storage operation has been successful
	if err != nil {
		responses.SendInternalServerError(writer, err.Error())
		return
	}
	responses.SendResponse(writer, responses.BuildOkResponseWithData("configurations", configurations))
}

// DisableClusterConfiguration method disables cluster configuration
func (s Server) DisableClusterConfiguration(writer http.ResponseWriter, request *http.Request) {
	// cluster name needs to be specified in request
	cluster, found := mux.Vars(request)["cluster"]
	if !found {
		responses.Send(http.StatusBadRequest, writer, "Cluster ID needs to be specified")
		return
	}

	// username needs to be specified in request
	username, foundUsername := request.URL.Query()["username"]

	// reason needs to be specified in request
	reason, foundReason := request.URL.Query()["reason"]

	if !foundUsername {
		responses.SendError(writer, "User name needs to be specified\n")
		return
	}

	if !foundReason {
		responses.SendError(writer, "Reason needs to be specified\n")
		return
	}

	// try to write information about DisableClusterConfiguration operation into Splunk
	err := s.Splunk.LogAction("DisableClusterConfiguration", username[0], cluster)
	// and check whether the Splunk operation was successful
	checkSplunkOperation(err)

	// perform storage operation, read the configuration
	configurations, err := s.Storage.DisableClusterConfiguration(cluster, username[0], reason[0])

	// check if storage operation has been successful
	if err != nil {
		responses.SendInternalServerError(writer, err.Error())
		return
	}
	responses.SendResponse(writer, responses.BuildOkResponseWithData("configurations", configurations))
}
