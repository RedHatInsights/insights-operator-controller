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
	"net/http"

	"github.com/RedHatInsights/insights-operator-controller/storage"
	"github.com/RedHatInsights/insights-operator-utils/responses"
	"github.com/gorilla/mux"
)

// GetAllTriggers - return list of all triggers
func (s Server) GetAllTriggers(writer http.ResponseWriter, request *http.Request) {
	triggers, err := s.Storage.ListAllTriggers()
	if err != nil {
		responses.SendInternalServerError(writer, err.Error())
		return
	}
	responses.SendResponse(writer, responses.BuildOkResponseWithData("triggers", triggers))
}

// GetTrigger - return single trigger by id
func (s Server) GetTrigger(writer http.ResponseWriter, request *http.Request) {
	id, err := retrieveIDRequestParameter(request)
	if err != nil {
		responses.Send(http.StatusBadRequest, writer, err.Error())
		return
	}

	trigger, err := s.Storage.GetTriggerByID(id)
	if err == storage.ErrNoSuchObj {
		responses.Send(http.StatusNotFound, writer, fmt.Sprintf("No such trigger for ID %v", id))
	} else if err != nil {
		responses.SendInternalServerError(writer, err.Error())
	} else {
		responses.SendResponse(writer, responses.BuildOkResponseWithData("trigger", trigger))
	}
}

// DeleteTrigger - delete single trigger
func (s Server) DeleteTrigger(writer http.ResponseWriter, request *http.Request) {
	id, err := retrieveIDRequestParameter(request)
	if err != nil {
		responses.Send(http.StatusBadRequest, writer, err.Error())
		return
	}

	// try to record the action DeleteTrigger into Splunk
	err = s.Splunk.LogAction("DeleteTrigger", "tester", fmt.Sprint(id))
	// and check whether the Splunk operation was successful
	checkSplunkOperation(err)

	err = s.Storage.DeleteTriggerByID(id)
	if _, ok := err.(*storage.ItemNotFoundError); ok {
		responses.Send(
			http.StatusNotFound,
			writer,
			responses.BuildOkResponse(),
		)
	} else if err != nil {
		responses.Send(http.StatusInternalServerError, writer, err.Error())
	} else {
		responses.SendResponse(writer, responses.BuildOkResponse())
	}
}

// ActivateTrigger - active single trigger
func (s Server) ActivateTrigger(writer http.ResponseWriter, request *http.Request) {
	id, err := retrieveIDRequestParameter(request)
	if err != nil {
		responses.Send(http.StatusBadRequest, writer, err.Error())
		return
	}

	// try to record the action ActivateTrigger into Splunk
	err = s.Splunk.LogAction("ActivateTrigger", "tester", fmt.Sprint(id))
	// and check whether the Splunk operation was successful
	checkSplunkOperation(err)

	err = s.Storage.ChangeStateOfTriggerByID(id, 1)
	if _, ok := err.(*storage.ItemNotFoundError); ok {
		responses.Send(http.StatusNotFound, writer, err.Error())
	} else if err != nil {
		responses.SendInternalServerError(writer, err.Error())
	} else {
		responses.SendResponse(writer, responses.BuildOkResponse())
	}
}

// DeactivateTrigger - deactivate single trigger
func (s Server) DeactivateTrigger(writer http.ResponseWriter, request *http.Request) {
	id, err := retrieveIDRequestParameter(request)
	if err != nil {
		responses.Send(http.StatusBadRequest, writer, err.Error())
		return
	}

	s.Splunk.LogAction("DeactivateTrigger", "tester", fmt.Sprint(id))
	err = s.Storage.ChangeStateOfTriggerByID(id, 0)
	if _, ok := err.(*storage.ItemNotFoundError); ok {
		responses.Send(http.StatusNotFound, writer, err.Error())
	} else if err != nil {
		responses.SendInternalServerError(writer, err.Error())
	} else {
		responses.SendResponse(writer, responses.BuildOkResponse())
	}
}

// GetClusterTriggers - return list of triggers for single cluster
func (s Server) GetClusterTriggers(writer http.ResponseWriter, request *http.Request) {
	cluster, found := mux.Vars(request)["cluster"]
	if !found {
		responses.SendError(writer, "Cluster name needs to be specified")
		return
	}

	triggers, err := s.Storage.ListClusterTriggers(cluster)
	if _, ok := err.(*storage.ItemNotFoundError); ok {
		responses.Send(http.StatusNotFound, writer, err.Error())
	} else if err != nil {
		responses.SendInternalServerError(writer, err.Error())
	} else {
		responses.SendResponse(writer, responses.BuildOkResponseWithData("triggers", triggers))
	}
}

// RegisterClusterTrigger - register new trigger for cluster
func (s Server) RegisterClusterTrigger(writer http.ResponseWriter, request *http.Request) {
	cluster, found := mux.Vars(request)["cluster"]
	if !found {
		responses.SendError(writer, "Cluster name needs to be specified")
		return
	}

	triggerType, found := mux.Vars(request)["trigger"]
	if !found {
		responses.SendError(writer, "Trigger type needs to be specified")
		return
	}

	username, foundUsername := request.URL.Query()["username"]
	if !foundUsername {
		responses.SendError(writer, "User name needs to be specified\n")
		return
	}

	reason, foundReason := request.URL.Query()["reason"]
	if !foundReason {
		responses.SendError(writer, "Reason needs to be specified\n")
		return
	}

	link, foundReason := request.URL.Query()["link"]
	if !foundReason {
		responses.SendError(writer, "Link needs to be specified\n")
		return
	}

	s.Splunk.LogTriggerAction("RegisterTrigger", username[0], cluster, triggerType)
	err := s.Storage.NewTrigger(cluster, triggerType, username[0], reason[0], link[0])
	if _, ok := err.(*storage.ItemNotFoundError); ok {
		responses.Send(http.StatusNotFound, writer, err.Error())
	} else if err != nil {
		responses.SendInternalServerError(writer, err.Error())
	} else {
		responses.SendResponse(writer, responses.BuildOkResponse())
	}
}
