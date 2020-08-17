// Cluster handling REST API implementation

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
// https://redhatinsights.github.io/insights-operator-controller/packages/server/cluster.html

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/RedHatInsights/insights-operator-controller/storage"
	"github.com/RedHatInsights/insights-operator-controller/utils"
	"github.com/RedHatInsights/insights-operator-utils/responses"
	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
)

// GetClusters method reads list of all clusters from database and return it to a client.
func (s Server) GetClusters(writer http.ResponseWriter, request *http.Request) {
	// try to retrieve list of clusters from storage
	clusters, err := s.Storage.ListOfClusters()

	// check if the operation has been successful
	if err != nil {
		log.Println("Unable to get list of clusters", err)
		responses.SendInternalServerError(writer, err.Error())
	} else {
		responses.SendResponse(writer, responses.BuildOkResponseWithData("clusters", clusters))
	}
}

// NewCluster method creates a record with new cluster in a database. The updated list of all clusters is returned to client.
func (s Server) NewCluster(writer http.ResponseWriter, request *http.Request) {
	// get the cluster name from request
	clusterName, foundName := mux.Vars(request)["name"]

	if !foundName {
		log.Println("Cluster name is not provided")
		// query parameter 'name' can't be found in request,
		// which might be caused by issue in Gorilla mux (not on client side)
		responses.Send(http.StatusBadRequest, writer, "Cluster name needs to be specified")
		return
	}

	if len(strings.TrimSpace(clusterName)) == 0 {
		responses.Send(http.StatusBadRequest, writer, "Cluster name shouldn't be empty")
		return
	}

	// try to record the action CreateNewCluster into Splunk
	err := s.Splunk.LogAction("CreateNewCluster", "tester", clusterName)
	// and check whether the Splunk operation was successful
	checkSplunkOperation(err)

	//err := storage.CreateNewCluster(clusterId, clusterName)
	err = s.Storage.RegisterNewCluster(clusterName)
	if err != nil {
		log.Println("Cannot create new cluster", err)
		responses.SendInternalServerError(writer, err.Error())
	}

	// try to retrieve list of clusters from storage
	clusters, err := s.Storage.ListOfClusters()

	// check if the operation has been successful
	if err != nil {
		log.Println("Unable to get list of clusters", err)
		responses.SendInternalServerError(writer, err.Error())
	} else {
		responses.SendCreated(writer, responses.BuildOkResponseWithData("clusters", clusters))
	}
}

// GetClusterByID method reads cluster specified by its ID and return it to a client.
func (s Server) GetClusterByID(writer http.ResponseWriter, request *http.Request) {
	// try to retrieve cluster ID from query
	id, err := retrieveIDRequestParameter(request)

	// check if the operation has been successful
	if _, ok := err.(*strconv.NumError); ok {
		log.Println("Bad cluster ID", err)
		responses.Send(http.StatusBadRequest, writer, "Bad cluster ID")
	} else if err != nil {
		log.Println("Cluster ID is not specified in a request", err)
		responses.Send(http.StatusBadRequest, writer, "Error reading cluster ID from request")
	} else {
		cluster, err := s.Storage.GetCluster(int(id))
		if _, ok := err.(*storage.ItemNotFoundError); ok {
			responses.Send(http.StatusNotFound, writer, err.Error())
		} else if err != nil {
			log.Println("Unable to read cluster from database", err)
			responses.SendInternalServerError(writer, err.Error())
		} else {
			responses.SendResponse(writer, responses.BuildOkResponseWithData("cluster", cluster))
		}
	}
}

