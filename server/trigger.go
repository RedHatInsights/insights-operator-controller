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
	"github.com/RedHatInsights/insights-operator-controller/utils"
	"github.com/gorilla/mux"
	"net/http"
)

// GetAllTriggers - return list of all triggers
func (s Server) GetAllTriggers(writer http.ResponseWriter, request *http.Request) {
	triggers, err := s.Storage.ListAllTriggers()
	if err != nil {
		utils.SendError(writer, err.Error())
		return
	}
	utils.SendResponse(writer, utils.BuildOkResponseWithData("triggers", triggers))
}

// GetTrigger - return single trigger by id
func (s Server) GetTrigger(writer http.ResponseWriter, request *http.Request) {
	id, found := mux.Vars(request)["id"]
	if !found {
		utils.SendError(writer, "Trigger ID needs to be specified")
		return
	}

	trigger, err := s.Storage.GetTriggerByID(id)
	if err != nil {
		utils.SendError(writer, err.Error())
		return
	}
	utils.SendResponse(writer, utils.BuildOkResponseWithData("trigger", trigger))
}

// DeleteTrigger - delete single trigger
func (s Server) DeleteTrigger(writer http.ResponseWriter, request *http.Request) {
	id, found := mux.Vars(request)["id"]
	if !found {
		utils.SendError(writer, "Trigger ID needs to be specified")
		return
	}

	s.Splunk.LogAction("DeleteTrigger", "tester", id)
	err := s.Storage.DeleteTriggerByID(id)
	if err != nil {
		utils.SendError(writer, err.Error())
		return
	}
	utils.SendResponse(writer, utils.BuildOkResponse())
}

// ActivateTrigger - active single trigger
func (s Server) ActivateTrigger(writer http.ResponseWriter, request *http.Request) {
	id, found := mux.Vars(request)["id"]
	if !found {
		utils.SendError(writer, "Trigger ID needs to be specified")
		return
	}

	s.Splunk.LogAction("ActivateTrigger", "tester", id)
	err := s.Storage.ChangeStateOfTriggerByID(id, 1)
	if err != nil {
		utils.SendError(writer, err.Error())
		return
	}
	utils.SendResponse(writer, utils.BuildOkResponse())
}

// DeactivateTrigger - deactivate single trigger
func (s Server) DeactivateTrigger(writer http.ResponseWriter, request *http.Request) {
	id, found := mux.Vars(request)["id"]
	if !found {
		utils.SendError(writer, "Trigger ID needs to be specified")
		return
	}

	s.Splunk.LogAction("DeactivateTrigger", "tester", id)
	err := s.Storage.ChangeStateOfTriggerByID(id, 0)
	if err != nil {
		utils.SendError(writer, err.Error())
		return
	}
	utils.SendResponse(writer, utils.BuildOkResponse())
}

// GetClusterTriggers - return list of triggers for single cluster
func (s Server) GetClusterTriggers(writer http.ResponseWriter, request *http.Request) {
	cluster, found := mux.Vars(request)["cluster"]
	if !found {
		utils.SendError(writer, "Cluster name needs to be specified")
		return
	}

	triggers, err := s.Storage.ListClusterTriggers(cluster)
	if err != nil {
		utils.SendError(writer, err.Error())
		return
	}
	utils.SendResponse(writer, utils.BuildOkResponseWithData("triggers", triggers))
}

// RegisterClusterTrigger - register new trigger for cluster
func (s Server) RegisterClusterTrigger(writer http.ResponseWriter, request *http.Request) {
	cluster, found := mux.Vars(request)["cluster"]
	if !found {
		utils.SendError(writer, "Cluster name needs to be specified")
		return
	}

	triggerType, found := mux.Vars(request)["trigger"]
	if !found {
		utils.SendError(writer, "Trigger type needs to be specified")
		return
	}

	username, foundUsername := request.URL.Query()["username"]
	if !foundUsername {
		utils.SendError(writer, "User name needs to be specified\n")
		return
	}

	reason, foundReason := request.URL.Query()["reason"]
	if !foundReason {
		utils.SendError(writer, "Reason needs to be specified\n")
		return
	}

	link, foundReason := request.URL.Query()["link"]
	if !foundReason {
		utils.SendError(writer, "Link needs to be specified\n")
		return
	}

	s.Splunk.LogTriggerAction("RegisterTrigger", username[0], cluster, triggerType)
	err := s.Storage.NewTrigger(cluster, triggerType, username[0], reason[0], link[0])
	if err != nil {
		utils.SendError(writer, err.Error())
		return
	}
	utils.SendResponse(writer, utils.BuildOkResponse())
}
