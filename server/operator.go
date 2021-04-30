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
// https://redhatinsights.github.io/insights-operator-controller/packages/server/operator.html

import (
	"fmt"
	"log"
	"net/http"

	"github.com/RedHatInsights/insights-operator-controller/storage"
	"github.com/RedHatInsights/insights-operator-utils/responses"
	"github.com/gorilla/mux"
)

// ReadConfigurationForOperator method reads configuration for the operator.
func (s Server) ReadConfigurationForOperator(writer http.ResponseWriter, request *http.Request) {
	// cluster name needs to be specified in request
	cluster, found := mux.Vars(request)["cluster"]
	if !found {
		log.Println("Cluster name is not provided")
		TryToSendBadRequestServerResponse(writer, "Cluster ID needs to be specified")
		return
	}

	// try to read cluster active configuration from storage
	configuration, err := s.Storage.GetClusterActiveConfiguration(cluster)

	// check if the storage operation has been successful
	if itemNotFoundError, ok := err.(*storage.ItemNotFoundError); ok {
		responses.Send(
			http.StatusNotFound,
			writer,
			fmt.Sprintf("unable to read any active configuration for the cluster %v",
				itemNotFoundError.ItemID),
		)
	} else if err != nil {
		log.Println("Cannot read cluster configuration", err)
		TryToSendInternalServerError(writer, err.Error())
	} else {
		sendConfiguration(writer, configuration)
	}
}

// RegisterCluster method registers new cluster.
func (s Server) RegisterCluster(writer http.ResponseWriter, request *http.Request) {
	// cluster name needs to be specified in request
	clusterName, foundName := mux.Vars(request)["cluster"]

	// check parameters provided by client
	if !foundName {
		log.Println("Cluster name is not provided")
		TryToSendBadRequestServerResponse(writer, "Cluster name needs to be specified")
		return
	}

	// try to record the action RegisterCluster into Splunk
	err := s.Splunk.LogAction("RegisterCluster", "tester", clusterName)
	if err != nil {
		log.Println("(not critical) Log into splunk failed", err)
	}

	// register new cluster in the storage
	err = s.Storage.RegisterNewCluster(clusterName)

	// check if the storage operation has been successful
	if err != nil {
		log.Println("Cannot create new cluster", err)
		TryToSendInternalServerError(writer, err.Error())
	}
	responses.SendCreated(writer, responses.BuildOkResponse())
}

// GetActiveTriggersForCluster method returns list of triggers for single cluster
func (s Server) GetActiveTriggersForCluster(writer http.ResponseWriter, request *http.Request) {
	// cluster name needs to be specified in request
	cluster, found := mux.Vars(request)["cluster"]
	if !found {
		TryToSendBadRequestServerResponse(writer, "Cluster name needs to be specified")
		return
	}

	// try to read list of active cluster triggers
	triggers, err := s.Storage.ListActiveClusterTriggers(cluster)

	// check if the storage operation has been successful
	if _, ok := err.(*storage.ItemNotFoundError); ok {
		responses.Send(http.StatusNotFound, writer, err.Error())
	} else if err != nil {
		TryToSendInternalServerError(writer, err.Error())
	} else {
		TryToSendOKServerResponse(writer, responses.BuildOkResponseWithData("triggers", triggers))
	}
}

// AckTriggerForCluster method perform ack for single cluster's trigger
func (s Server) AckTriggerForCluster(writer http.ResponseWriter, request *http.Request) {
	// cluster name needs to be specified in request
	cluster, found := mux.Vars(request)["cluster"]
	if !found {
		TryToSendBadRequestServerResponse(writer, "Cluster name needs to be specified")
		return
	}

	// trigger ID needs to be specified in request
	triggerID, err := retrievePositiveIntRequestParameter(request, "trigger")
	if err != nil {
		responses.Send(http.StatusBadRequest, writer, err.Error())
		return
	}

	// try to ack cluster in storage
	err = s.Storage.AckTrigger(cluster, triggerID)

	// check if the storage operation has been successful
	if _, ok := err.(*storage.ItemNotFoundError); ok {
		responses.Send(http.StatusNotFound, writer, err.Error())
	} else if err != nil {
		TryToSendInternalServerError(writer, err.Error())
	} else {
		TryToSendOKServerResponse(writer, responses.BuildOkResponse())
	}
}
