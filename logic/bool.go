package logic

import (
	"reflect"
)

// Convert a value to a boolean.
// Follows Python-like rules for truthiness.
func ToBoolean(value any) bool {
	if value == nil {
		return false
	}

	// basic types
	switch v := value.(type) {
	case bool:
		return v
	case int:
		return v != 0
	case string:
		return v != ""
	case float32, float64:
		return v != 0.0
	case []any:
		return len(v) > 0
	case map[any]any:
		return len(v) > 0
	}

	// have to use reflect to test generic slices/maps
	typ := reflect.TypeOf(value)
	refVal := reflect.ValueOf(value)
	switch typ.Kind() {
	case reflect.Slice, reflect.Array, reflect.Map:
		return refVal.Len() > 0
	default:
		return true
	}
}
