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

package typecollection

import (
	"sort"
	"sync"

	"github.com/dolthub/doltgresql/core/id"

	"github.com/dolthub/doltgresql/server/types"
)

// TypeCollection contains a collection of types.
type TypeCollection struct {
	schemaMap map[string]map[string]*types.DoltgresType
	mutex     *sync.RWMutex
}

// Clone returns a new *TypeCollection with the same contents as the original.
func (pgs *TypeCollection) Clone() *TypeCollection {
	pgs.mutex.Lock()
	defer pgs.mutex.Unlock()

	newCollection := &TypeCollection{
		schemaMap: make(map[string]map[string]*types.DoltgresType),
		mutex:     &sync.RWMutex{},
	}
	for schema, nameMap := range pgs.schemaMap {
		if len(nameMap) == 0 {
			continue
		} else if schema == "pg_catalog" {
			// TODO: technically, can create type in pg_catalog schema
			continue
		}
		clonedNameMap := make(map[string]*types.DoltgresType)
		for key, typ := range nameMap {
			clonedNameMap[key] = typ
		}
		newCollection.schemaMap[schema] = clonedNameMap
	}
	return newCollection
}

// CreateType creates a new type.
func (pgs *TypeCollection) CreateType(schName string, typ *types.DoltgresType) error {
	pgs.mutex.Lock()
	defer pgs.mutex.Unlock()

	nameMap, ok := pgs.schemaMap[schName]
	if !ok {
		nameMap = make(map[string]*types.DoltgresType)
		pgs.schemaMap[schName] = nameMap
	}
	if _, ok = nameMap[typ.Name()]; ok {
		return types.ErrTypeAlreadyExists.New(typ.Name())
	}
	nameMap[typ.Name()] = typ
	return nil
}

// DropType drops an existing type.
func (pgs *TypeCollection) DropType(schName, typName string) error {
	pgs.mutex.Lock()
	defer pgs.mutex.Unlock()

	if nameMap, ok := pgs.schemaMap[schName]; ok {
		if _, ok = nameMap[typName]; ok {
			delete(nameMap, typName)
			return nil
		}
	}
	return types.ErrTypeDoesNotExist.New(typName)
}

// GetAllTypes returns a map containing all types in the collection, grouped by the schema they're contained in.
// Each type array is also sorted by the type name. It includes built-in types.
func (pgs *TypeCollection) GetAllTypes() (typesMap map[string][]*types.DoltgresType, schemaNames []string, totalCount int) {
	pgs.mutex.RLock()
	defer pgs.mutex.RUnlock()

	pgs.addSupportedBuiltInTypes()
	typesMap = make(map[string][]*types.DoltgresType)
	for schemaName, nameMap := range pgs.schemaMap {
		schemaNames = append(schemaNames, schemaName)
		typs := make([]*types.DoltgresType, 0, len(nameMap))
		for _, typ := range nameMap {
			typs = append(typs, typ)
		}
		totalCount += len(typs)
		sort.Slice(typs, func(i, j int) bool {
			return typs[i].Name() < typs[j].Name()
		})
		typesMap[schemaName] = typs
	}

	sort.Slice(schemaNames, func(i, j int) bool {
		return schemaNames[i] < schemaNames[j]
	})
	return
}

// GetDomainType returns a domain type with the given schema and name.
// Returns nil if the type cannot be found. It checks for domain type.
func (pgs *TypeCollection) GetDomainType(internalID id.Type) (*types.DoltgresType, bool) {
	t, exists := pgs.GetType(internalID)
	if !exists {
		return nil, exists
	}
	if t.TypType == types.TypeType_Domain {
		return t, exists
	}
	return nil, false
}

// GetType returns the type with the given schema and name.
// Returns nil if the type cannot be found.
func (pgs *TypeCollection) GetType(internalID id.Type) (*types.DoltgresType, bool) {
	pgs.mutex.RLock()
	defer pgs.mutex.RUnlock()

	pgs.addSupportedBuiltInTypes()
	if nameMap, ok := pgs.schemaMap[internalID.SchemaName()]; ok {
		if typ, ok := nameMap[internalID.TypeName()]; ok {
			return typ, true
		}
	}

	return nil, false
}

// GetTypeByID returns the type matching given ID.
func (pgs *TypeCollection) GetTypeByID(internalID id.Id) (*types.DoltgresType, bool) {
	pgs.mutex.RLock()
	defer pgs.mutex.RUnlock()

	pgs.addSupportedBuiltInTypes()
	for _, nameMap := range pgs.schemaMap {
		for _, typ := range nameMap {
			if typ.ID.AsId() == internalID {
				return typ, true
			}
		}
	}
	return nil, false
}

// HasType checks if a type exists with given schema and type name.
func (pgs *TypeCollection) HasType(schema string, typName string) bool {
	pgs.mutex.Lock()
	defer pgs.mutex.Unlock()

	pgs.addSupportedBuiltInTypes()
	nameMap, ok := pgs.schemaMap[schema]
	if !ok {
		nameMap = make(map[string]*types.DoltgresType)
		pgs.schemaMap[schema] = nameMap
	}
	_, ok = nameMap[typName]
	return ok
}

// IterateTypes iterates over all types in the collection.
func (pgs *TypeCollection) IterateTypes(f func(schema string, typ *types.DoltgresType) error) error {
	pgs.mutex.Lock()
	defer pgs.mutex.Unlock()

	pgs.addSupportedBuiltInTypes()
	for schema, nameMap := range pgs.schemaMap {
		for _, t := range nameMap {
			if err := f(schema, t); err != nil {
				return err
			}
		}
	}
	return nil
}

// addSupportedBuiltInTypes adds supported built-in types in the type collection map
// with 'pg_catalog' schema as key. It doesn't add if 'pg_catalog' entry exists in the map.
func (pgs *TypeCollection) addSupportedBuiltInTypes() {
	if _, ok := pgs.schemaMap["pg_catalog"]; !ok {
		// add built-in types
		pgCatTypeMap := make(map[string]*types.DoltgresType)
		for _, t := range types.GetAllBuitInTypes() {
			pgCatTypeMap[t.Name()] = t
		}
		pgs.schemaMap["pg_catalog"] = pgCatTypeMap
	}
}
