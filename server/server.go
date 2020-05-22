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

package server

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/RedHatInsights/insights-operator-controller/logging"
	"github.com/RedHatInsights/insights-operator-controller/storage"
	"github.com/RedHatInsights/insights-operator-utils/env"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

// Server - configuration of server
type Server struct {
	Address  string
	UseHTTPS bool
	Storage  storage.Storage
	Splunk   logging.Client
	TLSCert  string
	TLSKey   string

	ClusterQuery *storage.ClusterQuery
}

// APIPrefix is appended before all REST API endpoint addresses
var APIPrefix = env.GetEnv("CONTROLLER_PREFIX", "/api/v1/")

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

// checkSplunkOperation checks whether the Splunk operation (write/record) was successful
func checkSplunkOperation(err error) {
	if err != nil {
		log.Println("(not critical) Log into splunk failed", err)
	}
}

func (s Server) createTLSServer(router http.Handler) *http.Server {
	caCert, err := ioutil.ReadFile(s.TLSCert)
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	tlsConfig := &tls.Config{
		ClientCAs:  caCertPool,
		ClientAuth: tls.RequireAndVerifyClientCert,
	}
	tlsConfig.BuildNameToCertificate()
	server := &http.Server{
		Addr:      s.Address,
		TLSConfig: tlsConfig,
		Handler:   router,
	}

	return server
}

func countEndpoint(request *http.Request, start time.Time) {
	url := request.URL.String()
	duration := time.Since(start)
	log.Printf("Time to serve the page: %s\n", duration)

	apiRequests.With(prometheus.Labels{"url": url}).Inc()

	apiResponses.With(prometheus.Labels{"url": url}).Observe(float64(duration.Microseconds()))
}

// retrievePositiveIntRequestParameter gets param with paramName converts to int and checks
// if it's positive
func retrievePositiveIntRequestParameter(request *http.Request, paramName string) (int64, error) {
	strValue, found := mux.Vars(request)[paramName]
	if !found {
		return 0, fmt.Errorf("'%v' param not found", paramName)
	}
	intValue, err := strconv.ParseInt(strValue, 10, 0)
	if err != nil {
		return 0, err
	}
	if intValue < 0 {
		return 0, fmt.Errorf("'%v' param cannot be negative", paramName)
	}
	return intValue, nil
}

func retrieveIDRequestParameter(request *http.Request) (int64, error) {
	return retrievePositiveIntRequestParameter(request, "id")
}

func (s Server) mainEndpoint(writer http.ResponseWriter, request *http.Request) {
	start := time.Now()
	io.WriteString(writer, "Hello world!\n")
	countEndpoint(request, start)
}

func logRequestHandler(writer http.ResponseWriter, request *http.Request, nextHandler http.Handler) {
	log.Println("Request URI: " + request.RequestURI)
	log.Println("Request method: " + request.Method)
	nextHandler.ServeHTTP(writer, request)
}

// LogRequest - middleware for loging requests
func (s Server) LogRequest(nextHandler http.Handler) http.Handler {
	return http.HandlerFunc(
		func(writer http.ResponseWriter, request *http.Request) {
			logRequestHandler(writer, request, nextHandler)
		})
}

