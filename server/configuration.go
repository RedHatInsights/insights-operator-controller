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

func sendConfiguration(writer http.ResponseWriter, configuration string) {
	resp := responses.BuildOkResponse()
	resp["configuration"] = configuration
	responses.SendResponse(writer, resp)
}

// GetConfiguration - return single configuration by id
func (s Server) GetConfiguration(writer http.ResponseWriter, request *http.Request) {
	id, err := retrieveIDRequestParameter(request)
	if err != nil {
		responses.Send(http.StatusBadRequest, writer, err.Error())
		return
	}

	configuration, err := s.Storage.GetClusterConfigurationByID(id)
	if _, ok := err.(*storage.ItemNotFoundError); ok {
		responses.Send(http.StatusNotFound, writer, err.Error())
	} else if err != nil {
		responses.SendInternalServerError(writer, err.Error())
	} else {
		sendConfiguration(writer, configuration)
	}
}

// DeleteConfiguration - remove single configuration by id
func (s Server) DeleteConfiguration(writer http.ResponseWriter, request *http.Request) {
	id, err := retrieveIDRequestParameter(request)
	if err != nil {
		responses.Send(http.StatusBadRequest, writer, err.Error())
		return
	}

	// try to record the action DeleteConfigurationById into Splunk
	err = s.Splunk.LogAction("DeleteClusterConfigurationById", "tester", fmt.Sprint(id))
	if err != nil {
		log.Println("Unable to write log into Splunk", err)
	}

	err = s.Storage.DeleteClusterConfigurationByID(id)
	if _, ok := err.(*storage.ItemNotFoundError); ok {
		responses.Send(http.StatusNotFound, writer, err.Error())
	} else if err != nil {
		responses.SendInternalServerError(writer, err.Error())
	} else {
		responses.SendResponse(writer, responses.BuildOkResponse())
	}
}

// GetAllConfigurations - return list of all configurations
func (s Server) GetAllConfigurations(writer http.ResponseWriter, request *http.Request) {
	configuration, err := s.Storage.ListAllClusterConfigurations()
	if err != nil {
		responses.SendInternalServerError(writer, err.Error())
		return
	}
	responses.SendResponse(writer, responses.BuildOkResponseWithData("configuration", configuration))
}

// GetClusterConfiguration - return list of configuration for single cluster
func (s Server) GetClusterConfiguration(writer http.ResponseWriter, request *http.Request) {
	cluster, found := mux.Vars(request)["cluster"]
	if !found {
		responses.Send(http.StatusBadRequest, writer, "Cluster ID needs to be specified")
		return
	}

	configuration, err := s.Storage.ListClusterConfiguration(cluster)
	if _, ok := err.(*storage.ItemNotFoundError); ok {
		responses.Send(http.StatusNotFound, writer, err.Error())
	} else if err != nil {
		responses.SendInternalServerError(writer, err.Error())
	} else {
		responses.SendResponse(writer, responses.BuildOkResponseWithData("configuration", configuration))
	}
}

// EnableOrDisableConfiguration - enable or disable single configuration
func (s Server) EnableOrDisableConfiguration(writer http.ResponseWriter, request *http.Request, active string) {
	id, err := retrieveIDRequestParameter(request)
	if err != nil {
		responses.Send(http.StatusBadRequest, writer, err.Error())
		return
	}

	if active == "0" {
		err = s.Splunk.LogAction("DisableClusterConfiguration", "tester", fmt.Sprint(id))
		if err != nil {
			log.Println("Unable to write log into Splunk", err)
		}
	} else {
		err = s.Splunk.LogAction("EnableClusterConfiguration", "tester", fmt.Sprint(id))
		if err != nil {
			log.Println("Unable to write log into Splunk", err)
		}
	}
	err = s.Storage.EnableOrDisableClusterConfigurationByID(id, active)
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
		responses.Send(http.StatusBadRequest, writer, "Cluster ID needs to be specified")
		return
	}

	username, foundUsername := request.URL.Query()["username"]
	reason, foundReason := request.URL.Query()["reason"]
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

	configuration, err := ioutil.ReadAll(request.Body)
	if err != nil || len(configuration) == 0 {
		responses.SendError(writer, "Configuration needs to be provided in the request body")
		return
	}

	configurations, err := s.Storage.CreateClusterConfiguration(cluster, username[0], reason[0], description[0], string(configuration))
	if err != nil {
		responses.SendInternalServerError(writer, err.Error())
		return
	}
	s.Splunk.LogAction("NewClusterConfiguration", "tester", string(configuration))
	responses.SendResponse(writer, responses.BuildOkResponseWithData("configurations", configurations))
}

// EnableClusterConfiguration - enable cluster configuration
func (s Server) EnableClusterConfiguration(writer http.ResponseWriter, request *http.Request) {
	cluster, found := mux.Vars(request)["cluster"]
	if !found {
		responses.Send(http.StatusBadRequest, writer, "Cluster ID needs to be specified")
		return
	}

	username, foundUsername := request.URL.Query()["username"]
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
	if err != nil {
		log.Println("Unable to write log into Splunk", err)
	}

	configurations, err := s.Storage.EnableClusterConfiguration(cluster, username[0], reason[0])
	if err != nil {
		responses.SendInternalServerError(writer, err.Error())
		return
	}
	responses.SendResponse(writer, responses.BuildOkResponseWithData("configurations", configurations))
}

// DisableClusterConfiguration - disable cluster configuration
func (s Server) DisableClusterConfiguration(writer http.ResponseWriter, request *http.Request) {
	cluster, found := mux.Vars(request)["cluster"]
	if !found {
		responses.Send(http.StatusBadRequest, writer, "Cluster ID needs to be specified")
		return
	}

	username, foundUsername := request.URL.Query()["username"]
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
	if err != nil {
		log.Println("Unable to write log into Splunk", err)
	}

	configurations, err := s.Storage.DisableClusterConfiguration(cluster, username[0], reason[0])
	if err != nil {
		responses.SendInternalServerError(writer, err.Error())
		return
	}
	responses.SendResponse(writer, responses.BuildOkResponseWithData("configurations", configurations))
}
