package server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/redhatinsighs/insights-operator-controller/storage"
	"io"
	"io/ioutil"
	"net/http"
)

func getConfiguration(writer http.ResponseWriter, request *http.Request, storage storage.Storage) {
	id, found := mux.Vars(request)["id"]
	if !found {
		writer.WriteHeader(http.StatusBadRequest)
		io.WriteString(writer, "Configuration ID needs to be specified")
		return
	}

	configuration, err := storage.GetClusterConfigurationById(id)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		io.WriteString(writer, err.Error())
		return
	}
	io.WriteString(writer, configuration)
}

func getAllConfigurations(writer http.ResponseWriter, request *http.Request, storage storage.Storage) {
	configuration, err := storage.ListAllClusterConfigurations()
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		io.WriteString(writer, err.Error())
		return
	}
	json.NewEncoder(writer).Encode(configuration)
}

func getClusterConfiguration(writer http.ResponseWriter, request *http.Request, storage storage.Storage) {
	cluster, found := mux.Vars(request)["cluster"]
	if !found {
		writer.WriteHeader(http.StatusBadRequest)
		io.WriteString(writer, "Cluster ID needs to be specified")
		return
	}

	configuration, err := storage.ListClusterConfiguration(cluster)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		io.WriteString(writer, err.Error())
		return
	}
	json.NewEncoder(writer).Encode(configuration)
}

func enableOrDisableConfiguration(writer http.ResponseWriter, request *http.Request, storage storage.Storage, active string) {
	id, found := mux.Vars(request)["id"]
	if !found {
		writer.WriteHeader(http.StatusBadRequest)
		io.WriteString(writer, "Configuration ID needs to be specified")
		return
	}

	err := storage.EnableOrDisableClusterConfigurationById(id, active)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		io.WriteString(writer, err.Error())
		return
	}
	if active == "0" {
		io.WriteString(writer, "disabled")
	} else {
		io.WriteString(writer, "enabled")
	}
}

func enableConfiguration(writer http.ResponseWriter, request *http.Request, storage storage.Storage) {
	enableOrDisableConfiguration(writer, request, storage, "1")
}

func disableConfiguration(writer http.ResponseWriter, request *http.Request, storage storage.Storage) {
	enableOrDisableConfiguration(writer, request, storage, "0")
}

func newClusterConfiguration(writer http.ResponseWriter, request *http.Request, storage storage.Storage) {
	cluster, found := mux.Vars(request)["cluster"]
	if !found {
		writer.WriteHeader(http.StatusBadRequest)
		io.WriteString(writer, "Cluster ID needs to be specified")
		return
	}

	username, foundUsername := request.URL.Query()["username"]
	reason, foundReason := request.URL.Query()["reason"]
	description, foundDescription := request.URL.Query()["description"]

	if !foundUsername {
		writer.WriteHeader(http.StatusBadRequest)
		io.WriteString(writer, "User name needs to be specified\n")
		return
	}

	if !foundReason {
		writer.WriteHeader(http.StatusBadRequest)
		io.WriteString(writer, "Reason needs to be specified\n")
		return
	}

	if !foundDescription {
		writer.WriteHeader(http.StatusBadRequest)
		io.WriteString(writer, "Description needs to be specified\n")
		return
	}

	configuration, err := ioutil.ReadAll(request.Body)
	if err != nil || len(configuration) == 0 {
		writer.WriteHeader(http.StatusBadRequest)
		io.WriteString(writer, "Configuration needs to be provided in the request body")
		return
	}

	configurations, err := storage.CreateClusterConfiguration(cluster, username[0], reason[0], description[0], string(configuration))
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		io.WriteString(writer, err.Error())
		return
	}
	json.NewEncoder(writer).Encode(configurations)
}

func enableClusterConfiguration(writer http.ResponseWriter, request *http.Request, storage storage.Storage) {
	cluster, found := mux.Vars(request)["cluster"]
	if !found {
		writer.WriteHeader(http.StatusBadRequest)
		io.WriteString(writer, "Cluster ID needs to be specified")
		return
	}

	username, foundUsername := request.URL.Query()["username"]
	reason, foundReason := request.URL.Query()["reason"]

	if !foundUsername {
		writer.WriteHeader(http.StatusBadRequest)
		io.WriteString(writer, "User name needs to be specified\n")
		return
	}

	if !foundReason {
		writer.WriteHeader(http.StatusBadRequest)
		io.WriteString(writer, "Reason needs to be specified\n")
		return
	}

	configurations, err := storage.EnableClusterConfiguration(cluster, username[0], reason[0])
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		io.WriteString(writer, err.Error())
		return
	}
	json.NewEncoder(writer).Encode(configurations)
}

func disableClusterConfiguration(writer http.ResponseWriter, request *http.Request, storage storage.Storage) {
	cluster, found := mux.Vars(request)["cluster"]
	if !found {
		writer.WriteHeader(http.StatusBadRequest)
		io.WriteString(writer, "Cluster ID needs to be specified")
		return
	}

	username, foundUsername := request.URL.Query()["username"]
	reason, foundReason := request.URL.Query()["reason"]

	if !foundUsername {
		writer.WriteHeader(http.StatusBadRequest)
		io.WriteString(writer, "User name needs to be specified\n")
		return
	}

	if !foundReason {
		writer.WriteHeader(http.StatusBadRequest)
		io.WriteString(writer, "Reason needs to be specified\n")
		return
	}

	configurations, err := storage.DisableClusterConfiguration(cluster, username[0], reason[0])
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		io.WriteString(writer, err.Error())
		return
	}
	json.NewEncoder(writer).Encode(configurations)
}
