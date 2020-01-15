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

package storage

import (
	"context"
	"database/sql"
	"reflect"
	"testing"

	sq "github.com/Masterminds/squirrel"
	"github.com/RedHatInsights/insights-operator-controller/utils"
	_ "github.com/lib/pq"           // PostgreSQL database driver
	_ "github.com/mattn/go-sqlite3" // SQLite database driver
)

func TestMap(t *testing.T) {
	c := &Cluster{}
	cqb := NewClusterQuery(Storage{})

	mappedCluster, err := cqb.Map([]ClusterCol{clusterColsDef.ID, clusterColsDef.Name}, c)
	if err != nil {
		t.Errorf("test should suceed, failed with: %s", err)
	}

	expCol := []struct {
		name string
		val  interface{}
	}{
		{name: "ID", val: &c.ID},
		{name: "Name", val: &c.Name},
	}

	for i := range expCol {
		if expCol[i].val != mappedCluster[i] {
			t.Errorf("the column %s is not mapped to proper structure field", expCol[i].name)
		}
	}
}

func TestQueryBy(t *testing.T) {

	tcs := []struct {
		name  string
		req   SearchClusterRequest
		query string
		args  []interface{}
	}{
		{
			name:  "ID",
			req:   SearchClusterRequest{ID: 1},
			query: "SELECT ID, Name FROM cluster WHERE ID = ?",
			args:  []interface{}{1},
		},
		{
			name:  "Name",
			req:   SearchClusterRequest{Name: "name"},
			query: "SELECT ID, Name FROM cluster WHERE Name = ?",
			args:  []interface{}{"name"},
		},
		{
			name:  "NameWithLimitAndOffset",
			req:   SearchClusterRequest{Name: "my cluster", Pagination: utils.Pagination{Offset: 1, Limit: 100}},
			query: "SELECT ID, Name FROM cluster WHERE Name = ? LIMIT 100 OFFSET 1",
			args:  []interface{}{"my cluster"},
		},
	}
	for _, tt := range tcs {
		t.Run(tt.name, func(t *testing.T) {
			storage := &TestStorage{}
			cqb := NewClusterQuery(storage)

			_, err := cqb.QueryOne(context.Background(), tt.req)
			if err != nil {
				t.Errorf("queryone failed with error %s", err)
			}
			if tt.query != storage.lastQuery {
				t.Errorf("expected query %s doesn't match actual query %s", tt.query, storage.lastQuery)
			}
			if !reflect.DeepEqual(tt.args, storage.lastArgs) {
				t.Errorf("expected args %s doesn't match actual args %s", tt.args, storage.lastArgs)
			}
		})
	}
}

type TestStorage struct {
	lastQuery string
	lastArgs  []interface{}
	lastErr   error
}

func (s *TestStorage) Connections() *sql.DB {
	return nil
}

func (s *TestStorage) Placeholder() sq.PlaceholderFormat {
	return sq.Question
}

func (s *TestStorage) QueryOne(ctx context.Context, selectCols []Column, selectBuilder sq.SelectBuilder, mapper func(Column, interface{}) (interface{}, error), res interface{}) error {
	s.lastQuery, s.lastArgs, s.lastErr = selectBuilder.ToSql()
	return s.lastErr
}

var _ Storager = (*TestStorage)(nil)
