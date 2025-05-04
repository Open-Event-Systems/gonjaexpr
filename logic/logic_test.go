package logic_test

import (
	"fmt"
	"testing"

	"github.com/Open-Event-Systems/gonjaexpr/logic"

	"github.com/nikolalohinski/gonja/v2/exec"
)

func TestLogicEval(t *testing.T) {
	cases := []struct {
		expr     logic.Evaluable
		expected any
	}{
		{logic.ValueExpr{0}, 0},
		{logic.ValueExpr{true}, true},
		{logic.ValueExpr{false}, false},
		{logic.ValueExpr{"test"}, "test"},
		{logic.NotExpr{logic.ValueExpr{true}}, false},
		{logic.NotExpr{logic.ValueExpr{""}}, true},
		{logic.AndExpr{}, true},
		{logic.OrExpr{}, false},
		{logic.AndExpr{[]logic.Evaluable{logic.ValueExpr{true}, logic.ValueExpr{false}}}, false},
		{logic.AndExpr{[]logic.Evaluable{logic.ValueExpr{true}, logic.ValueExpr{true}}}, true},
		{logic.OrExpr{[]logic.Evaluable{logic.ValueExpr{true}, logic.ValueExpr{false}}}, true},
		{logic.OrExpr{[]logic.Evaluable{logic.ValueExpr{false}, logic.ValueExpr{false}}}, false},
	}

	ctx := exec.NewContext(nil)

	for _, testCase := range cases {
		t.Run(fmt.Sprintf("%v", testCase.expr), func(t *testing.T) {
			res, err := testCase.expr.Evaluate(ctx)
			if err != nil {
				panic(err)
			}

			if res != testCase.expected {
				t.Errorf("expected %v, got %v", testCase.expected, res)
			}
		})
	}
}
