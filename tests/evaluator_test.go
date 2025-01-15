package lg_test

import (
	"encoding/json"
	"testing"

	"github.com/nanozuki/tenpen"
)

func TestEvaluator(t *testing.T) {
	tests := []struct {
		name    string
		rule    string
		facts   string
		want    string
		wantErr error
	}{
		{
			name:    "directly value",
			rule:    `"hello"`,
			facts:   "",
			want:    `"hello"`,
			wantErr: nil,
		},
		{
			name:    "use builtins get one value",
			rule:    `["$+", "#a", "#b"]`,
			facts:   `{"a": 1, "b": 2}`,
			want:    "3",
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule, err := tenpen.NewRule(tt.rule)
			if err != nil {
				t.Errorf("NewRule() error = %v", err)
				return
			}
			t.Logf("rule: %v", rule)
			got, err := rule.Eval(tt.facts)
			if err != tt.wantErr {
				t.Errorf("Evaluator.Eval() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				t.Errorf("Marshal() error = %v", err)
				return
			}
			if !isJSONEqual(got, tt.want) {
				t.Errorf("Evaluator.Eval() = %v, want %v", got, tt.want)
			}
		})
	}
}

func isJSONEqual(a, b string) bool {
	var x, y interface{}
	if err := json.Unmarshal([]byte(a), &x); err != nil {
		return false
	}
	if err := json.Unmarshal([]byte(b), &y); err != nil {
		return false
	}
	return isValueEqual(x, y)
}

func isValueEqual(x, y any) bool {
	switch x := x.(type) {
	case nil, bool, float64, string:
		return x == y
	case []interface{}:
		y, ok := y.([]interface{})
		if !ok {
			return false
		}
		if len(x) != len(y) {
			return false
		}
		for i := range x {
			if !isValueEqual(x[i], y[i]) {
				return false
			}
		}
		return true
	case map[string]interface{}:
		y, ok := y.(map[string]interface{})
		if !ok {
			return false
		}
		if len(x) != len(y) {
			return false
		}
		for k, v := range x {
			if !isValueEqual(v, y[k]) {
				return false
			}
		}
		return true
	default:
		panic("unreachable")
	}
}
