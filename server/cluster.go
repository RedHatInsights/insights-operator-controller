package server

import (
	"encoding/json"
	"github.com/redhatinsighs/insights-operator-controller/storage"
	"io"
	"net/http"
	"strconv"
)

// Read list of all clusters.
func getClusters(writer http.ResponseWriter, request *http.Request, storage storage.Storage) {
	clusters, err := storage.ListOfClusters()
	if err == nil {
		json.NewEncoder(writer).Encode(clusters)
	} else {
		writer.WriteHeader(http.StatusBadRequest)
		io.WriteString(writer, err.Error())
	}
}

// Read cluster specified by its ID.
func getClusterById(writer http.ResponseWriter, request *http.Request, storage storage.Storage) {
	id, err := retrieveIdRequestParameter(request)
	if err == nil {
		cluster, err := storage.GetCluster(int(id))
		if err == nil {
			json.NewEncoder(writer).Encode(cluster)
		} else {
			writer.WriteHeader(http.StatusBadRequest)
			io.WriteString(writer, err.Error())
		}
	} else {
		writer.WriteHeader(http.StatusBadRequest)
		io.WriteString(writer, "Error reading cluster ID from request\n")
	}
}

// Read cluster specified by its ID or name.
func searchCluster(writer http.ResponseWriter, request *http.Request, storage storage.Storage) {
	idParam, foundId := request.URL.Query()["id"]
	nameParam, foundName := request.URL.Query()["name"]

	if foundId {
		id, err := strconv.ParseInt(idParam[0], 10, 0)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			io.WriteString(writer, "Error reading cluster ID from query\n")
		} else {
			cluster, err := storage.GetCluster(int(id))
			if err == nil {
				json.NewEncoder(writer).Encode(cluster)
			} else {
				writer.WriteHeader(http.StatusBadRequest)
				io.WriteString(writer, err.Error())
			}
		}
	} else if foundName {
		cluster, err := storage.GetClusterByName(nameParam[0])
		if err == nil {
			json.NewEncoder(writer).Encode(cluster)
		} else {
			writer.WriteHeader(http.StatusBadRequest)
			io.WriteString(writer, err.Error())
		}
	} else {
		writer.WriteHeader(http.StatusBadRequest)
		io.WriteString(writer, "Either cluster ID or name needs to be specified")
	}
}
