package lg

import (
	"github.com/nanozuki/tenpen/tperr"
)

type Evaluator struct {
	rule Expr
	v    []Expr // v is the stack of values, last one is the runtime value
	f    []Expr // f is the stack of functions, last one is the runtime functions
}

func NewEvaluator(rule Expr, vars []Expr, functions []Expr) *Evaluator {
	e := &Evaluator{
		rule: rule,
		v:    vars,
		f:    functions,
	}
	switch rule.Type() {
	case ExprArray:
		e.v[len(e.v)-1] = Array{}
		e.f[len(e.f)-1] = Array{}
	case ExprObject:
		e.v[len(e.v)-1] = Object{}
		e.f[len(e.f)-1] = Object{}
	default:
		e.v[len(e.v)-1] = Null{}
		e.f[len(e.f)-1] = Null{}
	}
	return e
}

func (e *Evaluator) SubEvaller(scopedValue Expr) *Evaluator {
	return &Evaluator{
		rule: e.rule,
		v:    append(e.v, scopedValue),
		f:    e.f,
	}
}

func (e *Evaluator) Eval(expr Expr) Expr {
	return e.v[len(e.v)-1]
}

func (e *Evaluator) evalInLoc(expr Expr, loc Path) error {
	switch expr := expr.(type) {
	case Null, String, Number, Bool:
		return e.setVal(loc, expr)
	case Array:
		for i, ex := range expr {
			subloc := append(loc, NumberStep(i))
			return e.evalInLoc(ex, subloc)
		}
	case Object:
		for k, v := range expr {
			subloc := append(loc, StringStep(k))
			return e.evalInLoc(v, subloc)
		}
	case ValRef:
		val, err := e.getVal(Path(expr))
		if err != nil {
			return err
		}
		return e.setVal(loc, val)
	case FnRef:
		fn, err := e.getFn(Path(expr))
		if err != nil {
			return err
		}
		return e.setFn(loc, fn)
	case FnCall:
		panic("not implemented")
	case Fn:
		return e.setFn(loc, expr)
	default:
		panic("unreachable")
	}
	return nil
}

func (e *Evaluator) setVal(loc Path, value Expr) error {
	if len(loc) == 0 {
		e.v[len(e.v)-1] = value
		return nil
	}
	return loc.SetTo(e.v[len(e.v)-1], value)
}

func (e *Evaluator) getVal(loc Path) (Expr, error) {
	for i := len(e.v) - 1; i >= 0; i-- {
		if v, err := loc.GetFrom(e.v[i]); err == nil {
			return v, nil
		}
	}
	v, err := loc.GetFrom(e.rule) // TODO: check DAG
	if err != nil {
		return nil, err
	}
	if err := e.evalInLoc(v, loc); err != nil {
		return nil, err
	}
	return loc.GetFrom(e.v[len(e.v)-1])
}

func (e *Evaluator) setFn(loc Path, value Fn) error {
	if len(loc) == 0 {
		e.f[len(e.f)-1] = value
		return nil
	}
	return loc.SetTo(e.f[len(e.f)-1], value)
}

func (e *Evaluator) getFn(loc Path) (Fn, error) {
	for i := len(e.f) - 1; i >= 0; i-- {
		if v, err := loc.GetFrom(e.f[i]); err == nil && v.Type() == ExprFn {
			return v.(Fn), nil
		}
	}
	if f, err := loc.GetFrom(e.rule); err == nil && f.Type() == ExprFn {
		return f.(Fn), err
	}
	return nil, tperr.NoRefError()
}
