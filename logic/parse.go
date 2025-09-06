package logic

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/nikolalohinski/gonja/v2/exec"
)

var InvalidExpressionErr = errors.New("invalid expression")

func ParseBooleanExpression(eval *exec.Evaluator, value any) (Evaluable, error) {
	if value == nil {
		return &LogicAnd{
			Exprs: nil,
		}, nil
	}
	refTyp := reflect.TypeOf(value)
	if refTyp.Kind() == reflect.Slice || refTyp.Kind() == reflect.Array {
		evalArr, err := parseEvalSlice(eval, value)
		if err != nil {
			return nil, err
		}

		return &LogicAnd{
			Exprs: evalArr,
		}, nil
	}

	return ParseBooleanExpressionNoImplicitAnd(eval, value)
}

func ParseBooleanExpressionNoImplicitAnd(eval *exec.Evaluator, value any) (Evaluable, error) {
	if value == nil {
		return &Value{
			Value: nil,
		}, nil
	}
	
	refTyp := reflect.TypeOf(value)

	if refTyp.Kind() == reflect.Map {
		asMap, err := toMap(value)
		if err != nil {
			return nil, err
		}

		return parseMap(eval, asMap)
	}

	if refTyp.Kind() == reflect.String {
		return ParseTemplateExpression(eval, value.(string))
	}

	return &Value{
		Value: value,
	}, nil
}

func parseMap(eval *exec.Evaluator, value map[string]any) (Evaluable, error) {
	if items, ok := value["and"]; ok {
		evalItems, err := parseEvalSlice(eval, items)
		if err != nil {
			return nil, err
		}

		return &LogicAnd{
			Exprs: evalItems,
		}, nil
	} else if items, ok := value["or"]; ok {
		evalItems, err := parseEvalSlice(eval, items)
		if err != nil {
			return nil, err
		}

		return &LogicOr{
			Exprs: evalItems,
		}, nil
	} else if item, ok := value["not"]; ok {
		refTyp := reflect.TypeOf(item)
		if refTyp.Kind() == reflect.Slice || refTyp.Kind() == reflect.Array {
			evalItems, err := parseEvalSlice(eval, item)
			if err != nil {
				return nil, err
			}

			return &LogicNot{
				Expr: &LogicOr{
					Exprs: evalItems,
				},
			}, nil
		} else {
			expr, err := ParseBooleanExpressionNoImplicitAnd(eval, item)
			if err != nil {
				return nil, err
			}

			return &LogicNot{
				Expr: expr,
			}, nil
		}
	} else {
		return nil, fmt.Errorf("%w: invalid expression: %v", InvalidExpressionErr, value)
	}
}

func parseEvalSlice(eval *exec.Evaluator, val any) ([]Evaluable, error) {
	asSlice, err := toSlice(val)
	if err != nil {
		return nil, err
	}

	parsed := make([]Evaluable, 0, len(asSlice))
	for _, item := range asSlice {
		parsedItem, err := ParseBooleanExpressionNoImplicitAnd(eval, item)
		if err != nil {
			return nil, err
		}

		parsed = append(parsed, parsedItem)
	}

	return parsed, nil
}

func toSlice(val any) ([]any, error) {
	refVal := reflect.ValueOf(val)
	if refVal.Kind() != reflect.Slice && refVal.Kind() != reflect.Array {
		return nil, fmt.Errorf("%w: not a slice: %v", InvalidExpressionErr, val)
	}

	var slice []any

	for _, item := range refVal.Seq2() {
		slice = append(slice, item.Interface())
	}

	return slice, nil
}

func toMap(val any) (map[string]any, error) {
	refVal := reflect.ValueOf(val)
	if refVal.Kind() != reflect.Map {
		return nil, fmt.Errorf("%w: not a map: %v", InvalidExpressionErr, val)
	}

	asMap := make(map[string]any, refVal.Len())

	for k, v := range refVal.Seq2() {
		if k.Kind() == reflect.String {
			asMap[k.Interface().(string)] = v.Interface()
		}
	}

	return asMap, nil
}
