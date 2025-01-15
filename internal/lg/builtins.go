package lg

import (
	"fmt"

	"github.com/nanozuki/tenpen/tperr"
)

var builtins = Object{
	"+": GoFn(add),
	"-": GoFn(sub),
	"*": GoFn(mul),
	"/": GoFn(div),
}

func add(e Evaller, args []Expr) (Expr, error) {
	sum := 0.0
	for _, arg := range args {
		n, ok := arg.(Number)
		fmt.Println(arg, n, ok)
		if !ok {
			return nil, tperr.InvalidTypeError()
		}
		sum += float64(n)
	}
	return Number(sum), nil
}

func sub(e Evaller, args []Expr) (Expr, error) {
	sum := 0.0
	for i, arg := range args {
		n, ok := arg.(Number)
		if !ok {
			return nil, tperr.InvalidTypeError()
		}
		if i == 0 {
			sum = float64(n)
		} else {
			sum -= float64(n)
		}
	}
	return Number(sum), nil
}

func mul(e Evaller, args []Expr) (Expr, error) {
	sum := 1.0
	for _, arg := range args {
		n, ok := arg.(Number)
		if !ok {
			return nil, tperr.InvalidTypeError()
		}
		sum *= float64(n)
	}
	return Number(sum), nil
}

func div(e Evaller, args []Expr) (Expr, error) {
	sum := 1.0
	for i, arg := range args {
		n, ok := arg.(Number)
		if !ok {
			return nil, tperr.InvalidTypeError()
		}
		if i == 0 {
			sum = float64(n)
		} else {
			sum /= float64(n)
		}
	}
	return Number(sum), nil
}
