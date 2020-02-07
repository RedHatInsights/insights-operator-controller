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

package tests

import (
	"encoding/json"
	"fmt"
	"github.com/verdverm/frisby"
)

// Cluster represents cluster record in the controller service.
//     ID: unique key
//     Name: cluster GUID in the following format:
//         c8590f31-e97e-4b85-b506-c45ce1911a12
type Cluster struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// ClusterResponse represents default response for cluster request
type ClusterResponse struct {
	Status   string    `json:"status"`
	Clusters []Cluster `json:"clusters"`
}

func readListOfClusters(f *frisby.Frisby) []Cluster {
	f.Get(API_URL + "/client/cluster")
	f.Send()
	f.ExpectStatus(200)
	f.ExpectHeader("Content-Type", "application/json; charset=utf-8")

	response := ClusterResponse{}
	text, err := f.Resp.Content()
	if err != nil {
		f.AddError(err.Error())
	} else {
		json.Unmarshal(text, &response)
	}
	return response.Clusters
}

func createCluster(f *frisby.Frisby, clusterID string, clusterName string) {
	f.Post(API_URL + "client/cluster/" + clusterName)
	f.Send()
	f.ExpectStatus(201)
}

func deleteCluster(f *frisby.Frisby, clusterID string) {
	f.Delete(API_URL + "client/cluster/" + clusterID)
	f.Send()
	f.ExpectHeader("Content-Type", "application/json; charset=utf-8")
}

func compareClusters(f *frisby.Frisby, clusters []Cluster, expected []Cluster) {
	if len(clusters) != len(expected) {
		f.AddError(fmt.Sprintf("%d clusters are expected, but got %d", len(expected), len(clusters)))
		return
	}

	for i := 0; i < len(expected); i++ {
		if clusters[i] != expected[i] {
			f.AddError(fmt.Sprintf("Different cluster info returned: %v != %v", clusters[i], expected[i]))
		}
	}
}

func compareClustersWithoutID(f *frisby.Frisby, clusters []Cluster, expected []Cluster) {
	if len(clusters) != len(expected) {
		f.AddError(fmt.Sprintf("%d clusters are expected, but got %d", len(expected), len(clusters)))
		return
	}

	for i := 0; i < len(expected); i++ {
		// we are not interested in comparing IDs
		clusters[i].ID = 0
		expected[i].ID = 0
		if clusters[i] != expected[i] {
			f.AddError(fmt.Sprintf("Different cluster info returned: %v != %v", clusters[i], expected[i]))
		}
	}
}

func checkInitialListOfClusters() {
	f := frisby.Create("Check the initial list of clusters")

	clusters := readListOfClusters(f)
	expected := []Cluster{
		{0, "00000000-0000-0000-0000-000000000000"},
		{1, "00000000-0000-0000-0000-000000000001"},
		{2, "00000000-0000-0000-0000-000000000002"},
		{3, "00000000-0000-0000-0000-000000000003"},
		{4, "00000000-0000-0000-0000-000000000004"},
	}
	compareClusters(f, clusters, expected)
}

func checkAddCluster() {
	f := frisby.Create("Check adding new cluster")

	clusters := readListOfClusters(f)
	expected := []Cluster{
		{0, "00000000-0000-0000-0000-000000000000"},
		{1, "00000000-0000-0000-0000-000000000001"},
		{2, "00000000-0000-0000-0000-000000000002"},
		{3, "00000000-0000-0000-0000-000000000003"},
		{4, "00000000-0000-0000-0000-000000000004"},
	}
	compareClusters(f, clusters, expected)

	createCluster(f, "50", "00000000-0000-0000-0000-000000000005")

	clusters = readListOfClusters(f)
	expected = []Cluster{
		{0, "00000000-0000-0000-0000-000000000000"},
		{1, "00000000-0000-0000-0000-000000000001"},
		{2, "00000000-0000-0000-0000-000000000002"},
		{3, "00000000-0000-0000-0000-000000000003"},
		{4, "00000000-0000-0000-0000-000000000004"},
		{5, "00000000-0000-0000-0000-000000000005"},
	}
	compareClusters(f, clusters, expected)
}

func checkDeleteCluster() {
	f := frisby.Create("Check deleting existing cluster")

	clusters := readListOfClusters(f)
	expected := []Cluster{
		{0, "00000000-0000-0000-0000-000000000000"},
		{1, "00000000-0000-0000-0000-000000000001"},
		{2, "00000000-0000-0000-0000-000000000002"},
		{3, "00000000-0000-0000-0000-000000000003"},
		{4, "00000000-0000-0000-0000-000000000004"},
		{5, "00000000-0000-0000-0000-000000000005"},
	}
	compareClusters(f, clusters, expected)

	deleteCluster(f, "5")
	f.ExpectStatus(200)

	clusters = readListOfClusters(f)
	expected = []Cluster{
		{0, "00000000-0000-0000-0000-000000000000"},
		{1, "00000000-0000-0000-0000-000000000001"},
		{2, "00000000-0000-0000-0000-000000000002"},
		{3, "00000000-0000-0000-0000-000000000003"},
		{4, "00000000-0000-0000-0000-000000000004"},
	}
	compareClusters(f, clusters, expected)
}

