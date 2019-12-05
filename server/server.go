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
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redhatinsighs/insights-operator-controller/logging"
	"github.com/redhatinsighs/insights-operator-controller/storage"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

// APIPrefix is appended before all REST API endpoint addresses
const APIPrefix = "/api/v1/"

// Environment CONTROLLER_ENV const for specifying production vs test environment
var Environment = os.Getenv("CONTROLLER_ENV")

var apiRequests = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "api_endpoints_requests",
	Help: "The total number requests per API endpoint",
}, []string{"url"})

var apiResponses = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Name:    "response_time",
	Help:    "Response time",
	Buckets: prometheus.LinearBuckets(0, 20, 20),
}, []string{"url"})

func countEndpoint(request *http.Request, start time.Time) {
	url := request.URL.String()
	duration := time.Since(start)
	log.Printf("Time to serve the page: %s\n", duration)

	apiRequests.With(prometheus.Labels{"url": url}).Inc()

	apiResponses.With(prometheus.Labels{"url": url}).Observe(float64(duration.Microseconds()))
}

func retrieveIDRequestParameter(request *http.Request) (int64, error) {
	idVar := mux.Vars(request)["id"]
	return strconv.ParseInt(idVar, 10, 0)
}

func mainEndpoint(writer http.ResponseWriter, request *http.Request) {
	start := time.Now()
	io.WriteString(writer, "Hello world!\n")
	countEndpoint(request, start)
}

func logRequestHandler(writer http.ResponseWriter, request *http.Request, nextHandler http.Handler) {
	log.Println("Request URI: " + request.RequestURI)
	log.Println("Request method: " + request.Method)
	nextHandler.ServeHTTP(writer, request)
}

func logRequest(nextHandler http.Handler) http.Handler {
	return http.HandlerFunc(
		func(writer http.ResponseWriter, request *http.Request) {
			logRequestHandler(writer, request, nextHandler)
		})
}

func addDefaultHeaders(nextHandler http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if Environment != "production" {
				if origin := r.Header.Get("Origin"); origin != "" {
					w.Header().Set("Access-Control-Allow-Origin", origin)
				}
				w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			nextHandler.ServeHTTP(w, r)
		})
}

