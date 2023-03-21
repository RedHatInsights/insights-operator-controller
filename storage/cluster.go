/*
Copyright © 2019, 2020, 2021, 2022, 2023 Red Hat, Inc.

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

// Generated documentation is available at:
// https://godoc.org/github.com/RedHatInsights/insights-operator-controller/storage
//
// Documentation in literate-programming-style is available at:
// https://redhatinsights.github.io/insights-operator-controller/packages/storage/cluster.html

import (
	"context"
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/RedHatInsights/insights-operator-controller/utils"
	_ "github.com/lib/pq"           // PostgreSQL database driver
	_ "github.com/mattn/go-sqlite3" // SQLite database driver
)

// SearchClusterRequest defines type safe SearchCluster request, it is reused and defines request validation tags
type SearchClusterRequest struct {
	utils.Pagination
	ID   int    `schema:"id"`
	Name string `schema:"name"`
}

// ClusterQuery is Sql query model for Cluster
type ClusterQuery struct {
	storage Storager

	Cols          Cols
	selectColumns []ClusterCol
	TableName     string
}

var clusterTableName = "cluster"

// Cols defines which columns exist in Clusters table, just for tyoe safe operations with them
type Cols struct {
	ID   ClusterCol
	Name ClusterCol
}

// ClusterCol is type of cluster column
type ClusterCol Column

// cols is filling Cols structure with actual column names
var clusterColsDef Cols = Cols{
	ID:   ClusterCol("ID"),
	Name: ClusterCol("Name"),
}

// mapCol defines mapping from type safe column to a struct field in result Cluster.
// Used by db row.Scan
func mapCol(col ClusterCol, r *Cluster) (interface{}, error) {
	switch col {
	case clusterColsDef.ID:
		return &r.ID, nil
	case clusterColsDef.Name:
		return &r.Name, nil
	default:
		return nil, fmt.Errorf("unknown col %s", col)
	}
}

// Map creates a list of destination struct fields using columns to select
func (c *ClusterQuery) Map(cols []ClusterCol, r *Cluster) ([]interface{}, error) {
	var mappedCols []interface{}
	for _, c := range cols {
		mc, err := mapCol(c, r)
		if err != nil {
			return nil, err
		}
		mappedCols = append(mappedCols, mc)
	}
	return mappedCols, nil
}

var clusterCols = []ClusterCol{clusterColsDef.ID, clusterColsDef.Name}

// NewClusterQuery creates new ClusterQuery Sql model
func NewClusterQuery(s Storager) *ClusterQuery {
	c := &ClusterQuery{storage: s, Cols: clusterColsDef}
	c.selectColumns = clusterCols
	c.TableName = clusterTableName

	return c
}

// Storager exposes interface for testing
type Storager interface {
	Connections() *sql.DB
	Placeholder() sq.PlaceholderFormat
	QueryOne(context.Context, []Column, sq.SelectBuilder, func(Column, interface{}) (interface{}, error), interface{}) error
}

// ClusterQueryBuilder is just a typed wrapped to build queries more conviniently using only exposed methods
type ClusterQueryBuilder struct {
	sb sq.SelectBuilder
}

// Query exposes typed Cluster queryBuilder
func (c *ClusterQuery) Query() ClusterQueryBuilder {
	qb := ClusterQueryBuilder{}
	builder := sq.StatementBuilder.PlaceholderFormat(c.storage.Placeholder())
	qb.sb = builder.Select(ColNames(c.selectColumns...)...).From(c.TableName).RunWith(c.storage.Connections())
	return qb
}

// Equals add a SQL Where Predicate using AND (if any exists) using Equals (=) operand.
// Ignores Zero values in values
// For example: WHERE col = 3
func (b ClusterQueryBuilder) Equals(col ClusterCol, v interface{}) ClusterQueryBuilder {
	// We will ignore Zero values and skip condition
	if utils.ZeroValue(v) {
		return b
	}
	// squirell QueryBuilder is immutable, so after adding new condition it needs to be assigned back for further changes
	// also only result q := Equals (q from that case) contains change
	// AND this condition
	b.sb = b.sb.Where(sq.Eq{string(col): v})
	return b
}

// WithPaging is setting how many recors (limit) and from which record (offset)
// This can be used for paging.
// It skips Zero values
func (b ClusterQueryBuilder) WithPaging(limit, offset int) ClusterQueryBuilder {
	if limit != 0 {
		b.sb = b.sb.Limit(uint64(limit))
	}
	if offset != 0 {
		b.sb = b.sb.Offset(uint64(offset))
	}
	return b
}

// mapCol defines mapping from type safe column to a struct field in result Cluster.
// Used by db row.Scan
func (c *ClusterQuery) mapCol(storageCol Column, cluster interface{}) (interface{}, error) {
	col := ClusterCol(storageCol)
	r := cluster.(*Cluster)
	var sc interface{}
	switch col {
	case clusterColsDef.ID:
		sc = &r.ID
	case clusterColsDef.Name:
		sc = &r.Name
	default:
		return nil, fmt.Errorf("unknown col %s", col)
	}
	return sc, nil
}

// QueryOne will query DB with generated command and return one row as Cluster
func (c *ClusterQuery) QueryOne(ctx context.Context, req SearchClusterRequest) (*Cluster, error) {
	if c == nil {
		panic("ClusterQuery must not be nil. Make sure the Server has a pointer reference to ClusterQuery by calling NewClusterQuery().")
	}
	qb := c.Query().
		Equals(c.Cols.ID, req.ID).
		Equals(c.Cols.Name, req.Name).
		WithPaging(req.Limit, req.Offset)

	cluster := &Cluster{}
	err := c.storage.QueryOne(ctx, storageCols(c.selectColumns), qb.sb, c.mapCol, cluster)
	if err != nil {
		return nil, err
	}
	return cluster, nil
}

func storageCols(cols []ClusterCol) []Column {
	c := []Column{}
	for _, s := range cols {
		c = append(c, Column(s))
	}
	return c
}

// ColNames Creates a list fo string column names from list of typed ClusterCols
func ColNames(cols ...ClusterCol) []string {
	cns := []string{}
	for _, c := range cols {
		cns = append(cns, string(c))
	}
	return cns
}
