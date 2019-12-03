package utils

import (
	"encoding/json"
	"net/http"
)

func BuildResponse(status string) map[string]interface{} {
	return map[string]interface{}{"status": status}
}

func BuildOkResponse() map[string]interface{} {
	return map[string]interface{}{"status": "ok"}
}

func BuildOkResponseWithData(dataName string, data interface{}) map[string]interface{} {
	resp := map[string]interface{}{"status": "ok"}
	resp[dataName] = data
	return resp
}

func SendResponse(w http.ResponseWriter, data map[string]interface{}) {
	json.NewEncoder(w).Encode(data)
}

func SendCreated(w http.ResponseWriter, data map[string]interface{}) {
	w.WriteHeader(http.StatusCreated)
	SendResponse(w, data)
}

func SendAccepted(w http.ResponseWriter, data map[string]interface{}) {
	w.WriteHeader(http.StatusAccepted)
	SendResponse(w, data)
}

func SendError(w http.ResponseWriter, err string) {
	w.WriteHeader(http.StatusBadRequest)
	SendResponse(w, BuildResponse(err))
}

func SendForbidden(w http.ResponseWriter, err string) {
	w.WriteHeader(http.StatusForbidden)
	SendResponse(w, BuildResponse(err))
}

func SendInternalServerError(w http.ResponseWriter, err string) {
	w.WriteHeader(http.StatusInternalServerError)
	SendResponse(w, BuildResponse(err))
}
