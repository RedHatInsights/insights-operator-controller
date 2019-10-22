package server

import (
	"github.com/gorilla/mux"
	"github.com/redhatinsighs/insights-operator-controller/storage"
	"io"
	"log"
	"net/http"
)

// Read configuration for the operator
func readConfigurationForOperator(writer http.ResponseWriter, request *http.Request, storage storage.Storage) {
	cluster, found := mux.Vars(request)["cluster"]
	if !found {
		log.Println("Cluster name is not provided")
		writer.WriteHeader(http.StatusBadRequest)
		io.WriteString(writer, "Cluster ID needs to be specified")
		return
	}

	configuration, err := storage.GetClusterActiveConfiguration(cluster)
	if err != nil {
		log.Println("Cannot read cluster configuration", err)
		writer.WriteHeader(http.StatusBadRequest)
		io.WriteString(writer, err.Error())
		return
	}
	writer.WriteHeader(http.StatusOK)
	io.WriteString(writer, configuration)
}

func registerCluster(writer http.ResponseWriter, request *http.Request, storage storage.Storage) {
	clusterName, foundName := mux.Vars(request)["cluster"]

	// check parameters provided by client
	if !foundName {
		log.Println("Cluster name is not provided")
		writer.WriteHeader(http.StatusBadRequest)
		io.WriteString(writer, "Cluster name needs to be specified")
		return
	}

	err := storage.RegisterNewCluster(clusterName)
	if err != nil {
		log.Println("Cannot create new cluster", err)
		writer.WriteHeader(http.StatusInternalServerError)
		io.WriteString(writer, err.Error())
	}
	writer.WriteHeader(http.StatusOK)
	io.WriteString(writer, "Registered")
}
