package tenpen

type Rule interface {
	Eval(envs ...string) (string, error)
}

type Runtime map[string]any

type SingleRule struct {
	Expression
}

func (r SingleRule) Eval(envs ...string) (string, error) {
	return r.Expression.Eval(envs)
}

type ArrayRule struct {
	Expressions []Expression
}

type ObjectRule struct {
	Expressions map[string]Expression
}
