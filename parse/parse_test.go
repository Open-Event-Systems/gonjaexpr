package parse_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/Open-Event-Systems/gonjaexpr/parse"

	"github.com/nikolalohinski/gonja/v2"
	"github.com/nikolalohinski/gonja/v2/exec"
	"github.com/nikolalohinski/gonja/v2/loaders"
)

func TestParseCondition(t *testing.T) {
	cfg := gonja.DefaultConfig
	env := gonja.DefaultEnvironment
	loader := loaders.MustNewMemoryLoader(nil)
	eval := exec.Evaluator{
		Config: cfg,
		Environment: env,
		Loader: loader,
	}
	cases := []struct{expr string; expected any}{
		{"123", 123.0},
		{"\"number\"", 123},
		{"{\"and\": [\"number\", true]}", true},
		{"{\"or\": [\"test\", false]}", true},
		{"{\"not\": [\"number\"]}", false},
		{"[\"true\", \"false\"]", false},
		{"[\"true\", 1]", true},
	}

	
	ctx := exec.NewContext(map[string]any{
		"test": "test",
		"number": 123,
	})

	for _, testCase := range cases {
		t.Run(fmt.Sprintf("%v", testCase.expr), func (t *testing.T) {
			var fromJson any
			err := json.Unmarshal([]byte(testCase.expr), &fromJson)
			if err != nil {
				panic(err)
			}

			parsed, err := parse.ParseCondition(eval, fromJson)
			
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