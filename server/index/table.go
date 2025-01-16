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

package index

import (
	"fmt"

	"github.com/dolthub/dolt/go/libraries/doltcore/sqle"
	"github.com/dolthub/dolt/go/libraries/doltcore/sqle/index"
	"github.com/dolthub/go-mysql-server/sql"
)

// DoltgresTable is a wrapper around a DoltTable that allows for indexing operations to use the correct indexing
// interface for Postgres compatibility.
type DoltgresTable struct {
	*sqle.DoltTable
}

var _ sql.Table = (*DoltgresTable)(nil)
var _ sql.ProjectedTable = (*DoltgresTable)(nil)
var _ sql.IndexSearchableTable = (*DoltgresTable)(nil)

// IndexedAccess implements the sql.IndexSearchableTable interface.
func (dt *DoltgresTable) IndexedAccess(lookup sql.IndexLookup) sql.IndexedTable {
	if dgRanges, ok := lookup.Ranges.(index.DoltgresRangeCollection); ok {
		return &IndexedDoltgresTable{
			IndexedDoltTable: dt.DoltTable.IndexedAccess(lookup).(*sqle.IndexedDoltTable),
			idx:              lookup.Index,
			rc:               dgRanges,
		}
	} else {
		return dt.DoltTable.IndexedAccess(lookup)
	}
}

// LookupForExpressions implements the sql.IndexSearchableTable interface.
func (dt *DoltgresTable) LookupForExpressions(ctx *sql.Context, exprs ...sql.Expression) (sql.IndexLookup, *sql.FuncDepSet, sql.Expression, bool, error) {
	// GMS splits conjunctions before calling LookupForExpressions, however it is not aware of Doltgres expressions and
	// cannot properly split those.
	exprs = SplitConjunctions(exprs)
	allIndexes, err := dt.DoltTable.GetIndexes(ctx)
	if err != nil {
		return sql.IndexLookup{}, nil, nil, false, err
	}
	if len(allIndexes) == 0 {
		return sql.IndexLookup{}, nil, nil, false, nil
	}
	// Specially handle OR expressions here, although we need to build proper support into the index builder
	if len(exprs) == 1 {
		exprs = SplitDisjunction(exprs[0])
		if len(exprs) > 1 {
			var lookup sql.IndexLookup
			for _, expr := range exprs {
				indexBuilder, err := NewIndexBuilder(ctx, allIndexes)
				if err != nil {
					return sql.IndexLookup{}, nil, nil, false, err
				}
				for _, andExpr := range SplitConjunction(expr) {
					indexBuilder.AddExpression(ctx, andExpr)
				}
				if lookup.Index == nil {
					lookup = indexBuilder.GetLookup(ctx)
				} else {
					newLookup := indexBuilder.GetLookup(ctx)
					// If we're looking at two different indexes, then we'll just return nil and do a table scan
					if lookup.Index.ID() != newLookup.Index.ID() || lookup.Index.Table() != newLookup.Index.Table() {
						return sql.IndexLookup{}, nil, nil, false, nil
					}
					lookup.Ranges = append(lookup.Ranges.(index.DoltgresRangeCollection), newLookup.Ranges.(index.DoltgresRangeCollection)...)
				}
			}
			return lookup, nil, nil, true, nil
		}
	}
	indexBuilder, err := NewIndexBuilder(ctx, allIndexes)
	if err != nil {
		return sql.IndexLookup{}, nil, nil, false, err
	}
	for _, expr := range exprs {
		indexBuilder.AddExpression(ctx, expr)
	}
	return indexBuilder.GetLookup(ctx), nil, nil, true, nil
}

// PreciseMatch implements the sql.IndexSearchableTable interface.
func (dt *DoltgresTable) PreciseMatch() bool {
	return false
}

// SkipIndexCosting implements the sql.IndexSearchableTable interface.
func (dt *DoltgresTable) SkipIndexCosting() bool {
	return true
}

// WithProjections implements the sql.ProjectedTable interface.
func (dt *DoltgresTable) WithProjections(colNames []string) sql.Table {
	return &DoltgresTable{dt.DoltTable.WithProjections(colNames).(*sqle.DoltTable)}
}

// IndexedDoltgresTable is a DoltgresTable with an associated index.
type IndexedDoltgresTable struct {
	*sqle.IndexedDoltTable
	idx sql.Index
	rc  index.DoltgresRangeCollection
}

var _ sql.IndexedTable = (*IndexedDoltgresTable)(nil)

// LookupPartitions implements the sql.IndexedTable interface.
func (idt *IndexedDoltgresTable) LookupPartitions(ctx *sql.Context, lookup sql.IndexLookup) (sql.PartitionIter, error) {
	return idt.IndexedDoltTable.LookupPartitions(ctx, lookup)
}

// Partitions implements the sql.Table interface.
func (idt *IndexedDoltgresTable) Partitions(ctx *sql.Context) (sql.PartitionIter, error) {
	return nil, fmt.Errorf("%T: Partitions is invalid on this table", idt)
}

// PartitionRows implements the sql.Table interface.
func (idt *IndexedDoltgresTable) PartitionRows(ctx *sql.Context, partition sql.Partition) (sql.RowIter, error) {
	return idt.IndexedDoltTable.PartitionRows(ctx, partition)
}

// PreciseMatch implements the sql.IndexSearchableTable interface.
func (idt *IndexedDoltgresTable) PreciseMatch() bool {
	for _, rang := range idt.rc {
		if !rang.PreciseMatch {
			return false
		}
	}
	return true
}

// WithProjections implements the sql.ProjectedTable interface.
func (idt *IndexedDoltgresTable) WithProjections(colNames []string) sql.Table {
	return &DoltgresTable{idt.DoltTable.WithProjections(colNames).(*sqle.DoltTable)}
}
