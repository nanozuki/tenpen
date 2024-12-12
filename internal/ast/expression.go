package ast

type ExprType int

const (
	ExprTypeEmpty ExprType = iota
	ExprTypeScalar
	ExprTypeArray
	ExprTypeObject
	ExprTypeRef
	ExprTypeFunCall
)

type ExprEnv map[string]Expr

type Expr interface {
	ExprType() ExprType
	Eval(runtime ObjectValue) (Value, error)
}

// ScalarExpr is Null, String, Number, Bool
func (n NullValue) ExprType() ExprType                      { return ExprTypeScalar }
func (n NullValue) Eval(runtime ObjectValue) (Value, error) { return n, nil }

func (s StringValue) ExprType() ExprType                      { return ExprTypeScalar }
func (s StringValue) Eval(runtime ObjectValue) (Value, error) { return s, nil }

func (n NumberValue) ExprType() ExprType                      { return ExprTypeScalar }
func (n NumberValue) Eval(runtime ObjectValue) (Value, error) { return n, nil }

func (b BoolValue) ExprType() ExprType                      { return ExprTypeScalar }
func (b BoolValue) Eval(runtime ObjectValue) (Value, error) { return b, nil }

// ArrayExpr is Array of Expr
type ArrayExpr []Expr

func (a ArrayExpr) ExprType() ExprType { return ExprTypeArray }

func (a ArrayExpr) Eval(runtime ObjectValue) (Value, error) {
	var values []Value
	for _, expr := range a {
		value, err := expr.Eval(runtime)
		if err != nil {
			return nil, err
		}
		values = append(values, value)
	}
	return ArrayValue(values), nil
}

type ObjectExpr map[string]Expr

func (o ObjectExpr) ExprType() ExprType { return ExprTypeObject }

func (o ObjectExpr) Eval(runtime ObjectValue) (Value, error) {
	values := make(ObjectValue)
	for key, expr := range o {
		value, err := expr.Eval(runtime)
		if err != nil {
			return nil, err
		}
		values[key] = value
	}
	return values, nil
}
