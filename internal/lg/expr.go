package lg

import (
	"fmt"
	"strings"
)

type ExprType int

const (
	ExprNull ExprType = iota
	ExprString
	ExprNumber
	ExprBool
	ExprArray
	ExprObject
	ExprValRef
	ExprFnRef
	ExprFnCall
	ExprFn
)

type Expr interface {
	String() string
	Type() ExprType
}

type Null struct{}

func (n Null) String() string {
	return "<null>"
}

func (n Null) Type() ExprType {
	return ExprNull
}

type String string

func (s String) String() string {
	return fmt.Sprintf("string<%s>", string(s))
}

func (s String) Type() ExprType {
	return ExprString
}

type Number float64

func (n Number) String() string {
	return fmt.Sprintf("number<%f>", float64(n))
}

func (n Number) Type() ExprType {
	return ExprNumber
}

type Bool bool

func (b Bool) String() string {
	return fmt.Sprintf("bool<%t>", bool(b))
}

func (b Bool) Type() ExprType {
	return ExprBool
}

type Array []Expr

func (a Array) String() string {
	b := strings.Builder{}
	b.WriteByte('[')
	for i, v := range a {
		if i > 0 {
			b.WriteByte(',')
		}
		b.WriteString(v.String())
	}
	b.WriteByte(']')
	return b.String()
}

func (a Array) Type() ExprType {
	return ExprArray
}

type Object map[string]Expr

func (o Object) String() string {
	b := strings.Builder{}
	b.WriteByte('{')
	i := 0
	for k, v := range o {
		if i > 0 {
			b.WriteByte(',')
		}
		b.WriteByte('"')
		b.WriteString(k)
		b.WriteByte('"')
		b.WriteByte(':')
		b.WriteString(v.String())
		i++
	}
	b.WriteByte('}')
	return b.String()
}

func (o Object) Type() ExprType { return ExprObject }

type ValRef Path

func (v ValRef) Type() ExprType { return ExprValRef }
func (v ValRef) String() string { return "#" + Path(v).String() }

type FnRef Path

func (f FnRef) Type() ExprType { return ExprFnRef }
func (f FnRef) String() string { return "$" + Path(f).String() }

type FnCall struct {
	FnRef FnRef
	Args  []Expr
}

func (f FnCall) String() string {
	b := strings.Builder{}
	b.WriteString("call<")
	b.WriteString(f.FnRef.String())
	b.WriteByte('(')
	for i, arg := range f.Args {
		if i > 0 {
			b.WriteByte(',')
		}
		b.WriteString(arg.String())
	}
	b.WriteByte(')')
	b.WriteByte('>')
	return b.String()
}

func (f FnCall) Type() ExprType {
	return ExprFnCall
}

type Fn interface {
	Expr
	Apply(e Evaller, args []Expr) (Expr, error)
}

type Evaller interface {
	SubEvaller(scopedValue Expr) Evaller
	Eval(expr Expr) (Expr, error)
}

type TenpenFn struct {
	Args []String
	Body Expr
}

func (f TenpenFn) String() string {
	b := strings.Builder{}
	b.WriteString("def<[")
	for i, arg := range f.Args {
		if i > 0 {
			b.WriteByte(',')
		}
		b.WriteString(arg.String())
	}
	b.WriteString("](")
	b.WriteString(f.Body.String())
	b.WriteString(")>")
	return b.String()
}

func (f TenpenFn) Type() ExprType {
	return ExprFn
}

func (f TenpenFn) Apply(e Evaller, args []Expr) (Expr, error) {
	for i := len(args); i < len(f.Args); i++ {
		args = append(args, Null{})
	}
	scope := Object{}
	for i, arg := range f.Args {
		scope[string(arg)] = args[i]
	}
	return e.SubEvaller(scope).Eval(f.Body)
}

type GoFn func(e Evaller, args []Expr) (Expr, error)

func (f GoFn) String() string {
	return "<go-func>"
}

func (f GoFn) Type() ExprType {
	return ExprFn
}

func (f GoFn) Apply(e Evaller, args []Expr) (Expr, error) {
	return f(e, args)
}
