package utils

import (
	"encoding/json"
	"net/http"
)

// BuildResponse - build response for RestAPI request
func BuildResponse(status string) map[string]interface{} {
	return map[string]interface{}{"status": status}
}

// BuildOkResponse build simple "ok" response
func BuildOkResponse() map[string]interface{} {
	return map[string]interface{}{"status": "ok"}
}

// BuildOkResponseWithData build response with status "ok" and data
func BuildOkResponseWithData(dataName string, data interface{}) map[string]interface{} {
	resp := map[string]interface{}{"status": "ok"}
	resp[dataName] = data
	return resp
}

// SendResponse - return JSON response
func SendResponse(w http.ResponseWriter, data map[string]interface{}) {
	json.NewEncoder(w).Encode(data)
}

// SendCreated - return response with status Created
func SendCreated(w http.ResponseWriter, data map[string]interface{}) {
	w.WriteHeader(http.StatusCreated)
	SendResponse(w, data)
}

// SendAccepted - return response with status Accepted
func SendAccepted(w http.ResponseWriter, data map[string]interface{}) {
	w.WriteHeader(http.StatusAccepted)
	SendResponse(w, data)
}

// SendError - return error response
func SendError(w http.ResponseWriter, err string) {
	w.WriteHeader(http.StatusBadRequest)
	SendResponse(w, BuildResponse(err))
}

// SendForbidden - return response with status Forbidden
func SendForbidden(w http.ResponseWriter, err string) {
	w.WriteHeader(http.StatusForbidden)
	SendResponse(w, BuildResponse(err))
}

// SendInternalServerError - return response with status Internal Server Error
func SendInternalServerError(w http.ResponseWriter, err string) {
	w.WriteHeader(http.StatusInternalServerError)
	SendResponse(w, BuildResponse(err))
}
