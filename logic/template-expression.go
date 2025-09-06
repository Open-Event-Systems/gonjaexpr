package logic

import (
	"fmt"
	"log"

	"github.com/nikolalohinski/gonja/v2/exec"
	"github.com/nikolalohinski/gonja/v2/nodes"
	"github.com/nikolalohinski/gonja/v2/parser"
	"github.com/nikolalohinski/gonja/v2/tokens"
)

type TemplateExpression struct {
	node nodes.Expression
	eval *exec.Evaluator
}

// Parse an expression
func ParseTemplateExpression(eval *exec.Evaluator, expr string) (*TemplateExpression, error) {
	tmpl := fmt.Sprintf("{{ %s }}", expr)
	stream := tokens.Lex(tmpl, eval.Config)

	stream.Next() // read brackets

	parser := parser.NewParser("", stream, eval.Config, eval.Loader, eval.Environment.ControlStructures)
	exprObj, err := parser.ParseExpression()
	if err != nil {
		return nil, err
	}

	return &TemplateExpression{
		node: exprObj,
		eval: eval,
	}, nil
}

func (e *TemplateExpression) Evaluate(ctx *exec.Context) (any, error) {
	env := *e.eval.Environment
	env.Context = ctx
	evalWithCtx := exec.Evaluator{
		Config:      e.eval.Config,
		Environment: &env,
		Loader:      e.eval.Loader,
	}
	log.Printf("preerr %v", e.node)
	resVal := evalWithCtx.Eval(e.node)
	log.Printf("err %v", resVal.Val)
	if resVal.IsError() {
		return nil, resVal.Interface().(error)
	}

	return resVal.ToGoSimpleType(false), nil
}
