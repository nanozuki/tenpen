package lg_test

import (
	"testing"

	"github.com/nanozuki/tenpen/internal/lg"
)

func TestEvaluator(t *testing.T) {
	tests := []struct {
		name    string
		rule    string
		vars    []lg.Expr
		fns     []lg.Expr
		want    string
		wantErr error
	}{
		{
			name:    "directly value",
			rule:    `"hello"`,
			vars:    nil,
			fns:     nil,
			want:    `"hello"`,
			wantErr: nil,
		},
		{
			name: "use builtins get one value",
			rule: `["$+", "#a", "#b"]`,
			vars: []lg.Expr{lg.Object{
				"a": lg.Number(1),
				"b": lg.Number(2),
			}},
			fns:     nil,
			want:    "3",
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule, err := lg.Unmarshal([]byte(tt.rule))
			if err != nil {
				t.Errorf("Unmarshal() error = %v", err)
				return
			}
			t.Logf("rule: %v", rule)
			e := lg.NewEvaluator(rule, tt.vars, tt.fns)
			gotExpr, err := e.Eval(rule)
			if err != tt.wantErr {
				t.Errorf("Evaluator.Eval() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			got, err := lg.Marshal(gotExpr)
			if err != nil {
				t.Errorf("Marshal() error = %v", err)
				return
			}
			if string(got) != tt.want {
				t.Errorf("Evaluator.Eval() = %v, want %v", gotExpr, tt.want)
			}
		})
	}
}
