// Utils for REST API

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

package utils

import (
	"encoding/json"
	"net/http"
)

const (
	contentType = "Content-Type"
	appJSON     = "application/json; charset=utf-8"
)

// BuildResponse builds response for RestAPI request
func BuildResponse(status string) map[string]interface{} {
	return map[string]interface{}{"status": status}
}

// BuildOkResponse builds simple "ok" response
func BuildOkResponse() map[string]interface{} {
	return map[string]interface{}{"status": "ok"}
}

// BuildOkResponseWithData builds response with status "ok" and data
func BuildOkResponseWithData(dataName string, data interface{}) map[string]interface{} {
	resp := map[string]interface{}{"status": "ok"}
	resp[dataName] = data
	return resp
}

// SendResponse returns JSON response
func SendResponse(w http.ResponseWriter, data map[string]interface{}) {
	json.NewEncoder(w).Encode(data)
}

// SendCreated returns response with status Created
func SendCreated(w http.ResponseWriter, data map[string]interface{}) {
	w.WriteHeader(http.StatusCreated)
	SendResponse(w, data)
}

// SendAccepted returns response with status Accepted
func SendAccepted(w http.ResponseWriter, data map[string]interface{}) {
	w.WriteHeader(http.StatusAccepted)
	SendResponse(w, data)
}

// SendError returns error response
func SendError(w http.ResponseWriter, err string) {
	w.Header().Set(contentType, appJSON)
	w.WriteHeader(http.StatusBadRequest)
	SendResponse(w, BuildResponse(err))
}

// SendForbidden returns response with status Forbidden
func SendForbidden(w http.ResponseWriter, err string) {
	w.Header().Set(contentType, appJSON)
	w.WriteHeader(http.StatusForbidden)
	SendResponse(w, BuildResponse(err))
}

// SendInternalServerError returns response with status Internal Server Error
func SendInternalServerError(w http.ResponseWriter, err string) {
	w.Header().Set(contentType, appJSON)
	w.WriteHeader(http.StatusInternalServerError)
	SendResponse(w, BuildResponse(err))
}
