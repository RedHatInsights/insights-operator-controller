package server

import (
	"github.com/gorilla/mux"
	"github.com/redhatinsighs/insights-operator-controller/storage"
	"io"
	"net/http"
)

// Read configuration for the operator
func readConfigurationForOperator(writer http.ResponseWriter, request *http.Request, storage storage.Storage) {
	cluster, found := mux.Vars(request)["cluster"]
	if !found {
		writer.WriteHeader(http.StatusBadRequest)
		io.WriteString(writer, "Cluster ID needs to be specified")
		return
	}

	configuration, err := storage.GetClusterActiveConfiguration(cluster)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		io.WriteString(writer, err.Error())
		return
	}
	io.WriteString(writer, configuration)
}
