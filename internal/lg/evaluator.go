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
		f:    []Expr{builtins},
	}
	e.f = append(e.f, functions...)
	switch rule.Type() {
	case ExprArray:
		e.v = append(e.v, Array{})
		e.f = append(e.f, Array{})
	case ExprObject:
		e.v = append(e.v, Object{})
		e.f = append(e.f, Object{})
	}
	return e
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

func (e *Evaluator) SubEvaller(scopedValue Expr) Evaller {
	return &Evaluator{
		rule: e.rule,
		v:    append(e.v, scopedValue, Null{}),
		f:    append(e.f, Null{}),
	}
}

func (e *Evaluator) Eval(expr Expr) (Expr, error) {
	return e.eval(expr, []Step{})
}

func (e *Evaluator) eval(expr Expr, loc Path) (Expr, error) {
	switch expr := expr.(type) {
	case Null, String, Number, Bool:
		return expr, nil
	case Object:
		return e.evalObject(expr, loc)
	case Array:
		return e.evalArray(expr, loc)
	case ValRef:
		return e.getVal(Path(expr))
	case FnRef:
		// FnRef cannot be evaluated directly
		return nil, tperr.InvalidTypeError()
	case Fn:
		err := e.setFn(loc, expr)
		return Null{}, err
	case FnCall:
		return e.evalFnCall(expr, expr.Args, loc)
	default:
		panic("unreachable")
	}
}

func (e *Evaluator) evalObject(obj Object, loc Path) (Expr, error) {
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
				result, err := e.eval(obj[key], append(loc, StringStep(key)))
				if err != nil {
					return nil, err
				}
				if err := e.setVal(append(loc, StringStep(key)), result); err != nil {
					return nil, err
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
			return nil, tperr.CircularRefError()
		}
		keys = remains
	}
	return e.getVal(loc)
}

func (e *Evaluator) evalArray(arr Array, loc Path) (Expr, error) {
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
				result, err := e.eval(arr[i], append(loc, NumberStep(i)))
				if err != nil {
					return nil, err
				}
				if err := e.setVal(append(loc, NumberStep(i)), result); err != nil {
					return nil, err
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
			return nil, tperr.CircularRefError()
		}
		indices = remains
	}
	return e.getVal(loc)
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

func (e *Evaluator) evalFnCall(fnCall FnCall, args []Expr, loc Path) (Expr, error) {
	fn, err := e.getFn(Path(fnCall.FnRef))
	if err != nil {
		return nil, err
	}
	switch fn := fn.(type) {
	case GoFn:
		var args []Expr
		for i, arg := range fnCall.Args {
			evaluated, err := e.eval(arg, append(loc, NumberStep(i)))
			if err != nil {
				return nil, err
			}
			args = append(args, evaluated)
		}
		return fn(e, args)
	case TenpenFn:
		if len(args) < len(fn.Args) {
			return nil, tperr.InvalidArgError()
		}
		e.v = append(e.v, Object{})
		defer func() {
			e.v = e.v[:len(e.v)-1]
		}()
		for i, argName := range fn.Args {
			if err := e.setVal(append(loc, StringStep(argName)), args[i]); err != nil {
				return nil, err
			}
		}
		return e.eval(fn.Body, loc)
	default:
		return nil, tperr.InvalidTypeError()
	}
}
