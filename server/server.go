package server

import (
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redhatinsighs/insights-operator-controller/storage"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

const API_PREFIX = "/api/v1/"

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

func retrieveIdRequestParameter(request *http.Request) (int64, error) {
	id_var := mux.Vars(request)["id"]
	return strconv.ParseInt(id_var, 10, 0)
}

func mainEndpoint(writer http.ResponseWriter, request *http.Request) {
	start := time.Now()
	io.WriteString(writer, "Hello world!\n")
	countEndpoint(request, start)
}

func readConfigurationForOperator(writer http.ResponseWriter, request *http.Request) {
	start := time.Now()
	io.WriteString(writer, "readConfigurationForOperator")
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

func Initialize(address string, storage storage.Storage) {
	log.Println("Initializing HTTP server at", address)
	router := mux.NewRouter().StrictSlash(true)
	router.Use(logRequest)

	// common REST API endpoints
	router.HandleFunc(API_PREFIX, mainEndpoint)

	// REST API endpoints used by client

	clientRouter := router.PathPrefix(API_PREFIX + "client").Subrouter()

	// clusters-related operations
	// (handlers are implemented in the file cluster.go)
	clientRouter.HandleFunc("/cluster", func(w http.ResponseWriter, r *http.Request) { getClusters(w, r, storage) }).Methods("GET")
	clientRouter.HandleFunc("/cluster/{id:[0-9]+}/{name}", func(w http.ResponseWriter, r *http.Request) { newCluster(w, r, storage) }).Methods("POST")
	clientRouter.HandleFunc("/cluster/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) { getClusterById(w, r, storage) }).Methods("GET")
	clientRouter.HandleFunc("/cluster/search", func(w http.ResponseWriter, r *http.Request) { searchCluster(w, r, storage) }).Methods("GET")

	// configuration profiles
	clientRouter.HandleFunc("/profile", func(w http.ResponseWriter, r *http.Request) { listConfigurationProfiles(w, r, storage) }).Methods("GET")
	clientRouter.HandleFunc("/profile/{id}", func(w http.ResponseWriter, r *http.Request) { getConfigurationProfile(w, r, storage) }).Methods("GET")
	clientRouter.HandleFunc("/profile/{id}", func(w http.ResponseWriter, r *http.Request) { changeConfigurationProfile(w, r, storage) }).Methods("PUT")
	clientRouter.HandleFunc("/profile", func(w http.ResponseWriter, r *http.Request) { newConfigurationProfile(w, r, storage) }).Methods("POST")
	clientRouter.HandleFunc("/profile/{id}", func(w http.ResponseWriter, r *http.Request) { deleteConfigurationProfile(w, r, storage) }).Methods("DELETE")

	// clusters and its configurations
	clientRouter.HandleFunc("/cluster/{cluster}/configuration", func(w http.ResponseWriter, r *http.Request) { getClusterConfiguration(w, r, storage) }).Methods("GET")
	clientRouter.HandleFunc("/cluster/{cluster}/configuration", func(w http.ResponseWriter, r *http.Request) { newClusterConfiguration(w, r, storage) }).Methods("POST")
	clientRouter.HandleFunc("/cluster/{cluster}/configuration/enable", func(w http.ResponseWriter, r *http.Request) { enableClusterConfiguration(w, r, storage) }).Methods("PUT")
	clientRouter.HandleFunc("/cluster/{cluster}/configuration/disable", func(w http.ResponseWriter, r *http.Request) { disableClusterConfiguration(w, r, storage) }).Methods("PUT")

	// REST API endpoints used by operator
	operatorRouter := router.PathPrefix(API_PREFIX + "operator").Subrouter()
	operatorRouter.HandleFunc("/configuration", readConfigurationForOperator).Methods("GET")

	// Prometheus metrics
	router.Handle("/metrics", promhttp.Handler()).Methods("GET")

	log.Println("Starting HTTP server at", address)

	err := http.ListenAndServe(address, router)
	if err != nil {
		log.Fatal("Unable to initialize HTTP server", err)
		os.Exit(2)
	}
}
