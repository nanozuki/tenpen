package ast

type ExprType int

const (
	ExprTypeEmpty ExprType = iota
	ExprTypeValue          // Pure value expression, evaluation is return itself
	ExprTypeArray
	ExprTypeObject
	ExprTypeRef
	ExprTypeFunCall
)

type ExprEnv map[string]Expr

type Expr interface {
	ExprType() ExprType
	Eval(runtime Object) (Value, error)
}

// ScalarExpr is Null, String, Number, Bool
func (n Null) ExprType() ExprType                 { return ExprTypeValue }
func (n Null) Eval(runtime Object) (Value, error) { return n, nil }

func (s String) ExprType() ExprType                 { return ExprTypeValue }
func (s String) Eval(runtime Object) (Value, error) { return s, nil }

func (n Number) ExprType() ExprType                 { return ExprTypeValue }
func (n Number) Eval(runtime Object) (Value, error) { return n, nil }

func (b Bool) ExprType() ExprType                 { return ExprTypeValue }
func (b Bool) Eval(runtime Object) (Value, error) { return b, nil }

func (r Ref) ExprType() ExprType                 { return ExprTypeRef }
func (r Ref) Eval(runtime Object) (Value, error) { panic("not implemented") }

func (f Fn) ExprType() ExprType                 { return ExprTypeFunCall }
func (f Fn) Eval(runtime Object) (Value, error) { panic("not implemented") }

func (a Array) ExprType() ExprType                 { return ExprTypeValue }
func (a Array) Eval(runtime Object) (Value, error) { return a, nil }

func (o Object) ExprType() ExprType                 { return ExprTypeValue }
func (o Object) Eval(runtime Object) (Value, error) { return o, nil }

// ArrayExpr is Array of Expr
type ArrayExpr []Expr

func (a ArrayExpr) ExprType() ExprType { return ExprTypeArray }

func (a ArrayExpr) Eval(runtime Object) (Value, error) {
	var values []Value
	for _, expr := range a {
		value, err := expr.Eval(runtime)
		if err != nil {
			return nil, err
		}
		values = append(values, value)
	}
	return Array(values), nil
}

type ObjectExpr map[string]Expr

func (o ObjectExpr) ExprType() ExprType { return ExprTypeObject }

func (o ObjectExpr) Eval(runtime Object) (Value, error) {
	values := make(Object)
	for key, expr := range o {
		value, err := expr.Eval(runtime)
		if err != nil {
			return nil, err
		}
		values[key] = value
	}
	return values, nil
}

type FunCallExpr struct {
	Name Fn
	Args []Expr
}

func (f FunCallExpr) ExprType() ExprType { return ExprTypeFunCall }

func (f FunCallExpr) Eval(runtime Object) (Value, error) {
	panic("not implemented")
}
