package server

import (
	"encoding/json"
	"github.com/redhatinsighs/insights-operator-controller/logging"
	"github.com/redhatinsighs/insights-operator-controller/storage"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
)

// Read list of configuration profiles.
func listConfigurationProfiles(writer http.ResponseWriter, request *http.Request, storage storage.Storage) {
	profiles, err := storage.ListConfigurationProfiles()
	if err == nil {
		json.NewEncoder(writer).Encode(profiles)
	} else {
		writer.WriteHeader(http.StatusBadRequest)
		io.WriteString(writer, err.Error())
	}
}

// Read profile specified by its ID
func getConfigurationProfile(writer http.ResponseWriter, request *http.Request, storage storage.Storage) {
	id, err := retrieveIdRequestParameter(request)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		io.WriteString(writer, "Error reading profile ID from request\n")
		return
	}

	profile, err := storage.GetConfigurationProfile(int(id))
	if err == nil {
		json.NewEncoder(writer).Encode(profile)
	} else {
		writer.WriteHeader(http.StatusBadRequest)
		io.WriteString(writer, err.Error())
	}
}

// Create new configuration profile
func newConfigurationProfile(writer http.ResponseWriter, request *http.Request, storage storage.Storage, splunk logging.Client) {
	username, foundUsername := request.URL.Query()["username"]
	description, foundDescription := request.URL.Query()["description"]

	if !foundUsername {
		writer.WriteHeader(http.StatusBadRequest)
		io.WriteString(writer, "User name needs to be specified\n")
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

	splunk.LogAction("NewConfigurationProfile", username[0], string(configuration))
	profiles, err := storage.StoreConfigurationProfile(username[0], description[0], string(configuration))
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		io.WriteString(writer, err.Error())
	} else {
		writer.WriteHeader(http.StatusCreated)
		json.NewEncoder(writer).Encode(profiles)
	}
}

// Delete configuration profile
func deleteConfigurationProfile(writer http.ResponseWriter, request *http.Request, storage storage.Storage, splunk logging.Client) {
	id, err := retrieveIdRequestParameter(request)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		io.WriteString(writer, "Error reading profile ID from request\n")
		return
	}

	splunk.LogAction("DeleteConfigurationProfile", "tester", strconv.Itoa(int(id)))
	profiles, err := storage.DeleteConfigurationProfile(int(id))
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		io.WriteString(writer, err.Error())
	} else {
		writer.WriteHeader(http.StatusOK)
		json.NewEncoder(writer).Encode(profiles)
	}
}

// Change configuration profile
func changeConfigurationProfile(writer http.ResponseWriter, request *http.Request, storage storage.Storage, splunk logging.Client) {
	id, err := retrieveIdRequestParameter(request)
	username, foundUsername := request.URL.Query()["username"]
	description, foundDescription := request.URL.Query()["description"]

	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		io.WriteString(writer, "Error reading profile ID from request\n")
		return
	}

	if !foundUsername {
		writer.WriteHeader(http.StatusBadRequest)
		io.WriteString(writer, "User name needs to be specified\n")
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

	splunk.LogAction("ChangeConfigurationProfile", username[0], string(configuration))
	profiles, err := storage.ChangeConfigurationProfile(int(id), username[0], description[0], string(configuration))
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		io.WriteString(writer, err.Error())
	} else {
		writer.WriteHeader(http.StatusAccepted)
		json.NewEncoder(writer).Encode(profiles)
	}
}
