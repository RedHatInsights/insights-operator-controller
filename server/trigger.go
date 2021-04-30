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
// https://redhatinsights.github.io/insights-operator-controller/packages/server/trigger.html

import (
	"fmt"
	"net/http"

	"github.com/RedHatInsights/insights-operator-controller/storage"
	"github.com/RedHatInsights/insights-operator-utils/responses"
	"github.com/gorilla/mux"
)

// GetAllTriggers method returns list of all triggers
func (s Server) GetAllTriggers(writer http.ResponseWriter, request *http.Request) {
	// try to read list of all triggers from storage
	triggers, err := s.Storage.ListAllTriggers()

	// check if the storage operation has been successful
	if err != nil {
		TryToSendInternalServerError(writer, err.Error())
		return
	}
	responses.SendOK(writer, responses.BuildOkResponseWithData("triggers", triggers))
}

// GetTrigger method returns single trigger by id
func (s Server) GetTrigger(writer http.ResponseWriter, request *http.Request) {
	// trigger ID needs to be specified in request parameter
	id, err := retrieveIDRequestParameter(request)
	if err != nil {
		responses.Send(http.StatusBadRequest, writer, err.Error())
		return
	}

	// try to read trigger identified by its ID from storage
	trigger, err := s.Storage.GetTriggerByID(id)

	// check if the storage operation has been successful
	if err == storage.ErrNoSuchObj {
		responses.Send(http.StatusNotFound, writer, fmt.Sprintf("No such trigger for ID %v", id))
	} else if err != nil {
		TryToSendInternalServerError(writer, err.Error())
	} else {
		responses.SendOK(writer, responses.BuildOkResponseWithData("trigger", trigger))
	}
}

// DeleteTrigger method deletes single trigger
func (s Server) DeleteTrigger(writer http.ResponseWriter, request *http.Request) {
	// trigger ID needs to be specified in request parameter
	id, err := retrieveIDRequestParameter(request)
	if err != nil {
		responses.Send(http.StatusBadRequest, writer, err.Error())
		return
	}

	// try to record the action DeleteTrigger into Splunk
	err = s.Splunk.LogAction("DeleteTrigger", "tester", fmt.Sprint(id))
	// and check whether the Splunk operation was successful
	checkSplunkOperation(err)

	// try to delete trigger identified by its ID from storage
	err = s.Storage.DeleteTriggerByID(id)

	// check if the storage operation has been successful
	if _, ok := err.(*storage.ItemNotFoundError); ok {
		responses.Send(
			http.StatusNotFound,
			writer,
			responses.BuildOkResponse(),
		)
	} else if err != nil {
		responses.Send(http.StatusInternalServerError, writer, err.Error())
	} else {
		responses.SendOK(writer, responses.BuildOkResponse())
	}
}

// ActivateTrigger method actives single trigger
func (s Server) ActivateTrigger(writer http.ResponseWriter, request *http.Request) {
	// trigger ID needs to be specified in request parameter
	id, err := retrieveIDRequestParameter(request)
	if err != nil {
		responses.Send(http.StatusBadRequest, writer, err.Error())
		return
	}

	// try to record the action ActivateTrigger into Splunk
	err = s.Splunk.LogAction("ActivateTrigger", "tester", fmt.Sprint(id))
	// and check whether the Splunk operation was successful
	checkSplunkOperation(err)

	// try to activate trigger identified by its ID from storage
	err = s.Storage.ChangeStateOfTriggerByID(id, 1)

	// check if the storage operation has been successful
	if _, ok := err.(*storage.ItemNotFoundError); ok {
		responses.Send(http.StatusNotFound, writer, err.Error())
	} else if err != nil {
		TryToSendInternalServerError(writer, err.Error())
	} else {
		responses.SendOK(writer, responses.BuildOkResponse())
	}
}

// DeactivateTrigger method deactivates single trigger
func (s Server) DeactivateTrigger(writer http.ResponseWriter, request *http.Request) {
	// trigger ID needs to be specified in request parameter
	id, err := retrieveIDRequestParameter(request)
	if err != nil {
		responses.Send(http.StatusBadRequest, writer, err.Error())
		return
	}

	// try to record the action DeactivateTrigger into Splunk
	err = s.Splunk.LogAction("DeactivateTrigger", "tester", fmt.Sprint(id))
	// and check whether the Splunk operation was successful
	checkSplunkOperation(err)

	// try to deactivate trigger identified by its ID from storage
	err = s.Storage.ChangeStateOfTriggerByID(id, 0)

	// check if the storage operation has been successful
	if _, ok := err.(*storage.ItemNotFoundError); ok {
		responses.Send(http.StatusNotFound, writer, err.Error())
	} else if err != nil {
		TryToSendInternalServerError(writer, err.Error())
	} else {
		responses.SendOK(writer, responses.BuildOkResponse())
	}
}

// GetClusterTriggers method returns list of triggers for single cluster
func (s Server) GetClusterTriggers(writer http.ResponseWriter, request *http.Request) {
	// cluster name needs to be specified in request parameter
	cluster, found := mux.Vars(request)["cluster"]
	if !found {
		responses.SendBadRequest(writer, "Cluster name needs to be specified")
		return
	}

	// try to read list of all triggers for specified cluster from storage
	triggers, err := s.Storage.ListClusterTriggers(cluster)

	// check if the storage operation has been successful
	if _, ok := err.(*storage.ItemNotFoundError); ok {
		responses.Send(http.StatusNotFound, writer, err.Error())
	} else if err != nil {
		TryToSendInternalServerError(writer, err.Error())
	} else {
		responses.SendOK(writer, responses.BuildOkResponseWithData("triggers", triggers))
	}
}

// RegisterClusterTrigger method registers new trigger for cluster
func (s Server) RegisterClusterTrigger(writer http.ResponseWriter, request *http.Request) {
	// cluster name needs to be specified in request parameter
	cluster, found := mux.Vars(request)["cluster"]
	if !found {
		responses.SendBadRequest(writer, "Cluster name needs to be specified")
		return
	}

	// trigger type needs to be specified in request parameter
	triggerType, found := mux.Vars(request)["trigger"]
	if !found {
		responses.SendBadRequest(writer, "Trigger type needs to be specified")
		return
	}

	// user name needs to be specified in request parameter
	username, foundUsername := request.URL.Query()["username"]
	if !foundUsername {
		responses.SendBadRequest(writer, "User name needs to be specified\n")
		return
	}

	// reason needs to be specified in request parameter
	reason, foundReason := request.URL.Query()["reason"]
	if !foundReason {
		responses.SendBadRequest(writer, "Reason needs to be specified\n")
		return
	}

	// link needs to be specified in request parameter
	link, foundReason := request.URL.Query()["link"]
	if !foundReason {
		responses.SendBadRequest(writer, "Link needs to be specified\n")
		return
	}

	// try to record the action RegisterTrigger into Splunk
	err := s.Splunk.LogTriggerAction("RegisterTrigger", username[0], cluster, triggerType)
	// and check whether the Splunk operation was successful
	checkSplunkOperation(err)

	// try to create new trigger in storage
	err = s.Storage.NewTrigger(cluster, triggerType, username[0], reason[0], link[0])

	// check if the storage operation has been successful
	if _, ok := err.(*storage.ItemNotFoundError); ok {
		responses.Send(http.StatusNotFound, writer, err.Error())
	} else if err != nil {
		TryToSendInternalServerError(writer, err.Error())
	} else {
		responses.SendOK(writer, responses.BuildOkResponse())
	}
}
