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
	"io"
	"net/http"
)

func getAllTriggers(writer http.ResponseWriter, request *http.Request, storage storage.Storage) {
	triggers, err := storage.ListAllTriggers()
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		io.WriteString(writer, err.Error())
		return
	}
	addJsonHeader(writer)
	writer.WriteHeader(http.StatusOK)
	addJson(writer, triggers)
}

func getTrigger(writer http.ResponseWriter, request *http.Request, storage storage.Storage) {
	id, found := mux.Vars(request)["id"]
	if !found {
		writer.WriteHeader(http.StatusBadRequest)
		io.WriteString(writer, "Trigger ID needs to be specified")
		return
	}

	triggers, err := storage.GetTriggerById(id)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		io.WriteString(writer, err.Error())
		return
	}
	addJsonHeader(writer)
	writer.WriteHeader(http.StatusOK)
	addJson(writer, triggers)
}

func deleteTrigger(writer http.ResponseWriter, request *http.Request, storage storage.Storage, splunk logging.Client) {
	id, found := mux.Vars(request)["id"]
	if !found {
		writer.WriteHeader(http.StatusBadRequest)
		io.WriteString(writer, "Trigger ID needs to be specified")
		return
	}

	splunk.LogAction("DeleteTrigger", "tester", id)
	err := storage.DeleteTriggerById(id)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		io.WriteString(writer, err.Error())
		return
	}
	writer.WriteHeader(http.StatusOK)
	io.WriteString(writer, "Deleted")
}

func activateTrigger(writer http.ResponseWriter, request *http.Request, storage storage.Storage, splunk logging.Client) {
	id, found := mux.Vars(request)["id"]
	if !found {
		writer.WriteHeader(http.StatusBadRequest)
		io.WriteString(writer, "Trigger ID needs to be specified")
		return
	}

	splunk.LogAction("ActivateTrigger", "tester", id)
	err := storage.ChangeStateOfTriggerById(id, 1)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		io.WriteString(writer, err.Error())
		return
	}
	addJsonHeader(writer)
	writer.WriteHeader(http.StatusOK)
	addJson(writer, OkStatus)
}

func deactivateTrigger(writer http.ResponseWriter, request *http.Request, storage storage.Storage, splunk logging.Client) {
	id, found := mux.Vars(request)["id"]
	if !found {
		writer.WriteHeader(http.StatusBadRequest)
		io.WriteString(writer, "Trigger ID needs to be specified")
		return
	}

	splunk.LogAction("DeactivateTrigger", "tester", id)
	err := storage.ChangeStateOfTriggerById(id, 0)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		io.WriteString(writer, err.Error())
		return
	}
	addJsonHeader(writer)
	writer.WriteHeader(http.StatusOK)
	addJson(writer, OkStatus)
}

func getClusterTriggers(writer http.ResponseWriter, request *http.Request, storage storage.Storage) {
	cluster, found := mux.Vars(request)["cluster"]
	if !found {
		writer.WriteHeader(http.StatusBadRequest)
		io.WriteString(writer, "Cluster name needs to be specified")
		return
	}

	triggers, err := storage.ListClusterTriggers(cluster)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		io.WriteString(writer, err.Error())
		return
	}
	addJsonHeader(writer)
	writer.WriteHeader(http.StatusOK)
	addJson(writer, triggers)
}

func registerClusterTrigger(writer http.ResponseWriter, request *http.Request, storage storage.Storage, splunk logging.Client) {
	cluster, found := mux.Vars(request)["cluster"]
	if !found {
		writer.WriteHeader(http.StatusBadRequest)
		io.WriteString(writer, "Cluster name needs to be specified")
		return
	}

	triggerType, found := mux.Vars(request)["trigger"]
	if !found {
		writer.WriteHeader(http.StatusBadRequest)
		io.WriteString(writer, "Trigger type needs to be specified")
		return
	}

	username, foundUsername := request.URL.Query()["username"]
	if !foundUsername {
		writer.WriteHeader(http.StatusBadRequest)
		io.WriteString(writer, "User name needs to be specified\n")
		return
	}

	reason, foundReason := request.URL.Query()["reason"]
	if !foundReason {
		writer.WriteHeader(http.StatusBadRequest)
		io.WriteString(writer, "Reason needs to be specified\n")
		return
	}

	link, foundReason := request.URL.Query()["link"]
	if !foundReason {
		writer.WriteHeader(http.StatusBadRequest)
		io.WriteString(writer, "Link needs to be specified\n")
		return
	}

	splunk.LogTriggerAction("RegisterTrigger", username[0], cluster, triggerType)
	err := storage.NewTrigger(cluster, triggerType, username[0], reason[0], link[0])
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		io.WriteString(writer, err.Error())
		return
	}
	writer.WriteHeader(http.StatusOK)
}