// DeleteCluster method deletes a cluster
func (s Server) DeleteCluster(writer http.ResponseWriter, request *http.Request) {
	clusterID, err := retrieveIDRequestParameter(request)
	if err != nil {
		log.Println("Cluster ID is not provided or not an integer")
		responses.Send(http.StatusBadRequest, writer, "Cluster ID needs to be specified and to be an integer")
		return
	}

	// try to record the action DeleteCluster into Splunk
	err = s.Splunk.LogAction("DeleteCluster", "tester", fmt.Sprint(clusterID))
	// and check whether the Splunk operation was successful
	checkSplunkOperation(err)

	// delete cluster in database
	err = s.Storage.DeleteCluster(clusterID)

	// check if the storage operation has been successful
	if _, ok := err.(*storage.ItemNotFoundError); ok {
		responses.Send(http.StatusNotFound, writer, err.Error())
	} else if err != nil {
		log.Println("Cannot delete cluster", err)
		responses.SendInternalServerError(writer, err.Error())
	} else {
		clusters, err := s.Storage.ListOfClusters()
		if err != nil {
			log.Println("Unable to get list of clusters", err)
			responses.SendInternalServerError(writer, err.Error())
		} else {
			responses.SendResponse(writer, responses.BuildOkResponseWithData("clusters", clusters))
		}
	}
}

// DeleteClusterByName method deletes a cluster
func (s Server) DeleteClusterByName(writer http.ResponseWriter, request *http.Request) {
	// get the cluster name from request
	clusterName, foundName := mux.Vars(request)["name"]
	if !foundName {
		log.Println("Cluster name is not provided")
		responses.Send(http.StatusBadRequest, writer, "Cluster name needs to be specified")
		return
	}

	// try to record the action DeleteCluster into Splunk
	err := s.Splunk.LogAction("DeleteCluster", "tester", clusterName)
	// and check whether the Splunk operation was successful
	checkSplunkOperation(err)

	// delete cluster in database
	err = s.Storage.DeleteClusterByName(clusterName)

	// check if the storage operation has been successful
	if _, ok := err.(*storage.ItemNotFoundError); ok {
		responses.Send(http.StatusNotFound, writer, err.Error())
	} else if err != nil {
		log.Println("Cannot delete cluster", err)
		responses.Send(http.StatusNotFound, writer, err.Error())
	} else {
		clusters, err := s.Storage.ListOfClusters()
		if err != nil {
			log.Println("Unable to get list of clusters", err)
			responses.SendInternalServerError(writer, err.Error())
		} else {
			responses.SendResponse(writer, responses.BuildOkResponseWithData("clusters", clusters))
		}
	}
}

// SearchCluster method searchs for a cluster specified by its ID or name.
func (s Server) SearchCluster(writer http.ResponseWriter, request *http.Request) {
	var (
		req     storage.SearchClusterRequest
		cluster *storage.Cluster
	)

	err := utils.DecodeValidRequest(&req, SearchClusterTemplate, request.URL.Query())
	if err != nil {
		log.Println(err)
		responses.Send(http.StatusBadRequest, writer, err.Error())
		return
	}
	// either cluster id or its name needs to be specified
	cluster, err = s.ClusterQuery.QueryOne(request.Context(), req)

	// check if the storage operation has been successful
	if err == storage.ErrNoSuchObj {
		responses.Send(http.StatusNotFound, writer, err.Error())
		return
	}

	if err != nil {
		log.Println("Unable to read cluster from database", err)
		responses.SendInternalServerError(writer, err.Error())
		return
	}

	responses.SendResponse(writer, responses.BuildOkResponseWithData("cluster", cluster))
}

// SearchClusterTemplate defines validation rules and messages for SearchCluster
var SearchClusterTemplate = utils.MergeMaps(map[string]interface{}{
	// all acceptable fields are listed
	// case sensitive
	"id":   "int~Error reading and decoding cluster ID from query",
	"name": "",
	"":     "oneOfIdOrName~Either cluster ID or name needs to be specified",
}, utils.PaginationTemplate)

// oneOfIDOrNameValidation validates that id or name is filled
func oneOfIDOrNameValidation(i interface{}, context interface{}) bool {
	// Tag oneOfIdOrName
	v, ok := context.(map[string]interface{})
	if !ok {
		return false
	}
	// the int validation is done next by validator, we are just checking if its filled
	if id, ok := v["id"].(string); ok && len(id) != 0 {
		return true
	}
	if name, ok := v["name"].(string); ok && len(name) != 0 {
		return true
	}
	return false
}

// init function is called during module initialization
func init() {
	govalidator.CustomTypeTagMap.Set("oneOfIdOrName", govalidator.CustomTypeValidator(oneOfIDOrNameValidation))
}
