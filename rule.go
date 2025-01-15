package tenpen

import "github.com/nanozuki/tenpen/internal/lg"

type Rule struct {
	expr   lg.Expr
	engine *Engine
}

func NewRule(rule string) (*Rule, error) {
	return defaultEngine.NewRule(rule)
}

func (r *Rule) Eval(facts string) (string, error) {
	var vals []lg.Expr
	if facts != "" {
		val, err := lg.ExprFromBytes([]byte(facts))
		if err != nil {
			return "", err
		}
		vals = append(vals, val)
	}
	e := lg.NewEvaluator(r.expr, vals, r.engine.funs)
	gotExpr, err := e.Eval(r.expr)
	if err != nil {
		return "", err
	}
	got, err := lg.ExprToBytes(gotExpr)
	if err != nil {
		return "", err
	}
	return string(got), nil
}
