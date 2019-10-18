package server

import (
	"encoding/json"
	"github.com/redhatinsighs/insights-operator-controller/storage"
	"io/ioutil"
	"net/http"
	"testing"
)

const API_URL = "http://localhost:8080/api/v1/"

func readListOfClusters(t *testing.T) []storage.Cluster {
	response, err := http.Get(API_URL + "client/cluster")
	if err != nil {
		t.Errorf("Communication error with the server %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected HTTP status 200 OK, got %d", response.StatusCode)
	}

	body, readErr := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	if readErr != nil {
		t.Errorf("Unable to read response body")
	}

	clusters := []storage.Cluster{}
	err = json.Unmarshal(body, &clusters)
	if err != nil {
		t.Error(err)
	}
	return clusters
}

func compareClusters(t *testing.T, clusters []storage.Cluster, expected []storage.Cluster) {
	if len(clusters) != len(expected) {
		t.Errorf("%d clusters are expected, but got %d", len(expected), len(clusters))
	}

	for i := 0; i < len(expected); i++ {
		if clusters[i] != expected[i] {
			t.Errorf("Different cluster info returned: %v != %v", clusters[i], expected[i])
		}
	}
}

func TestGetListOfClusters(t *testing.T) {
	clusters := readListOfClusters(t)

	expected := []storage.Cluster{
		{0, "cluster0"},
		{1, "cluster1"},
		{2, "cluster2"},
		{3, "cluster3"},
		{4, "cluster4"},
	}
	compareClusters(t, clusters, expected)
}

func TestAddCluster(t *testing.T) {
	var client http.Client

	request, err := http.NewRequest("POST", API_URL+"client/cluster/5/cluster5", nil)
	if err != nil {
		t.Errorf("Error creating request %v", err)
	}

	response, err := client.Do(request)
	if response.StatusCode != http.StatusCreated {
		t.Errorf("Expected HTTP status 201 Created, got %d", response.StatusCode)
	}

	clusters := readListOfClusters(t)

	expected := []storage.Cluster{
		{0, "cluster0"},
		{1, "cluster1"},
		{2, "cluster2"},
		{3, "cluster3"},
		{4, "cluster4"},
		{5, "cluster5"},
	}
	compareClusters(t, clusters, expected)
}

func TestDeleteCluster(t *testing.T) {
	var client http.Client

	request, err := http.NewRequest("DELETE", API_URL+"client/cluster/5", nil)
	if err != nil {
		t.Errorf("Error creating request %v", err)
	}

	response, err := client.Do(request)
	if err != nil {
		t.Errorf("Communication error with the server %v", err)
	}
	if response.StatusCode != http.StatusAccepted {
		t.Errorf("Expected HTTP status 202 Accepted, got %d", response.StatusCode)
	}

	clusters := readListOfClusters(t)

	expected := []storage.Cluster{
		{0, "cluster0"},
		{1, "cluster1"},
		{2, "cluster2"},
		{3, "cluster3"},
		{4, "cluster4"},
	}
	compareClusters(t, clusters, expected)
}

func TestDeleteAnotherCluster(t *testing.T) {
	var client http.Client

	request, err := http.NewRequest("DELETE", API_URL+"client/cluster/4", nil)
	if err != nil {
		t.Errorf("Error creating request %v", err)
	}

	response, err := client.Do(request)
	if err != nil {
		t.Errorf("Communication error with the server %v", err)
	}
	if response.StatusCode != http.StatusAccepted {
		t.Errorf("Expected HTTP status 202 Accepted, got %d", response.StatusCode)
	}

	clusters := readListOfClusters(t)

	expected := []storage.Cluster{
		{0, "cluster0"},
		{1, "cluster1"},
		{2, "cluster2"},
		{3, "cluster3"},
	}
	compareClusters(t, clusters, expected)
}
