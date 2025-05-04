package expression

import (
	"fmt"

	"github.com/nikolalohinski/gonja/v2/config"
	"github.com/nikolalohinski/gonja/v2/exec"
	"github.com/nikolalohinski/gonja/v2/nodes"
	"github.com/nikolalohinski/gonja/v2/parser"
	"github.com/nikolalohinski/gonja/v2/tokens"
)

// A Jinja2 expression.
type Expression struct {
	source string
	expr   nodes.Expression
	eval   exec.Evaluator
}

// Parse a Jinja2 expression string.
func ParseExpression(config *config.Config, env *exec.Environment, eval exec.Evaluator, expr string) (*Expression, error) {
	stream := tokens.Lex(fmt.Sprintf("{{ %s }}", expr), config)
	stream.Next()
	parser := parser.NewParser("", stream, config, eval.Loader, env.ControlStructures)
	exprObj, err := parser.ParseExpression()
	if err != nil {
		return nil, err
	}

	return &Expression{
		source: expr,
		expr:   exprObj,
		eval:   eval,
	}, nil
}

func (e *Expression) Evaluate(context *exec.Context) (any, error) {
	eval := exec.Evaluator{
		Config: e.eval.Config,
		Environment: &exec.Environment{
			Filters: e.eval.Environment.Filters,
			ControlStructures: e.eval.Environment.ControlStructures,
			Tests: e.eval.Environment.Tests,
			Methods: e.eval.Environment.Methods,
			Context: context,
		},
		Loader: e.eval.Loader,
	}
	resVal := eval.Eval(e.expr)

	if resVal.IsError() {
		return nil, resVal.Interface().(error)
	}

	return resVal.ToGoSimpleType(false), nil
}
