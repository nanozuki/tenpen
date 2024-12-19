package lg

import (
	"strconv"
	"strings"

	"github.com/nanozuki/tenpen/tperr"
)

type ExprType int

const (
	ExprNull ExprType = iota
	ExprString
	ExprNumber
	ExprBool
	ExprArray
	ExprObject
	ExprRef
	ExprFn
	ExprFnCall
	ExprFnDef
)

type Expr interface {
	Type() ExprType
}

type Null struct{}

func (n Null) Type() ExprType {
	return ExprNull
}

type String string

func (s String) Type() ExprType {
	return ExprString
}

type Number float64

func (n Number) Type() ExprType {
	return ExprNumber
}

type Bool bool

func (b Bool) Type() ExprType {
	return ExprBool
}

type Array []Expr

func (a Array) Type() ExprType {
	return ExprArray
}

func (a Array) GetByRef(ref Ref) (Expr, error) {
	switch {
	case len(ref) == 0 || ref[0].StepType() == StepTypeString:
		return nil, tperr.NoRefError()
	case len(ref) == 1:
		index := int(ref[0].(NumberStep))
		if index < 0 || index >= len(a) {
			return Null{}, nil
		}
		return a[index], nil
	default: // len(ref) > 1
		child := a[int(ref[0].(NumberStep))]
		switch child := child.(type) {
		case Array:
			return child.GetByRef(ref[1:])
		case Object:
			return child.GetByRef(ref[1:])
		default:
			return nil, tperr.NoRefError()
		}
	}
}

type Object map[string]Expr

func (o Object) Type() ExprType { return ExprObject }

func (o Object) GetByRef(ref Ref) (Expr, error) {
	switch {
	case len(ref) == 0 || ref[0].StepType() == StepTypeNumber:
		return nil, tperr.NoRefError()
	case len(ref) == 1:
		child, ok := o[string(ref[0].(StringStep))]
		if !ok {
			return Null{}, nil
		}
		return child, nil
	default: // len(ref) > 1
		child := o[string(ref[0].(StringStep))]
		switch child := child.(type) {
		case Array:
			return child.GetByRef(ref[1:])
		case Object:
			return child.GetByRef(ref[1:])
		default:
			return nil, tperr.NoRefError()
		}
	}
}

type Step interface {
	StepType() StepType
}
type StepType int

const (
	StepTypeString StepType = iota
	StepTypeNumber
)

type StringStep string

func (s StringStep) StepType() StepType { return StepTypeString }

type NumberStep int

func (n NumberStep) StepType() StepType { return StepTypeNumber }

type Ref []Step

func (r Ref) Type() ExprType {
	return ExprRef
}

func (r Ref) String() string {
	var b strings.Builder
	b.WriteRune('#')
	for i, s := range r {
		if i > 0 {
			b.WriteRune('.')
		}
		switch s := s.(type) {
		case StringStep:
			b.WriteString(string(s))
		case NumberStep:
			b.WriteString(strconv.Itoa(int(s)))
		}
	}
	return b.String()
}

type Fn []string

func (f Fn) Type() ExprType {
	return ExprFn
}

func (f Fn) String() string {
	var b strings.Builder
	b.WriteRune('$')
	for i, s := range f {
		if i > 0 {
			b.WriteRune('.')
		}
		b.WriteString(s)
	}
	return b.String()
}

type FnCall struct {
	Fn   Fn
	Args []Expr
}

func (f FnCall) Type() ExprType {
	return ExprFnCall
}

type FnDef struct {
	Args []String
	Body Expr
}

func (f FnDef) Type() ExprType {
	return ExprFnDef
}
