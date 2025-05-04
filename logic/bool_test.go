package logic_test

import (
	"fmt"
	"testing"

	"github.com/Open-Event-Systems/gonjaexpr/logic"
)

func TestToBoolean(t *testing.T) {
	cases := []struct{val any; expected bool}{
		{nil, false},
		{false, false},
		{0, false},
		{0.0, false},
		{"", false},
		{make([]string, 0), false},
		{map[string]bool{}, false},
		{true, true},
		{1, true},
		{-1, true},
		{0.1, true},
		{-0.1, true},
		{"false", true},
		{[]bool{false}, true},
		{map[string]bool{"0": false}, true},
	}

	for _, testCase := range cases {
		t.Run(fmt.Sprintf("%v", testCase.val), func (t *testing.T) {
			res := logic.ToBoolean(testCase.val)
			if res != testCase.expected {
				t.Errorf("ToBoolean(%v): expected %v, got %v", testCase.val, testCase.expected, res)
			}
		})
	}
}