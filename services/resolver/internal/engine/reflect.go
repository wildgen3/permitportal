package engine

import "reflect"

// CEL converts to native Go types by reflect.Type, so the two shapes the engine reads
// out of an activation are named once here.
var (
	mapStringStringType = reflect.TypeOf(map[string]string{})
	sliceStringType     = reflect.TypeOf([]string{})
)
