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
	"log"
	"net/http"
)

// ReadConfigurationForOperator - read configuration for the operator.
func (s Server) ReadConfigurationForOperator(writer http.ResponseWriter, request *http.Request) {
	cluster, found := mux.Vars(request)["cluster"]
	if !found {
		log.Println("Cluster name is not provided")
		utils.SendError(writer, "Cluster ID needs to be specified")
		return
	}

	configuration, err := s.Storage.GetClusterActiveConfiguration(cluster)
	if err != nil {
		log.Println("Cannot read cluster configuration", err)
		utils.SendError(writer, err.Error())
		return
	}
	sendConfiguration(writer, configuration)
}

// RegisterCluster - register new cluster.
func (s Server) RegisterCluster(writer http.ResponseWriter, request *http.Request) {
	clusterName, foundName := mux.Vars(request)["cluster"]

	// check parameters provided by client
	if !foundName {
		log.Println("Cluster name is not provided")
		utils.SendError(writer, "Cluster name needs to be specified")
		return
	}

	s.Splunk.LogAction("RegisterCluster", "tester", clusterName)
	err := s.Storage.RegisterNewCluster(clusterName)
	if err != nil {
		log.Println("Cannot create new cluster", err)
		utils.SendInternalServerError(writer, err.Error())
	}
	utils.SendCreated(writer, utils.BuildOkResponse())
}

// GetActiveTriggersForCluster - return list of triggers for single cluster
func (s Server) GetActiveTriggersForCluster(writer http.ResponseWriter, request *http.Request) {
	cluster, found := mux.Vars(request)["cluster"]
	if !found {
		utils.SendError(writer, "Cluster name needs to be specified")
		return
	}

	triggers, err := s.Storage.ListActiveClusterTriggers(cluster)
	if err != nil {
		utils.SendError(writer, err.Error())
		return
	}
	utils.SendResponse(writer, utils.BuildOkResponseWithData("triggers", triggers))
}

// AckTriggerForCluster - ack single cluster's trigger
func (s Server) AckTriggerForCluster(writer http.ResponseWriter, request *http.Request) {
	cluster, found := mux.Vars(request)["cluster"]
	if !found {
		utils.SendError(writer, "Cluster name needs to be specified")
		return
	}

	triggerID, found := mux.Vars(request)["trigger"]
	if !found {
		utils.SendError(writer, "Trigger ID needs to be specified")
		return
	}

	err := s.Storage.AckTrigger(cluster, triggerID)
	if err != nil {
		utils.SendError(writer, err.Error())
		return
	}
	utils.SendResponse(writer, utils.BuildOkResponse())
}
