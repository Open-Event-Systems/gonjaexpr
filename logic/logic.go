package logic

import (
	"github.com/nikolalohinski/gonja/v2/exec"
)

// Something that can be evaluated.
type Evaluable interface {
	Evaluate(context *exec.Context) (any, error)
}

// Literal value.
type ValueExpr struct {
	Value any
}

func (e ValueExpr) Evaluate(context *exec.Context) (any, error) {
	return e.Value, nil
}

// Logic OR
type OrExpr struct {
	Exprs []Evaluable
}

func (expr OrExpr) Evaluate(context *exec.Context) (any, error) {
	for _, e := range expr.Exprs {
		res, err := e.Evaluate(context)
		if err != nil {
			return false, err
		}

		if ToBoolean(res) {
			return true, nil
		}
	}
	return false, nil
}

// Logic AND
type AndExpr struct {
	Exprs []Evaluable
}

func (expr AndExpr) Evaluate(context *exec.Context) (any, error) {
	for _, e := range expr.Exprs {
		res, err := e.Evaluate(context)
		if err != nil {
			return false, err
		}

		if !ToBoolean(res) {
			return false, nil
		}
	}
	return true, nil
}

// Logic NOT
type NotExpr struct {
	Expr Evaluable
}

func (expr NotExpr) Evaluate(context *exec.Context) (any, error) {
	res, err := expr.Expr.Evaluate(context)
	if err != nil {
		return false, err
	}

	return !ToBoolean(res), nil
}
