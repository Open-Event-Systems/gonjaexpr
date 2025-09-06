package logic

import (
	"reflect"

	"github.com/nikolalohinski/gonja/v2/exec"
)

type Evaluable interface {
	Evaluate(ctx *exec.Context) (any, error)
}

type Value struct {
	Value any
}

func (v *Value) Evaluate(ctx *exec.Context) (any, error) {
	return v.Value, nil
}

// Treat value as a boolean, follows Python-like rules for truthiness.
func ToBoolean(value any) bool {
	switch v := value.(type) {
	case nil:
		return false
	case int, uint:
		return v != 0
	case float32, float64:
		return v != 0
	case bool:
		return v
	case string:
		return len(v) > 0
	}

	vtype := reflect.TypeOf(value)
	refVal := reflect.ValueOf(value)

	switch vtype.Kind() {
	case reflect.Slice, reflect.Array, reflect.Map:
		return refVal.Len() > 0
	default:
		return true
	}
}
