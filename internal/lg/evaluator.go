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

func (e *Evaluator) SubEvaller(scopedValue Expr) Evaller {
	return &Evaluator{
		rule: e.rule,
		v:    append(e.v, scopedValue, Null{}),
		f:    append(e.f, Null{}),
	}
}

func (e *Evaluator) Eval(expr Expr) (Expr, error) {
	err := e.evalInLoc(expr, []Step{})
	if err != nil {
		return nil, err
	}
	return e.v[len(e.v)-1], nil
}

func (e *Evaluator) evalInLoc(expr Expr, loc Path) error {
	switch expr := expr.(type) {
	case Null, String, Number, Bool:
		return e.setVal(loc, expr)
	case Array:
		return e.evalArray(expr, loc)
	case Object:
		return e.evalObject(expr, loc)
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
		fn, err := e.getFn(Path(expr.FnRef))
		if err != nil {
			return err
		}
		result, err := fn.Apply(e, expr.Args)
		if err != nil {
			return err
		}
		return e.setVal(loc, result)
	case Fn:
		return e.setFn(loc, expr)
	default:
		panic("unreachable")
	}
}

func (e *Evaluator) evalObject(obj Object, loc Path) error {
	deps := make(map[string]map[Step]struct{})
	var keys []string
	for key, expr := range obj {
		deps[key] = make(map[Step]struct{})
		makeExprDeps(deps[key], expr, loc)
		keys = append(keys, key)
	}
	for len(keys) > 0 {
		var remains []string
		for _, key := range keys {
			if len(deps[key]) == 0 {
				if err := e.evalInLoc(obj[key], append(loc, StringStep(key))); err != nil {
					return err
				}
				for k := range deps {
					delete(deps[k], StringStep(key))
				}
				delete(deps, key)
			} else {
				remains = append(remains, key)
			}
		}
		if len(remains) == len(keys) {
			return tperr.CircularRefError()
		}
		keys = remains
	}
	return nil
}

func (e *Evaluator) evalArray(arr Array, loc Path) error {
	deps := make(map[int]map[Step]struct{})
	var indices []int
	for i, expr := range arr {
		deps[i] = make(map[Step]struct{})
		makeExprDeps(deps[i], expr, loc)
		indices = append(indices, i)
	}
	for len(indices) > 0 {
		var remains []int
		for _, i := range indices {
			if len(deps[i]) == 0 {
				if err := e.evalInLoc(arr[i], append(loc, NumberStep(i))); err != nil {
					return err
				}
				for k := range deps {
					delete(deps[k], NumberStep(i))
				}
				delete(deps, i)
			} else {
				remains = append(remains, i)
			}
		}
		if len(remains) == len(indices) {
			return tperr.CircularRefError()
		}
		indices = remains
	}
	return nil
}

func makeExprDeps(deps map[Step]struct{}, expr Expr, parent Path) {
	switch expr := expr.(type) {
	case Null, String, Number, Bool:
		return
	case Array:
		for _, ex := range expr {
			makeExprDeps(deps, ex, parent)
		}
	case Object:
		for _, ex := range expr {
			makeExprDeps(deps, ex, parent)
		}
	case ValRef:
		if Path(expr).IsChildOf(parent) {
			deps[expr[len(parent)]] = struct{}{}
		}
	case FnRef:
		if Path(expr).IsChildOf(parent) {
			deps[expr[len(parent)]] = struct{}{}
		}
	case FnCall:
		for _, arg := range expr.Args {
			makeExprDeps(deps, arg, parent)
		}
	case TenpenFn:
		dd := make(map[Step]struct{})
		makeExprDeps(dd, expr.Body, parent)
		for _, arg := range expr.Args {
			delete(dd, StringStep(arg))
		}
		for d := range dd {
			deps[d] = struct{}{}
		}
	}
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
	return nil, tperr.NoRefError()
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
	return nil, tperr.NoRefError()
}
