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

package functions

import (
	"time"

	"github.com/dolthub/go-mysql-server/sql"

	"github.com/dolthub/doltgresql/postgres/parser/sem/tree"
	"github.com/dolthub/doltgresql/postgres/parser/timeofday"
	"github.com/dolthub/doltgresql/server/functions/framework"
	pgtypes "github.com/dolthub/doltgresql/server/types"
)

// initTime registers the functions to the catalog.
func initTime() {
	framework.RegisterFunction(time_in)
	framework.RegisterFunction(time_out)
	framework.RegisterFunction(time_recv)
	framework.RegisterFunction(time_send)
	framework.RegisterFunction(timetypmodin)
	framework.RegisterFunction(timetypmodout)
	framework.RegisterFunction(time_cmp)
}

// time_in represents the PostgreSQL function of time type IO input.
var time_in = framework.Function3{
	Name:       "time_in",
	Return:     pgtypes.Time,
	Parameters: [3]*pgtypes.DoltgresType{pgtypes.Cstring, pgtypes.Oid, pgtypes.Int32},
	Strict:     true,
	Callable: func(ctx *sql.Context, _ [4]*pgtypes.DoltgresType, val1, val2, val3 any) (any, error) {
		input := val1.(string)
		//oid := val2.(id.Internal)
		//typmod := val3.(int32)
		// TODO: decode typmod to precision
		p := 6
		//if b.Precision == -1 {
		//	p = b.Precision
		//}
		t, _, err := tree.ParseDTime(nil, input, tree.TimeFamilyPrecisionToRoundDuration(int32(p)))
		if err != nil {
			return nil, err
		}
		return timeofday.TimeOfDay(*t).ToTime(), nil
	},
}

// time_out represents the PostgreSQL function of time type IO output.
var time_out = framework.Function1{
	Name:       "time_out",
	Return:     pgtypes.Cstring,
	Parameters: [1]*pgtypes.DoltgresType{pgtypes.Time},
	Strict:     true,
	Callable: func(ctx *sql.Context, _ [2]*pgtypes.DoltgresType, val any) (any, error) {
		return val.(time.Time).Format("15:04:05.999999999"), nil
	},
}

// time_recv represents the PostgreSQL function of time type IO receive.
var time_recv = framework.Function3{
	Name:       "time_recv",
	Return:     pgtypes.Time,
	Parameters: [3]*pgtypes.DoltgresType{pgtypes.Internal, pgtypes.Oid, pgtypes.Int32},
	Strict:     true,
	Callable: func(ctx *sql.Context, _ [4]*pgtypes.DoltgresType, val1, val2, val3 any) (any, error) {
		data := val1.([]byte)
		//oid := val2.(id.Internal)
		//typmod := val3.(int32)
		// TODO: decode typmod to precision
		if len(data) == 0 {
			return nil, nil
		}
		t := time.Time{}
		if err := t.UnmarshalBinary(data); err != nil {
			return nil, err
		}
		return t, nil
	},
}

// time_send represents the PostgreSQL function of time type IO send.
var time_send = framework.Function1{
	Name:       "time_send",
	Return:     pgtypes.Bytea,
	Parameters: [1]*pgtypes.DoltgresType{pgtypes.Time},
	Strict:     true,
	Callable: func(ctx *sql.Context, _ [2]*pgtypes.DoltgresType, val any) (any, error) {
		return val.(time.Time).MarshalBinary()
	},
}

// timetypmodin represents the PostgreSQL function of time type IO typmod input.
var timetypmodin = framework.Function1{
	Name:       "timetypmodin",
	Return:     pgtypes.Int32,
	Parameters: [1]*pgtypes.DoltgresType{pgtypes.CstringArray},
	Strict:     true,
	Callable: func(ctx *sql.Context, _ [2]*pgtypes.DoltgresType, val any) (any, error) {
		// TODO: typmod=(precision<<16)∣scale
		return nil, nil
	},
}

// timetypmodout represents the PostgreSQL function of time type IO typmod output.
var timetypmodout = framework.Function1{
	Name:       "timetypmodout",
	Return:     pgtypes.Cstring,
	Parameters: [1]*pgtypes.DoltgresType{pgtypes.Int32},
	Strict:     true,
	Callable: func(ctx *sql.Context, _ [2]*pgtypes.DoltgresType, val any) (any, error) {
		// TODO
		// Precision = typmod & 0xFFFF
		// Scale = (typmod >> 16) & 0xFFFF
		return nil, nil
	},
}

// time_cmp represents the PostgreSQL function of time type compare.
var time_cmp = framework.Function2{
	Name:       "time_cmp",
	Return:     pgtypes.Int32,
	Parameters: [2]*pgtypes.DoltgresType{pgtypes.Time, pgtypes.Time},
	Strict:     true,
	Callable: func(ctx *sql.Context, _ [3]*pgtypes.DoltgresType, val1, val2 any) (any, error) {
		ab := val1.(time.Time)
		bb := val2.(time.Time)
		return int32(ab.Compare(bb)), nil
	},
}