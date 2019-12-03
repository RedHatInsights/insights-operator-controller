/*
Copyright © 2019 Red Hat, Inc.

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
	"github.com/gorilla/mux"
	"github.com/redhatinsighs/insights-operator-controller/logging"
	"github.com/redhatinsighs/insights-operator-controller/storage"
	u "github.com/redhatinsighs/insights-operator-controller/utils"
	"log"
	"net/http"
)

// Read configuration for the operator.
func readConfigurationForOperator(writer http.ResponseWriter, request *http.Request, storage storage.Storage) {
	cluster, found := mux.Vars(request)["cluster"]
	if !found {
		log.Println("Cluster name is not provided")
		u.SendError(writer, "Cluster ID needs to be specified")
		return
	}

	configuration, err := storage.GetClusterActiveConfiguration(cluster)
	if err != nil {
		log.Println("Cannot read cluster configuration", err)
		u.SendError(writer, err.Error())
		return
	}
	sendConfiguration(writer, configuration)
}

// Register new cluster.
func registerCluster(writer http.ResponseWriter, request *http.Request, storage storage.Storage, splunk logging.Client) {
	clusterName, foundName := mux.Vars(request)["cluster"]

	// check parameters provided by client
	if !foundName {
		log.Println("Cluster name is not provided")
		u.SendError(writer, "Cluster name needs to be specified")
		return
	}

	splunk.LogAction("RegisterCluster", "tester", clusterName)
	err := storage.RegisterNewCluster(clusterName)
	if err != nil {
		log.Println("Cannot create new cluster", err)
		u.SendInternalServerError(writer, err.Error())
	}
	u.SendCreated(writer, u.BuildOkResponse())
}

func getActiveTriggersForCluster(writer http.ResponseWriter, request *http.Request, storage storage.Storage) {
	cluster, found := mux.Vars(request)["cluster"]
	if !found {
		u.SendError(writer, "Cluster name needs to be specified")
		return
	}

	triggers, err := storage.ListActiveClusterTriggers(cluster)
	if err != nil {
		u.SendError(writer, err.Error())
		return
	}
	u.SendResponse(writer, u.BuildOkResponseWithData("triggers", triggers))
}

func ackTriggerForCluster(writer http.ResponseWriter, request *http.Request, storage storage.Storage) {
	cluster, found := mux.Vars(request)["cluster"]
	if !found {
		u.SendError(writer, "Cluster name needs to be specified")
		return
	}

	triggerID, found := mux.Vars(request)["trigger"]
	if !found {
		u.SendError(writer, "Trigger ID needs to be specified")
		return
	}

	err := storage.AckTrigger(cluster, triggerID)
	if err != nil {
		u.SendError(writer, err.Error())
		return
	}
	u.SendResponse(writer, u.BuildOkResponse())
}
