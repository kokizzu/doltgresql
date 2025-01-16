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

package output

import "testing"

func TestAlterStatistics(t *testing.T) {
	tests := []QueryParses{
		Unimplemented("ALTER STATISTICS name OWNER TO new_owner"),
		Unimplemented("ALTER STATISTICS name OWNER TO CURRENT_ROLE"),
		Unimplemented("ALTER STATISTICS name OWNER TO CURRENT_USER"),
		Unimplemented("ALTER STATISTICS name OWNER TO SESSION_USER"),
		Unimplemented("ALTER STATISTICS name RENAME TO new_name"),
		Unimplemented("ALTER STATISTICS name SET SCHEMA new_schema"),
		Unimplemented("ALTER STATISTICS name SET STATISTICS 1"),
	}
	RunTests(t, tests)
}