// AddDefaultHeaders - middleware for adding headers that should be in any response
func (s Server) AddDefaultHeaders(nextHandler http.Handler) http.Handler {
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
func (s Server) Initialize() {
	log.Println("Environment: ", Environment)
	log.Println("API Prefix: ", APIPrefix)
	log.Println("Initializing HTTP server at", s.Address)
	s.ClusterQuery = storage.NewClusterQuery(s.Storage)
	router := mux.NewRouter().StrictSlash(true)
	router.Use(s.LogRequest)
	if Environment == "production" {
		router.Use(s.JWTAuthentication)
	}
	router.Use(s.AddDefaultHeaders)

	// common REST API endpoints
	router.HandleFunc(APIPrefix, s.mainEndpoint).Methods("GET")

	// REST API endpoints used by client
	clientRouter := router.PathPrefix(APIPrefix + "client").Subrouter()

	// clusters-related operations
	// (handlers are implemented in the file cluster.go)
	clientRouter.HandleFunc("/cluster", s.GetClusters).Methods("GET")
	clientRouter.HandleFunc("/cluster/{name}", s.NewCluster).Methods("POST")
	clientRouter.HandleFunc("/cluster/{id:[0-9]+}", s.GetClusterByID).Methods("GET")
	clientRouter.HandleFunc("/cluster/{id:[0-9]+}", s.DeleteCluster).Methods("DELETE")
	clientRouter.HandleFunc("/cluster/search", s.SearchCluster).Methods("GET")

	// configuration profiles
	// (handlers are implemented in the file profile.go)
	clientRouter.HandleFunc("/profile", s.ListConfigurationProfiles).Methods("GET")
	clientRouter.HandleFunc("/profile/{id}", s.GetConfigurationProfile).Methods("GET")
	clientRouter.HandleFunc("/profile/{id}", s.ChangeConfigurationProfile).Methods("PUT")
	clientRouter.HandleFunc("/profile", s.NewConfigurationProfile).Methods("POST")
	clientRouter.HandleFunc("/profile/{id}", s.DeleteConfigurationProfile).Methods("DELETE")

	// configurations
	// (handlers are implemented in the file configuration.go)
	clientRouter.HandleFunc("/configuration", s.GetAllConfigurations).Methods("GET")
	clientRouter.HandleFunc("/configuration/{id}", s.GetConfiguration).Methods("GET")
	clientRouter.HandleFunc("/configuration/{id}", s.DeleteConfiguration).Methods("DELETE")
	clientRouter.HandleFunc("/configuration/{id}/enable", s.EnableConfiguration).Methods("PUT")
	clientRouter.HandleFunc("/configuration/{id}/disable", s.DisableConfiguration).Methods("PUT")

	// clusters and its configurations
	// (handlers are implemented in the file configuration.go)
	clientRouter.HandleFunc("/cluster/{cluster}/configuration", s.GetClusterConfiguration).Methods("GET")
	clientRouter.HandleFunc("/cluster/{cluster}/configuration/create", s.NewClusterConfiguration).Methods("POST")
	clientRouter.HandleFunc("/cluster/{cluster}/configuration/enable", s.EnableClusterConfiguration).Methods("PUT")
	clientRouter.HandleFunc("/cluster/{cluster}/configuration/disable", s.DisableClusterConfiguration).Methods("PUT")

	// triggers
	clientRouter.HandleFunc("/trigger", s.GetAllTriggers).Methods("GET")
	clientRouter.HandleFunc("/trigger/{id}", s.GetTrigger).Methods("GET")
	clientRouter.HandleFunc("/trigger/{id}", s.DeleteTrigger).Methods("DELETE")
	clientRouter.HandleFunc("/trigger/{id}/activate", s.ActivateTrigger).Methods("PUT", "POST")
	clientRouter.HandleFunc("/trigger/{id}/deactivate", s.DeactivateTrigger).Methods("PUT", "POST")
	clientRouter.HandleFunc("/cluster/{cluster}/trigger", s.GetClusterTriggers).Methods("GET")
	clientRouter.HandleFunc("/cluster/{cluster}/trigger/{trigger}", s.RegisterClusterTrigger).Methods("POST")

	// REST API endpoints used by insights operator
	// (handlers are implemented in the file operator.go)
	operatorRouter := router.PathPrefix(APIPrefix + "operator").Subrouter()
	operatorRouter.HandleFunc("/register/{cluster}", s.RegisterCluster).Methods("GET", "PUT")
	operatorRouter.HandleFunc("/configuration/{cluster}", s.ReadConfigurationForOperator).Methods("GET")
	operatorRouter.HandleFunc("/triggers/{cluster}", s.GetActiveTriggersForCluster).Methods("GET")
	operatorRouter.HandleFunc("/trigger/{cluster}/ack/{trigger}", s.AckTriggerForCluster).Methods("GET", "PUT")

	// Prometheus metrics
	router.Handle("/metrics", promhttp.Handler()).Methods("GET")

	log.Println("Starting HTTP server at", s.Address)

	// try to record the action StartService into Splunk
	err := s.Splunk.Log("Action", "starting service at address "+s.Address)
	// and check whether the Splunk operation was successful
	checkSplunkOperation(err)

	if s.UseHTTPS {
		server := s.createTLSServer(router)
		err = server.ListenAndServeTLS(s.TLSCert, s.TLSKey)
	} else {
		err = http.ListenAndServe(s.Address, router)
	}
	if err != nil {
		log.Fatal("Unable to initialize HTTP server", err)
		s.Splunk.Log("Error", "service can not be started at address "+s.Address)
		os.Exit(2)
	}
}
