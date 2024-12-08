package tenpen

type ExpressionType int

const (
	EmptyExpression ExpressionType = iota
	ValueExpression
	ArrayExpression
	ObjectExpression
	RefExpression
	FunCallExpression
	FunDefExpression
)

type ExpressionEnv map[string]Expression

type Expression interface {
	Type() ExpressionType
	Eval(envs ...ExpressionEnv) (Expression, error)
}
