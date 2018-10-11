package signature

import (
	"fmt"
	"reflect"

	"github.com/pkg/errors"
)

var idMap map[AlgorithmID]Algorithm
var idNameMap map[string]AlgorithmID

func init() {
	idMap = map[AlgorithmID]Algorithm{
		AlgorithmID(1): ed25519{},
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
		existingName := nameOf(existing)
		if existingName != nameOf(al) {
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

// get the name of an Algorithm
func nameOf(al Algorithm) string {
	ty := reflect.TypeOf(al)
	// arbitrary-depth dereference
	for ty.Kind() == reflect.Interface || ty.Kind() == reflect.Ptr {
		ty = ty.Elem()
	}
	return ty.Name()
}

// Get the ID associated with an Algorithm type
func idOf(al Algorithm) (AlgorithmID, error) {
	alName := nameOf(al)
	if len(alName) == 0 {
		return 0, errors.New("anonymous types are not Algorithms")
	}
	id, hasID := idNameMap[alName]
	if hasID {
		return id, nil
	}
	return 0, fmt.Errorf("Supplied algorithm `%s` not in `idMap`", alName)
}
