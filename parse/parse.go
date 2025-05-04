package parse

import (
	"fmt"

	"github.com/Open-Event-Systems/gonjaexpr/expression"
	"github.com/Open-Event-Systems/gonjaexpr/logic"

	"github.com/nikolalohinski/gonja/v2/exec"
)

// Parse a value or expression.
// A slice implicitly becomes an AndExpr.
func ParseCondition(eval exec.Evaluator, value any) (logic.Evaluable, error) {
	asSlice, ok := value.([]any)
	if ok {
		subExprs := make([]logic.Evaluable, 0, len(asSlice))
		for _, e := range asSlice {
			subVal, err := ParseValueOrExpression(eval, e)
			if err != nil {
				return nil, err
			}
			subExprs = append(subExprs, subVal)
		}
		return logic.AndExpr{Exprs: subExprs}, nil
	}
	return ParseValueOrExpression(eval, value)
}

// Parse a value or expression.
func ParseValueOrExpression(eval exec.Evaluator, value any) (logic.Evaluable, error) {
	if value == nil {
		return logic.ValueExpr{Value: nil}, nil
	}

	// basic types
	switch v := value.(type) {
	case string:
		expr, err := expression.ParseExpression(eval.Config, eval.Environment, eval, v)
		if err != nil {
			return nil, err
		}
		return expr, nil
	case bool, int, float32, float64:
		return logic.ValueExpr{Value: v}, nil
	case []any:
		resSlice := make([]logic.Evaluable, 0, len(v))
		for _, e := range v {
			subVal, err := ParseValueOrExpression(eval, e)
			if err != nil {
				return nil, err
			}
			resSlice = append(resSlice, subVal)
		}
		return logic.ValueExpr{Value: resSlice}, nil
	case map[string]any:
		if andVal, andOk := v["and"]; len(v) == 1 && andOk {
			exprs, exprsOk := v["and"].([]any)
			if !exprsOk {
				return nil, fmt.Errorf("not a slice: %v", andVal)
			}
			subExprs := make([]logic.Evaluable, 0, len(exprs))
			for _, e := range exprs {
				subVal, err := ParseValueOrExpression(eval, e)
				if err != nil {
					return nil, err
				}
				subExprs = append(subExprs, subVal)
			}
			return logic.AndExpr{Exprs: subExprs}, nil
		}
		if orVal, orOk := v["or"]; len(v) == 1 && orOk {
			exprs, exprsOk := v["or"].([]any)
			if !exprsOk {
				return nil, fmt.Errorf("not a slice: %v", orVal)
			}
			subExprs := make([]logic.Evaluable, 0, len(exprs))
			for _, e := range exprs {
				subVal, err := ParseValueOrExpression(eval, e)
				if err != nil {
					return nil, err
				}
				subExprs = append(subExprs, subVal)
			}
			return logic.OrExpr{Exprs: subExprs}, nil
		}
		if notVal, notOk := v["not"]; len(v) == 1 && notOk {
			subVal, err := ParseValueOrExpression(eval, notVal)
			if err != nil {
				return nil, err
			}
			return logic.NotExpr{Expr: subVal}, nil
		}
		return logic.ValueExpr{Value: v}, nil
	}

	return logic.ValueExpr{Value: value}, nil
}
