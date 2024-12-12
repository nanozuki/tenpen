package ast

type ValueType int

const (
	ValueTypeNull ValueType = iota
	ValueTypeString
	ValueTypeNumber
	ValueTypeBool
	ValueTypeArray
	ValueTypeObject
)

type Value interface {
	ValueType() ValueType
	Value() any
}

type NullValue struct{}

func (n NullValue) ValueType() ValueType { return ValueTypeNull }
func (n NullValue) Value() any           { return nil }

type StringValue string

func (s StringValue) ValueType() ValueType { return ValueTypeString }
func (s StringValue) Value() any           { return string(s) }

type NumberValue float64

func (n NumberValue) ValueType() ValueType { return ValueTypeNumber }
func (n NumberValue) Value() any           { return float64(n) }

type BoolValue bool

func (b BoolValue) ValueType() ValueType { return ValueTypeBool }
func (b BoolValue) Value() any           { return bool(b) }

type ArrayValue []Value

func (a ArrayValue) ValueType() ValueType { return ValueTypeArray }
func (a ArrayValue) Value() any           { return []Value(a) }

type ObjectValue map[string]Value

func (o ObjectValue) ValueType() ValueType { return ValueTypeObject }
func (o ObjectValue) Value() any           { return map[string]Value(o) }
