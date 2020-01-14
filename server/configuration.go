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
	"github.com/RedHatInsighs/insights-operator-controller/utils"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

func sendConfiguration(writer http.ResponseWriter, configuration string) {
	resp := utils.BuildOkResponse()
	resp["configuration"] = configuration
	utils.SendResponse(writer, resp)
}

// GetConfiguration - return single configuration by id
func (s Server) GetConfiguration(writer http.ResponseWriter, request *http.Request) {
	id, found := mux.Vars(request)["id"]
	if !found {
		utils.SendError(writer, "Configuration ID needs to be specified")
		return
	}

	configuration, err := s.Storage.GetClusterConfigurationByID(id)
	if err != nil {
		utils.SendError(writer, err.Error())
		return
	}
	sendConfiguration(writer, configuration)
}

// DeleteConfiguration - remove single configuration by id
func (s Server) DeleteConfiguration(writer http.ResponseWriter, request *http.Request) {
	id, found := mux.Vars(request)["id"]
	if !found {
		utils.SendError(writer, "Configuration ID needs to be specified")
		return
	}

	s.Splunk.LogAction("DeleteClusterConfigurationById", "tester", id)
	err := s.Storage.DeleteClusterConfigurationByID(id)
	if err != nil {
		utils.SendError(writer, err.Error())
		return
	}
	utils.SendResponse(writer, utils.BuildOkResponse())
}

// GetAllConfigurations - return list of all configurations
func (s Server) GetAllConfigurations(writer http.ResponseWriter, request *http.Request) {
	configuration, err := s.Storage.ListAllClusterConfigurations()
	if err != nil {
		utils.SendError(writer, err.Error())
		return
	}
	utils.SendResponse(writer, utils.BuildOkResponseWithData("configuration", configuration))
}

// GetClusterConfiguration - return list of configuration for single cluster
func (s Server) GetClusterConfiguration(writer http.ResponseWriter, request *http.Request) {
	cluster, found := mux.Vars(request)["cluster"]
	if !found {
		utils.SendError(writer, "Cluster ID needs to be specified")
		return
	}

	configuration, err := s.Storage.ListClusterConfiguration(cluster)
	if err != nil {
		utils.SendError(writer, err.Error())
		return
	}
	utils.SendResponse(writer, utils.BuildOkResponseWithData("configuration", configuration))
}

// EnableOrDisableConfiguration - enable or disable single configuration
func (s Server) EnableOrDisableConfiguration(writer http.ResponseWriter, request *http.Request, active string) {
	id, found := mux.Vars(request)["id"]
	if !found {
		utils.SendError(writer, "Configuration ID needs to be specified")
		return
	}

	if active == "0" {
		s.Splunk.LogAction("DisableClusterConfiguration", "tester", id)
	} else {
		s.Splunk.LogAction("EnableClusterConfiguration", "tester", id)
	}
	err := s.Storage.EnableOrDisableClusterConfigurationByID(id, active)
	if err != nil {
		utils.SendError(writer, err.Error())
		return
	}
	if active == "0" {
		sendConfiguration(writer, "disabled")
	} else {
		sendConfiguration(writer, "enabled")
	}
}

// EnableConfiguration - enable sinlge configuration
func (s Server) EnableConfiguration(writer http.ResponseWriter, request *http.Request) {
	s.EnableOrDisableConfiguration(writer, request, "1")
}

// DisableConfiguration - disable single configuration
func (s Server) DisableConfiguration(writer http.ResponseWriter, request *http.Request) {
	s.EnableOrDisableConfiguration(writer, request, "0")
}

// NewClusterConfiguration - create configuration for single cluster
func (s Server) NewClusterConfiguration(writer http.ResponseWriter, request *http.Request) {
	cluster, found := mux.Vars(request)["cluster"]
	if !found {
		utils.SendError(writer, "Cluster ID needs to be specified")
		return
	}

	username, foundUsername := request.URL.Query()["username"]
	reason, foundReason := request.URL.Query()["reason"]
	description, foundDescription := request.URL.Query()["description"]

	if !foundUsername {
		utils.SendError(writer, "User name needs to be specified\n")
		return
	}

	if !foundReason {
		utils.SendError(writer, "Reason needs to be specified\n")
		return
	}

	if !foundDescription {
		utils.SendError(writer, "Description needs to be specified\n")
		return
	}

	configuration, err := ioutil.ReadAll(request.Body)
	if err != nil || len(configuration) == 0 {
		utils.SendError(writer, "Configuration needs to be provided in the request body")
		return
	}

	configurations, err := s.Storage.CreateClusterConfiguration(cluster, username[0], reason[0], description[0], string(configuration))
	if err != nil {
		utils.SendError(writer, err.Error())
		return
	}
	s.Splunk.LogAction("NewClusterConfiguration", "tester", string(configuration))
	utils.SendResponse(writer, utils.BuildOkResponseWithData("configurations", configurations))
}

// EnableClusterConfiguration - enable cluster configuration
func (s Server) EnableClusterConfiguration(writer http.ResponseWriter, request *http.Request) {
	cluster, found := mux.Vars(request)["cluster"]
	if !found {
		utils.SendError(writer, "Cluster ID needs to be specified")
		return
	}

	username, foundUsername := request.URL.Query()["username"]
	reason, foundReason := request.URL.Query()["reason"]

	if !foundUsername {
		utils.SendError(writer, "User name needs to be specified\n")
		return
	}

	if !foundReason {
		utils.SendError(writer, "Reason needs to be specified\n")
		return
	}

	s.Splunk.LogAction("EnableClusterConfiguration", username[0], cluster)
	configurations, err := s.Storage.EnableClusterConfiguration(cluster, username[0], reason[0])
	if err != nil {
		utils.SendError(writer, err.Error())
		return
	}
	utils.SendResponse(writer, utils.BuildOkResponseWithData("configurations", configurations))
}

// DisableClusterConfiguration - disable cluster configuration
func (s Server) DisableClusterConfiguration(writer http.ResponseWriter, request *http.Request) {
	cluster, found := mux.Vars(request)["cluster"]
	if !found {
		utils.SendError(writer, "Cluster ID needs to be specified")
		return
	}

	username, foundUsername := request.URL.Query()["username"]
	reason, foundReason := request.URL.Query()["reason"]

	if !foundUsername {
		utils.SendError(writer, "User name needs to be specified\n")
		return
	}

	if !foundReason {
		utils.SendError(writer, "Reason needs to be specified\n")
		return
	}

	s.Splunk.LogAction("DisableClusterConfiguration", username[0], cluster)
	configurations, err := s.Storage.DisableClusterConfiguration(cluster, username[0], reason[0])
	if err != nil {
		utils.SendError(writer, err.Error())
		return
	}
	utils.SendResponse(writer, utils.BuildOkResponseWithData("configurations", configurations))
}
