package ast

type ValueType int

const (
	TypeNull ValueType = iota
	TypeString
	TypeNumber
	TypeBool
	TypeArray
	TypeObject
	TypeRef
	TypeFn
)

type Value interface {
	ValueType() ValueType
}

type Null struct{}

func (n Null) ValueType() ValueType { return TypeNull }

type String string

func (s String) ValueType() ValueType { return TypeString }

type Number float64

func (n Number) ValueType() ValueType { return TypeNumber }

type Bool bool

func (b Bool) ValueType() ValueType { return TypeBool }

type Array []Value

func (a Array) ValueType() ValueType { return TypeArray }

type Object map[string]Value

func (o Object) ValueType() ValueType { return TypeObject }

type Ref string

func (r Ref) ValueType() ValueType { return TypeRef }

type Fn string

func (f Fn) ValueType() ValueType { return TypeFn }