// Initialize perform the server initialization
func Initialize(address string, useHTTPS bool, storage storage.Storage, splunk logging.Client) {
	log.Println("Environment: ", Environment)
	log.Println("Initializing HTTP server at", address)
	router := mux.NewRouter().StrictSlash(true)
	router.Use(logRequest)
	if Environment == "production" {
		router.Use(JWTAuthentication)
	}
	router.Use(addDefaultHeaders)

	// common REST API endpoints
	router.HandleFunc(APIPrefix, mainEndpoint).Methods("GET")

	// REST API endpoints used by client
	clientRouter := router.PathPrefix(APIPrefix + "client").Subrouter()

	// clusters-related operations
	// (handlers are implemented in the file cluster.go)
	clientRouter.HandleFunc("/cluster", func(w http.ResponseWriter, r *http.Request) { getClusters(w, r, storage) }).Methods("GET")
	clientRouter.HandleFunc("/cluster/{id:[0-9]+}/{name}", func(w http.ResponseWriter, r *http.Request) { newCluster(w, r, storage, splunk) }).Methods("POST")
	clientRouter.HandleFunc("/cluster/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) { getClusterByID(w, r, storage) }).Methods("GET")
	clientRouter.HandleFunc("/cluster/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) { deleteCluster(w, r, storage, splunk) }).Methods("DELETE")
	clientRouter.HandleFunc("/cluster/search", func(w http.ResponseWriter, r *http.Request) { searchCluster(w, r, storage) }).Methods("GET")

	// configuration profiles
	// (handlers are implemented in the file profile.go)
	clientRouter.HandleFunc("/profile", func(w http.ResponseWriter, r *http.Request) { listConfigurationProfiles(w, r, storage) }).Methods("GET")
	clientRouter.HandleFunc("/profile/{id}", func(w http.ResponseWriter, r *http.Request) { getConfigurationProfile(w, r, storage) }).Methods("GET")
	clientRouter.HandleFunc("/profile/{id}", func(w http.ResponseWriter, r *http.Request) { changeConfigurationProfile(w, r, storage, splunk) }).Methods("PUT")
	clientRouter.HandleFunc("/profile", func(w http.ResponseWriter, r *http.Request) { newConfigurationProfile(w, r, storage, splunk) }).Methods("POST")
	clientRouter.HandleFunc("/profile/{id}", func(w http.ResponseWriter, r *http.Request) { deleteConfigurationProfile(w, r, storage, splunk) }).Methods("DELETE")

	// configurations
	// (handlers are implemented in the file configuration.go)
	clientRouter.HandleFunc("/configuration", func(w http.ResponseWriter, r *http.Request) { getAllConfigurations(w, r, storage) }).Methods("GET")
	clientRouter.HandleFunc("/configuration/{id}", func(w http.ResponseWriter, r *http.Request) { getConfiguration(w, r, storage) }).Methods("GET")
	clientRouter.HandleFunc("/configuration/{id}", func(w http.ResponseWriter, r *http.Request) { deleteConfiguration(w, r, storage, splunk) }).Methods("DELETE")
	clientRouter.HandleFunc("/configuration/{id}/enable", func(w http.ResponseWriter, r *http.Request) { enableConfiguration(w, r, storage, splunk) }).Methods("PUT")
	clientRouter.HandleFunc("/configuration/{id}/disable", func(w http.ResponseWriter, r *http.Request) { disableConfiguration(w, r, storage, splunk) }).Methods("PUT")

	// clusters and its configurations
	// (handlers are implemented in the file configuration.go)
	clientRouter.HandleFunc("/cluster/{cluster}/configuration", func(w http.ResponseWriter, r *http.Request) { getClusterConfiguration(w, r, storage) }).Methods("GET")
	clientRouter.HandleFunc("/cluster/{cluster}/configuration/create", func(w http.ResponseWriter, r *http.Request) { newClusterConfiguration(w, r, storage, splunk) }).Methods("POST")
	clientRouter.HandleFunc("/cluster/{cluster}/configuration/enable", func(w http.ResponseWriter, r *http.Request) { enableClusterConfiguration(w, r, storage, splunk) }).Methods("PUT")
	clientRouter.HandleFunc("/cluster/{cluster}/configuration/disable", func(w http.ResponseWriter, r *http.Request) { disableClusterConfiguration(w, r, storage, splunk) }).Methods("PUT")

	// triggers
	clientRouter.HandleFunc("/trigger", func(w http.ResponseWriter, r *http.Request) { getAllTriggers(w, r, storage) }).Methods("GET")
	clientRouter.HandleFunc("/trigger/{id}", func(w http.ResponseWriter, r *http.Request) { getTrigger(w, r, storage) }).Methods("GET")
	clientRouter.HandleFunc("/trigger/{id}", func(w http.ResponseWriter, r *http.Request) { deleteTrigger(w, r, storage, splunk) }).Methods("DELETE")
	clientRouter.HandleFunc("/trigger/{id}/activate", func(w http.ResponseWriter, r *http.Request) { activateTrigger(w, r, storage, splunk) }).Methods("PUT", "POST")
	clientRouter.HandleFunc("/trigger/{id}/deactivate", func(w http.ResponseWriter, r *http.Request) { deactivateTrigger(w, r, storage, splunk) }).Methods("PUT", "POST")
	clientRouter.HandleFunc("/cluster/{cluster}/trigger", func(w http.ResponseWriter, r *http.Request) { getClusterTriggers(w, r, storage) }).Methods("GET")
	clientRouter.HandleFunc("/cluster/{cluster}/trigger/{trigger}", func(w http.ResponseWriter, r *http.Request) { registerClusterTrigger(w, r, storage, splunk) }).Methods("POST")

	// REST API endpoints used by insights operator
	// (handlers are implemented in the file operator.go)
	operatorRouter := router.PathPrefix(APIPrefix + "operator").Subrouter()
	operatorRouter.HandleFunc("/register/{cluster}", func(w http.ResponseWriter, r *http.Request) { registerCluster(w, r, storage, splunk) }).Methods("GET", "PUT")
	operatorRouter.HandleFunc("/configuration/{cluster}", func(w http.ResponseWriter, r *http.Request) { readConfigurationForOperator(w, r, storage) }).Methods("GET")
	operatorRouter.HandleFunc("/triggers/{cluster}", func(w http.ResponseWriter, r *http.Request) { getActiveTriggersForCluster(w, r, storage) }).Methods("GET")
	operatorRouter.HandleFunc("/trigger/{cluster}/ack/{trigger}", func(w http.ResponseWriter, r *http.Request) { ackTriggerForCluster(w, r, storage) }).Methods("GET", "PUT")

	// Prometheus metrics
	router.Handle("/metrics", promhttp.Handler()).Methods("GET")

	log.Println("Starting HTTP server at", address)

	splunk.Log("Action", "starting service at address "+address)
	var err error

	if useHTTPS {
		err = http.ListenAndServeTLS(address, "server.crt", "server.key", router)
	} else {
		err = http.ListenAndServe(address, router)
	}
	if err != nil {
		log.Fatal("Unable to initialize HTTP server", err)
		splunk.Log("Error", "service can not be started at address "+address)
		os.Exit(2)
	}
}
