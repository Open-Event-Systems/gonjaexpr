package logic

import (
	"github.com/nikolalohinski/gonja/v2/exec"
)

type LogicAnd struct {
	Exprs []Evaluable
}

func (a *LogicAnd) Evaluate(ctx *exec.Context) (any, error) {
	for _, expr := range a.Exprs {
		res, err := expr.Evaluate(ctx)
		if err != nil {
			return nil, err
		}

		if !ToBoolean(res) {
			return false, nil
		}
	}

	return true, nil
}

type LogicOr struct {
	Exprs []Evaluable
}

func (o *LogicOr) Evaluate(ctx *exec.Context) (any, error) {
	for _, expr := range o.Exprs {
		res, err := expr.Evaluate(ctx)
		if err != nil {
			return nil, err
		}

		if ToBoolean(res) {
			return true, nil
		}
	}

	return false, nil
}

type LogicNot struct {
	Expr Evaluable
}

func (n *LogicNot) Evaluate(ctx *exec.Context) (any, error) {
	res, err := n.Expr.Evaluate(ctx)
	if err != nil {
		return nil, err
	}

	return !ToBoolean(res), nil
}

func Evaluate(ctx *exec.Context, value any) (any, error) {
	switch v := value.(type) {
	case Evaluable:
		return v.Evaluate(ctx)
	default:
		return v, nil
	}
}
