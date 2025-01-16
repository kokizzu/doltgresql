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

package _go

import (
	"testing"

	"github.com/dolthub/go-mysql-server/sql"
)

func TestAsOf(t *testing.T) {
	RunScripts(t, []ScriptTest{
		{
			Name: "Single table",
			SetUpScript: []string{
				`CREATE TABLE test (a INT)`,
				`INSERT INTO test VALUES (1)`,
				`CALL DOLT_COMMIT('-Am', 'new table')`,
				`INSERT INTO test VALUES (2)`,
				`CALL DOLT_COMMIT('-am', 'new row')`,
			},
			Assertions: []ScriptTestAssertion{
				{
					Query: `SELECT * FROM test AS OF 'HEAD^' as t1`,
					Expected: []sql.Row{
						{1},
					},
				},
				{
					Query: `SELECT * FROM test AS OF 'HEAD'`,
					Expected: []sql.Row{
						{1},
						{2},
					},
				},
			},
		},
		{
			Name: "Join",
			SetUpScript: []string{
				`CREATE TABLE test (a INT)`,
				`INSERT INTO test VALUES (1)`,
				`CALL DOLT_COMMIT('-Am', 'new table')`,
				`INSERT INTO test VALUES (2)`,
				`CALL DOLT_COMMIT('-am', 'new row')`,
				`CREATE TABLE test2 (b INT)`,
				`INSERT INTO test2 VALUES (1)`,
				`CALL DOLT_COMMIT('-Am', 'new table')`,
				`INSERT INTO test2 VALUES (2)`,
				`CALL DOLT_COMMIT('-am', 'new row')`,
			},
			Assertions: []ScriptTestAssertion{
				{
					Query: `SELECT * FROM test AS OF 'HEAD~3' t1 join test2 AS OF 'HEAD~' t2 on t1.a = t2.b`,
					Expected: []sql.Row{
						{1, 1},
					},
				},
				{
					Query: `SELECT * FROM test AS OF 'HEAD~3' t1 join test2 AS t2 on t1.a = t2.b`,
					Expected: []sql.Row{
						{1, 1},
					},
				},
				{
					Query: `SELECT * FROM test t1 join test2 AS OF 'HEAD~' AS t2 on t1.a = t2.b`,
					Expected: []sql.Row{
						{1, 1},
					},
				},
				{
					Query: `SELECT * FROM test AS OF 'HEAD~3' t1 cross join test2 AS OF 'HEAD~' t2`,
					Expected: []sql.Row{
						{1, 1},
					},
				},
				{
					Query: `SELECT * FROM test AS OF 'HEAD' t1 cross join test2 AS OF 'HEAD~' t2`,
					Expected: []sql.Row{
						{1, 1},
						{2, 1},
					},
				},
			},
		},
		{
			Name: "Syntax variations", // There are no unit tests for the parser, so test all variations of the AS OF syntax
			SetUpScript: []string{
				`CREATE TABLE test (a INT)`,
				`INSERT INTO test VALUES (1)`,
				`CALL DOLT_COMMIT('-Am', 'new table')`,
				`INSERT INTO test VALUES (2)`,
				`CALL DOLT_COMMIT('-am', 'new row')`,
				`CREATE TABLE test2 (b INT)`,
				`INSERT INTO test2 VALUES (1)`,
				`CALL DOLT_COMMIT('-Am', 'new table')`,
				`INSERT INTO test2 VALUES (2)`,
				`CALL DOLT_COMMIT('-am', 'new row')`,
			},
			Assertions: []ScriptTestAssertion{
				{
					Query: `SELECT * FROM test AS OF 'HEAD~3' AS t1 join test2 AS OF 'HEAD' AS t2 on t1.a = t2.b`,
					Expected: []sql.Row{
						{1, 1},
					},
				},
				{
					Query: `SELECT * FROM test AS OF 'HEAD~3' t1 join test2 AS OF 'HEAD' t2 on t1.a = t2.b`,
					Expected: []sql.Row{
						{1, 1},
					},
				},
				{
					Query: `SELECT * FROM test AS t1 join test2 AS OF 'HEAD~' t2 on t1.a = t2.b`,
					Expected: []sql.Row{
						{1, 1},
					},
				},
				{
					Query: `SELECT * FROM test AS OF SYSTEM TIME 'HEAD~3' join test2 AS t2 on test.a = t2.b`,
					Expected: []sql.Row{
						{1, 1},
					},
				},
			},
		},
	})
}
