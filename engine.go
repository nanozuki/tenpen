package tenpen

import "github.com/nanozuki/tenpen/internal/lg"

type Engine struct {
	funs []lg.Expr
}

func NewEngine() *Engine {
	return &Engine{
		funs: []lg.Expr{lg.Builtins},
	}
}

func (e *Engine) AddFunction(name string, fn lg.GoFn) {
	if len(e.funs) == 1 {
		e.funs = append(e.funs, lg.Object{})
	}
	last := e.funs[1].(lg.Object)
	last[name] = fn
}

func (e *Engine) AddModule(name string, funcs map[string]lg.GoFn) {
	if len(e.funs) == 1 {
		e.funs = append(e.funs, lg.Object{})
	}
	last := e.funs[1].(lg.Object)
	if _, ok := last[name]; !ok {
		last[name] = lg.Object{}
	}
	mod := last[name].(lg.Object)
	for k, v := range funcs {
		mod[k] = v
	}
}

func (e *Engine) NewRule(rule string) (*Rule, error) {
	expr, err := lg.ExprFromBytes([]byte(rule))
	if err != nil {
		return nil, err
	}
	return &Rule{
		expr:   expr,
		engine: e,
	}, nil
}

var defaultEngine = NewEngine()
