// Copyright 2024 Dolthub, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package dtables

import (
	"github.com/dolthub/go-mysql-server/sql"

	pgtypes "github.com/dolthub/doltgresql/server/types"
)

// getUnscopedDoltDiffSchema returns the schema for the diff table.
func getUnscopedDoltDiffSchema(dbName, tableName string) sql.Schema {
	return []*sql.Column{
		{Name: "commit_hash", Type: pgtypes.Text, Source: tableName, PrimaryKey: true, DatabaseSource: dbName},
		{Name: "table_name", Type: pgtypes.Text, Source: tableName, PrimaryKey: true, DatabaseSource: dbName},
		{Name: "committer", Type: pgtypes.Text, Source: tableName, PrimaryKey: false, DatabaseSource: dbName},
		{Name: "email", Type: pgtypes.Text, Source: tableName, PrimaryKey: false, DatabaseSource: dbName},
		{Name: "date", Type: pgtypes.Timestamp, Source: tableName, PrimaryKey: false, DatabaseSource: dbName},
		{Name: "message", Type: pgtypes.Text, Source: tableName, PrimaryKey: false, DatabaseSource: dbName},
		{Name: "data_change", Type: pgtypes.Bool, Source: tableName, PrimaryKey: false, DatabaseSource: dbName},
		{Name: "schema_change", Type: pgtypes.Bool, Source: tableName, PrimaryKey: false, DatabaseSource: dbName},
	}
}

// getDiffTableName returns the name of the diff table.
func getDiffTableName() string {
	return "diff"
}
