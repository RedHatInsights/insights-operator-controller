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

package tests

import (
	"encoding/json"
	"fmt"
	"github.com/verdverm/frisby"
)

type Cluster struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func readListOfClusters(f *frisby.Frisby) []Cluster {
	f.Get(API_URL + "/client/cluster")
	f.Send()
	f.ExpectStatus(200)
	f.ExpectHeader("Content-Type", "application/json; charset=utf-8")

	clusters := []Cluster{}
	text, err := f.Resp.Content()
	if err != nil {
		f.AddError(err.Error())
	} else {
		json.Unmarshal(text, &clusters)
	}
	return clusters
}

func createCluster(f *frisby.Frisby, clusterId string, clusterName string) {
	f.Post(API_URL + "client/cluster/" + clusterId + "/" + clusterName)
	f.Send()
	f.ExpectStatus(201)
}

func deleteCluster(f *frisby.Frisby, clusterId string) {
	f.Delete(API_URL + "client/cluster/" + clusterId)
	f.Send()
	f.ExpectStatus(202)
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

func checkInitialListOfClusters() {
	f := frisby.Create("Check the initial list of clusters")

	clusters := readListOfClusters(f)
	expected := []Cluster{
		{0, "cluster0"},
		{1, "cluster1"},
		{2, "cluster2"},
		{3, "cluster3"},
		{4, "cluster4"},
	}
	compareClusters(f, clusters, expected)
}

func checkAddCluster() {
	f := frisby.Create("Check adding new cluster")

	clusters := readListOfClusters(f)
	expected := []Cluster{
		{0, "cluster0"},
		{1, "cluster1"},
		{2, "cluster2"},
		{3, "cluster3"},
		{4, "cluster4"},
	}
	compareClusters(f, clusters, expected)

	createCluster(f, "5", "cluster5")

	clusters = readListOfClusters(f)
	expected = []Cluster{
		{0, "cluster0"},
		{1, "cluster1"},
		{2, "cluster2"},
		{3, "cluster3"},
		{4, "cluster4"},
		{5, "cluster5"},
	}
	compareClusters(f, clusters, expected)
}

func checkDeleteCluster() {
	f := frisby.Create("Check deleting existing cluster")

	clusters := readListOfClusters(f)
	expected := []Cluster{
		{0, "cluster0"},
		{1, "cluster1"},
		{2, "cluster2"},
		{3, "cluster3"},
		{4, "cluster4"},
		{5, "cluster5"},
	}
	compareClusters(f, clusters, expected)

	deleteCluster(f, "5")

	clusters = readListOfClusters(f)
	expected = []Cluster{
		{0, "cluster0"},
		{1, "cluster1"},
		{2, "cluster2"},
		{3, "cluster3"},
		{4, "cluster4"},
	}
	compareClusters(f, clusters, expected)
}

func checkDeleteAnotherCluster() {
	f := frisby.Create("Check deleting another existing cluster")

	clusters := readListOfClusters(f)
	expected := []Cluster{
		{0, "cluster0"},
		{1, "cluster1"},
		{2, "cluster2"},
		{3, "cluster3"},
		{4, "cluster4"},
	}
	compareClusters(f, clusters, expected)

	deleteCluster(f, "4")

	clusters = readListOfClusters(f)
	expected = []Cluster{
		{0, "cluster0"},
		{1, "cluster1"},
		{2, "cluster2"},
		{3, "cluster3"},
	}
	compareClusters(f, clusters, expected)
}

func checkDeleteNonexistentCluster() {
	f := frisby.Create("Check deleting nonexistent cluster")

	clusters := readListOfClusters(f)
	expected := []Cluster{
		{0, "cluster0"},
		{1, "cluster1"},
		{2, "cluster2"},
		{3, "cluster3"},
	}
	compareClusters(f, clusters, expected)

	deleteCluster(f, "40")

	clusters = readListOfClusters(f)
	expected = []Cluster{
		{0, "cluster0"},
		{1, "cluster1"},
		{2, "cluster2"},
		{3, "cluster3"},
	}
	compareClusters(f, clusters, expected)
}

func checkDeleteAllClusters() {
	f := frisby.Create("Check deleting all existing clusters")

	clusters := readListOfClusters(f)
	expected := []Cluster{
		{0, "cluster0"},
		{1, "cluster1"},
		{2, "cluster2"},
		{3, "cluster3"},
	}
	compareClusters(f, clusters, expected)

	deleteCluster(f, "0")
	deleteCluster(f, "1")
	deleteCluster(f, "2")
	deleteCluster(f, "3")

	clusters = readListOfClusters(f)

	expected = []Cluster{}
	compareClusters(f, clusters, expected)
}

func checkCreateNewCluster() {
	f := frisby.Create("Check creating new cluster")

	clusters := readListOfClusters(f)
	expected := []Cluster{
		{0, "cluster0"},
		{1, "cluster1"},
		{2, "cluster2"},
		{3, "cluster3"},
	}
	compareClusters(f, clusters, expected)

	createCluster(f, "5", "cluster5")

	clusters = readListOfClusters(f)
	expected = []Cluster{
		{0, "cluster0"},
		{1, "cluster1"},
		{2, "cluster2"},
		{3, "cluster3"},
		{5, "cluster5"},
	}
	compareClusters(f, clusters, expected)
}

func checkCreateCluster1234() {
	f := frisby.Create("Check creating new cluster")

	clusters := readListOfClusters(f)
	expected := []Cluster{
		{0, "cluster0"},
		{1, "cluster1"},
		{2, "cluster2"},
		{3, "cluster3"},
		{5, "cluster5"},
	}
	compareClusters(f, clusters, expected)

	createCluster(f, "1234", "cluster1234")

	clusters = readListOfClusters(f)
	expected = []Cluster{
		{0, "cluster0"},
		{1, "cluster1"},
		{2, "cluster2"},
		{3, "cluster3"},
		{5, "cluster5"},
		{1234, "cluster1234"},
	}
	compareClusters(f, clusters, expected)
}

func ClusterTests() {
	checkInitialListOfClusters()
	checkAddCluster()
	checkDeleteCluster()
	checkDeleteAnotherCluster()
	checkDeleteNonexistentCluster()
	// checkDeleteAllClusters() - not used ATM, DB constraint
	checkCreateNewCluster()
	checkCreateCluster1234()
}
