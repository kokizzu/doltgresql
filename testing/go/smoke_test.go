// Copyright 2023 Dolthub, Inc.
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

func TestSmokeTests(t *testing.T) {
	RunScripts(t, []ScriptTest{
		{
			Name: "Commit and diff across branches",
			SetUpScript: []string{
				"CREATE TABLE test (pk BIGINT PRIMARY KEY, v1 BIGINT);",
				"INSERT INTO test VALUES (1, 1), (2, 2);",
				"CALL DOLT_ADD('-A');",
				"CALL DOLT_COMMIT('-m', 'initial commit');",
				"CALL DOLT_BRANCH('other');",
				"UPDATE test SET v1 = 3;",
				"CALL DOLT_ADD('-A');",
				"CALL DOLT_COMMIT('-m', 'commit main');",
				"CALL DOLT_CHECKOUT('other');",
				"UPDATE test SET v1 = 4 WHERE pk = 2;",
				"CALL DOLT_ADD('-A');",
				"CALL DOLT_COMMIT('-m', 'commit other');",
			},
			Assertions: []ScriptTestAssertion{
				{
					Query:            "CALL DOLT_CHECKOUT('main');",
					SkipResultsCheck: true,
				},
				{
					Query: "SELECT * FROM test;",
					Expected: []sql.Row{
						{1, 3},
						{2, 3},
					},
				},
				{
					Query:            "CALL DOLT_CHECKOUT('other');",
					SkipResultsCheck: true,
				},
				{
					Query: "SELECT * FROM test;",
					Expected: []sql.Row{
						{1, 1},
						{2, 4},
					},
				},
				{
					Query: "SELECT from_pk, to_pk, from_v1, to_v1 FROM dolt_diff_test;",
					Expected: []sql.Row{
						{2, 2, 2, 4},
						{nil, 1, nil, 1},
						{nil, 2, nil, 2},
					},
				},
			},
		},
		{
			Name: "Unsupported MySQL statements",
			Assertions: []ScriptTestAssertion{
				{
					Query:       "SHOW CREATE TABLE;",
					ExpectedErr: true,
				},
			},
		},
	})
}
