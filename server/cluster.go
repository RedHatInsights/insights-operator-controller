// Cluster handling REST API implementation

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
	"strconv"
)

// Read list of all clusters from database and return it to a client.
func getClusters(writer http.ResponseWriter, request *http.Request, storage storage.Storage) {
	clusters, err := storage.ListOfClusters()
	if err != nil {
		log.Println("Unable to get list of clusters", err)
		u.SendInternalServerError(writer, err.Error())
	} else {
		u.SendResponse(writer, u.BuildOkResponseWithData("clusters", clusters))
	}
}

// Create a record with new cluster in a database. The updated list of all clusters is returned to client.
func newCluster(writer http.ResponseWriter, request *http.Request, storage storage.Storage, splunk logging.Client) {
	clusterName, foundName := mux.Vars(request)["name"]

	if !foundName {
		log.Println("Cluster name is not provided")
		u.SendError(writer, "Cluster name needs to be specified")
		return
	}

	splunk.LogAction("CreateNewCluster", "tester", clusterName)
	//err := storage.CreateNewCluster(clusterId, clusterName)
	err := storage.RegisterNewCluster(clusterName)
	if err != nil {
		log.Println("Cannot create new cluster", err)
		u.SendInternalServerError(writer, err.Error())
	}

	clusters, err := storage.ListOfClusters()
	if err != nil {
		log.Println("Unable to get list of clusters", err)
		u.SendError(writer, err.Error())
	} else {
		u.SendCreated(writer, u.BuildOkResponseWithData("clusters", clusters))
	}
}

// Read cluster specified by its ID and return it to a client.
func getClusterByID(writer http.ResponseWriter, request *http.Request, storage storage.Storage) {
	// try to retrieve cluster ID from query
	id, err := retrieveIDRequestParameter(request)
	if err != nil {
		log.Println("Cluster ID is not specified in a request", err)
		u.SendError(writer, "Error reading cluster ID from request")
	} else {
		cluster, err := storage.GetCluster(int(id))
		if err != nil {
			log.Println("Unable to read cluster from database", err)
			u.SendError(writer, err.Error())
		} else {
			u.SendResponse(writer, u.BuildOkResponseWithData("cluster", cluster))
		}
	}
}

// Delete a cluster
func deleteCluster(writer http.ResponseWriter, request *http.Request, storage storage.Storage, splunk logging.Client) {
	clusterID, foundID := mux.Vars(request)["id"]

	// check parameter provided by client
	if !foundID {
		log.Println("Cluster ID is not provided")
		u.SendError(writer, "Cluster ID needs to be specified")
		return
	}

	splunk.LogAction("DeleteCluster", "tester", clusterID)
	err := storage.DeleteCluster(clusterID)
	if err != nil {
		log.Println("Cannot delete cluster", err)
		u.SendError(writer, err.Error())
	}

	clusters, err := storage.ListOfClusters()
	if err != nil {
		log.Println("Unable to get list of clusters", err)
		u.SendError(writer, err.Error())
	} else {
		u.SendAccepted(writer, u.BuildOkResponseWithData("clusters", clusters))
	}
}

// Search for a cluster specified by its ID or name.
func searchCluster(writer http.ResponseWriter, request *http.Request, storage storage.Storage) {
	idParam, foundID := request.URL.Query()["id"]
	nameParam, foundName := request.URL.Query()["name"]

	// either cluster id or its name needs to be specified
	if foundID {
		id, err := strconv.ParseInt(idParam[0], 10, 0)
		if err != nil {
			log.Println("Error reading and decoding cluster ID from query", err)
			u.SendError(writer, "Error reading and decoding cluster ID from query\n")
		} else {
			cluster, err := storage.GetCluster(int(id))
			if err != nil {
				log.Println("Unable to read cluster from database", err)
				u.SendError(writer, err.Error())
			} else {
				u.SendResponse(writer, u.BuildOkResponseWithData("cluster", cluster))
			}
		}
	} else if foundName {
		cluster, err := storage.GetClusterByName(nameParam[0])
		if err != nil {
			log.Println("Unable to read cluster from database", err)
			u.SendError(writer, err.Error())
		} else {
			u.SendResponse(writer, u.BuildOkResponseWithData("cluster", cluster))
		}
	} else {
		u.SendError(writer, "Either cluster ID or name needs to be specified")
	}
}
