package expression_test

import (
	"fmt"
	"testing"

	"github.com/Open-Event-Systems/gonjaexpr/expression"

	"github.com/nikolalohinski/gonja/v2"
	"github.com/nikolalohinski/gonja/v2/exec"
	"github.com/nikolalohinski/gonja/v2/loaders"
)

func TestExpression(t *testing.T) {
	cfg := gonja.DefaultConfig
	env := gonja.DefaultEnvironment
	loader := loaders.MustNewMemoryLoader(nil)
	eval := exec.Evaluator{
		Config: cfg,
		Environment: env,
		Loader: loader,
	}
	cases := []struct{expr string; expected any}{
		{"false", false},
		{"true", true},
		{"0", 0},
		{"0.0", 0.0},
		{"123", 123},
		{"test", "test"},
		{"number", 123},
		{"1 and true", true},
		{"0 or false", false},
		{"\"string\"", "string"},
		{"[1, 2, 3] | length", 3},
		{"1 if test == \"test\" else 0", 1},
	}

	ctx := exec.NewContext(map[string]any{
		"test": "test",
		"number": 123,
	})

	for _, testCase := range cases {
		t.Run(fmt.Sprintf("%v", testCase.expr), func (t *testing.T) {
			parsed, err := expression.ParseExpression(cfg, env, eval, testCase.expr)
			if err != nil {
				panic(err)
			}

			res, err := parsed.Evaluate(ctx)
			if err != nil {
				panic(err)
			}

			if res != testCase.expected {
				t.Errorf("expected %v, got %v", testCase.expected, res)
			}
		})
	}
}