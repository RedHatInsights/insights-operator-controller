package server

import (
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io"
	"log"
	"net/http"
	"os"
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
	log.Printf("Request URL: %s\n", url)
	duration := time.Since(start)
	log.Printf("Time to serve the page: %s\n", duration)

	apiRequests.With(prometheus.Labels{"url": url}).Inc()

	apiResponses.With(prometheus.Labels{"url": url}).Observe(float64(duration.Microseconds()))
}

func mainEndpoint(writer http.ResponseWriter, request *http.Request) {
	start := time.Now()
	io.WriteString(writer, "Hello world!\n")
	countEndpoint(request, start)
}

func getClusters(writer http.ResponseWriter, request *http.Request) {
	start := time.Now()
	io.WriteString(writer, "Clusters\n")
	countEndpoint(request, start)
}

func getConfiguration(writer http.ResponseWriter, request *http.Request) {
	start := time.Now()
	io.WriteString(writer, "getConfiguration\n")
	countEndpoint(request, start)
}

func setConfiguration(writer http.ResponseWriter, request *http.Request) {
	start := time.Now()
	io.WriteString(writer, "setConfiguration\n")
	countEndpoint(request, start)
}

func readConfiguration(writer http.ResponseWriter, request *http.Request) {
	start := time.Now()
	io.WriteString(writer, "operatorConfiguration\n")
	countEndpoint(request, start)
}

func Initialize(address string) {
	log.Println("Initializing HTTP server at", address)
	router := mux.NewRouter().StrictSlash(true)

	// common REST API endpoints
	router.HandleFunc(API_PREFIX, mainEndpoint)

	// REST API endpoints used by client
	router.HandleFunc(API_PREFIX+"client/clusters", getClusters).Methods("GET")
	router.HandleFunc(API_PREFIX+"client/configuration/{id}", getConfiguration).Methods("GET")
	router.HandleFunc(API_PREFIX+"client/configuration/{id}", setConfiguration).Methods("POST")

	// REST API endpoints used by operator
	router.HandleFunc(API_PREFIX+"operator/configuration", readConfiguration).Methods("GET")

	// Prometheus metrics
	router.Handle("/metrics", promhttp.Handler()).Methods("GET")

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal("Unable to initialize HTTP server", err)
		os.Exit(2)
	}
}
