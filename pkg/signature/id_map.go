package signature

// ----- ---- --- -- -
// Copyright 2019, 2020 The Axiom Foundation. All Rights Reserved.
//
// Licensed under the Apache License 2.0 (the "License").  You may not use
// this file except in compliance with the License.  You can obtain a copy
// in the file LICENSE in the source distribution or at
// https://www.apache.org/licenses/LICENSE-2.0.txt
// - -- --- ---- -----


import (
	"fmt"
	"reflect"

	"github.com/ndau/ndaumath/pkg/signature/algorithms/ed25519"
	"github.com/ndau/ndaumath/pkg/signature/algorithms/null"
	"github.com/ndau/ndaumath/pkg/signature/algorithms/secp256k1"
	"github.com/pkg/errors"
)

// re-export package-native algorithms
var (
	Ed25519   = ed25519.Ed25519
	Secp256k1 = secp256k1.Secp256k1
	Null      = null.Null
)

var idMap map[AlgorithmID]Algorithm
var idNameMap map[string]AlgorithmID

func init() {
	idMap = map[AlgorithmID]Algorithm{
		AlgorithmID(0): Null,
		AlgorithmID(1): Ed25519,
		AlgorithmID(2): Secp256k1,
	}
	buildNameMap()
}

// RegisterAlgorithm makes it possible to serialize and deserialize custom Algorithms
//
// If you build a custom Algorithm, you probably want to call this in an init function
// All IDs < 128 are reserved for canonical implementations.
func RegisterAlgorithm(id AlgorithmID, al Algorithm) error {
	if id < AlgorithmID(128) {
		return fmt.Errorf("Reserved algorithm id %d < 128", id)
	}
	if existing, exists := idMap[id]; exists {
		existingName := NameOf(existing)
		if existingName != NameOf(al) {
			return fmt.Errorf("ID %d already in use by %s", id, existingName)
		}
	}

	idMap[id] = al
	buildNameMap()
	return nil
}

// idNameMap is always the inverse of IdMap
func buildNameMap() {
	var id AlgorithmID
	var example Algorithm
	idNameMap = make(map[string]AlgorithmID)
	for id, example = range idMap {
		idNameMap[reflect.TypeOf(example).Name()] = id
	}
}

// NameOf returns the name of an Algorithm
func NameOf(al Algorithm) string {
	ty := reflect.TypeOf(al)
	// arbitrary-depth dereference
	for ty.Kind() == reflect.Interface || ty.Kind() == reflect.Ptr {
		ty = ty.Elem()
	}
	return ty.Name()
}

// Get the ID associated with an Algorithm type
func idOf(al Algorithm) (AlgorithmID, error) {
	// if the algorithm field is nil, then it's the Null algorithm
	if al == nil {
		return 0, nil
	}
	alName := NameOf(al)
	if len(alName) == 0 {
		return 0, errors.New("anonymous types are not Algorithms")
	}
	id, hasID := idNameMap[alName]
	if hasID {
		return id, nil
	}
	return 0, fmt.Errorf("Supplied algorithm `%s` not in `idMap`", alName)
}
