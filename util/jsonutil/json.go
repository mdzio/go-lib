package jsonutil

import (
	"encoding/json"
	"reflect"
)

// Equal compares two byte slices with JSON content. Only if a and b contain valid
// JSON and the objects are the same, true is returned. Because the package reflect
// is used, no high performance should be expected.
func Equal(a, b []byte) bool {
	var aobj interface{}
	err := json.Unmarshal(a, &aobj)
	if err != nil {
		return false
	}
	var bobj interface{}
	err = json.Unmarshal(b, &bobj)
	if err != nil {
		return false
	}
	return reflect.DeepEqual(aobj, bobj)
}
