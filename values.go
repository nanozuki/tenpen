package tenpen

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
	Type() ValueType
	Value() any
}

type NullValue struct{}

func (n NullValue) Type() ValueType { return ValueTypeNull }
func (n NullValue) Value() any      { return nil }

type StringValue string

func (s StringValue) Type() ValueType { return ValueTypeString }
func (s StringValue) Value() any      { return string(s) }

type NumberValue float64

func (n NumberValue) Type() ValueType { return ValueTypeNumber }
func (n NumberValue) Value() any      { return float64(n) }

type BoolValue bool

func (b BoolValue) Type() ValueType { return ValueTypeBool }
func (b BoolValue) Value() any      { return bool(b) }

type ArrayValue []Value

func (a ArrayValue) Type() ValueType { return ValueTypeArray }
func (a ArrayValue) Value() any      { return []Value(a) }

type ObjectValue map[string]Value

func (o ObjectValue) Type() ValueType { return ValueTypeObject }
func (o ObjectValue) Value() any      { return map[string]Value(o) }