func checkDeleteAnotherCluster() {
	f := frisby.Create("Check deleting another existing cluster")

	clusters := readListOfClusters(f)
	expected := []Cluster{
		{0, "00000000-0000-0000-0000-000000000000"},
		{1, "00000000-0000-0000-0000-000000000001"},
		{2, "00000000-0000-0000-0000-000000000002"},
		{3, "00000000-0000-0000-0000-000000000003"},
		{4, "00000000-0000-0000-0000-000000000004"},
	}
	compareClusters(f, clusters, expected)

	deleteCluster(f, "4")
	f.ExpectStatus(200)

	clusters = readListOfClusters(f)
	expected = []Cluster{
		{0, "00000000-0000-0000-0000-000000000000"},
		{1, "00000000-0000-0000-0000-000000000001"},
		{2, "00000000-0000-0000-0000-000000000002"},
		{3, "00000000-0000-0000-0000-000000000003"},
	}
	compareClusters(f, clusters, expected)
}

func checkDeleteNonexistentCluster() {
	f := frisby.Create("Check deleting nonexistent cluster")

	clusters := readListOfClusters(f)
	expected := []Cluster{
		{0, "00000000-0000-0000-0000-000000000000"},
		{1, "00000000-0000-0000-0000-000000000001"},
		{2, "00000000-0000-0000-0000-000000000002"},
		{3, "00000000-0000-0000-0000-000000000003"},
	}
	compareClusters(f, clusters, expected)

	deleteCluster(f, "40")
	f.ExpectStatus(404)

	clusters = readListOfClusters(f)
	expected = []Cluster{
		{0, "00000000-0000-0000-0000-000000000000"},
		{1, "00000000-0000-0000-0000-000000000001"},
		{2, "00000000-0000-0000-0000-000000000002"},
		{3, "00000000-0000-0000-0000-000000000003"},
	}
	compareClusters(f, clusters, expected)
}

func checkDeleteAllClusters() {
	f := frisby.Create("Check deleting all existing clusters")

	clusters := readListOfClusters(f)
	expected := []Cluster{
		{0, "00000000-0000-0000-0000-000000000000"},
		{1, "00000000-0000-0000-0000-000000000001"},
		{2, "00000000-0000-0000-0000-000000000002"},
		{3, "00000000-0000-0000-0000-000000000003"},
	}
	compareClusters(f, clusters, expected)

	deleteCluster(f, "0")
	f.ExpectStatus(200)
	deleteCluster(f, "1")
	f.ExpectStatus(200)
	deleteCluster(f, "2")
	f.ExpectStatus(200)
	deleteCluster(f, "3")
	f.ExpectStatus(200)

	clusters = readListOfClusters(f)

	expected = []Cluster{}
	compareClusters(f, clusters, expected)
}

func checkCreateNewCluster() {
	f := frisby.Create("Check creating new cluster")

	clusters := readListOfClusters(f)
	expected := []Cluster{
		{0, "00000000-0000-0000-0000-000000000000"},
		{1, "00000000-0000-0000-0000-000000000001"},
		{2, "00000000-0000-0000-0000-000000000002"},
		{3, "00000000-0000-0000-0000-000000000003"},
	}
	compareClusters(f, clusters, expected)

	createCluster(f, "50", "00000000-0000-0000-0000-000000000005")

	clusters = readListOfClusters(f)
	expected = []Cluster{
		{0, "00000000-0000-0000-0000-000000000000"},
		{1, "00000000-0000-0000-0000-000000000001"},
		{2, "00000000-0000-0000-0000-000000000002"},
		{3, "00000000-0000-0000-0000-000000000003"},
		{4, "00000000-0000-0000-0000-000000000005"},
	}
	compareClustersWithoutID(f, clusters, expected)
}

func checkCreateCluster1234() {
	f := frisby.Create("Check creating new cluster")

	clusters := readListOfClusters(f)
	expected := []Cluster{
		{0, "00000000-0000-0000-0000-000000000000"},
		{1, "00000000-0000-0000-0000-000000000001"},
		{2, "00000000-0000-0000-0000-000000000002"},
		{3, "00000000-0000-0000-0000-000000000003"},
		{4, "00000000-0000-0000-0000-000000000005"},
	}
	compareClustersWithoutID(f, clusters, expected)

	createCluster(f, "1234", "00000001-0002-0003-0004-000000000005")

	clusters = readListOfClusters(f)
	expected = []Cluster{
		{0, "00000000-0000-0000-0000-000000000000"},
		{1, "00000000-0000-0000-0000-000000000001"},
		{2, "00000000-0000-0000-0000-000000000002"},
		{3, "00000000-0000-0000-0000-000000000003"},
		{4, "00000000-0000-0000-0000-000000000005"},
		{5, "00000001-0002-0003-0004-000000000005"},
	}
	compareClustersWithoutID(f, clusters, expected)
}

// ClusterTests run all cluster-related REST API tests.
func ClusterTests() {
	checkInitialListOfClusters()
	checkAddCluster()
	checkDeleteCluster()
	checkDeleteAnotherCluster()
	checkDeleteNonexistentCluster()
	// checkDeleteAllClusters() - not used ATM, DB constraint
	checkCreateNewCluster()
	checkCreateCluster1234()
	// TODO:
	// add new cluster with improper name
	// delete a cluster with improper name
}
