package logic_test

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/Open-Event-Systems/gonjaexpr/logic"
	"github.com/nikolalohinski/gonja/v2"
	"github.com/nikolalohinski/gonja/v2/exec"
	"gopkg.in/yaml.v3"
)

type testCase struct {
	Expected   bool `yaml:"expected"`
	Expression any  `yaml:"expression"`
}

func TestLogic(t *testing.T) {
	eval := &exec.Evaluator{
		Config:      gonja.DefaultConfig,
		Environment: gonja.DefaultEnvironment,
		Loader:      gonja.DefaultLoader,
	}

	ctx := exec.NewContext(map[string]interface{}{
		"t": true,
		"o": 1,
		"f": false,
	})

	run := func(t *testing.T, value any, expected bool) {
		parsed, err := logic.ParseBooleanExpression(eval, value)
		if err != nil {
			panic(err)
		}

		res, err := parsed.Evaluate(ctx)
		if err != nil {
			panic(err)
		}

		asBool := logic.ToBoolean(res)
		if asBool != expected {
			t.Fatalf("expected %v, got %v", expected, asBool)
		}
	}

	entries, err := os.ReadDir("tests")
	if err != nil {
		panic(err)
	}

	for _, entry := range entries {
		f, err := os.Open("tests/" + entry.Name())
		if err != nil {
			panic(err)
		}
		defer f.Close()

		dec := yaml.NewDecoder(f)

		for {
			var val testCase
			err = dec.Decode(&val)
			if err == io.EOF {
				break
			}

			if err != nil {
				panic(err)
			}

			buf := bytes.NewBuffer(nil)
			enc := yaml.NewEncoder(buf)
			err = enc.Encode(val.Expression)
			if err != nil {
				panic(err)
			}

			t.Run(buf.String(), func(t *testing.T) { run(t, val.Expression, val.Expected) })
		}
	}
}
