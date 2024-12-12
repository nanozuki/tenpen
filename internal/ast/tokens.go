package ast

type TokenType int

const (
	TokenTypeNull TokenType = iota
	TokenTypeString
	TokenTypeNumber
	TokenTypeBool
	TokenTypeRef
	TokenTypeFun
)

type Token interface {
	Type() TokenType
	Value() any
}

type NullToken struct{}

func (n NullToken) Type() TokenType { return TokenTypeNull }
func (n NullToken) Value() any      { return nil }

type StringToken string

func (s StringToken) Type() TokenType { return TokenTypeString }
func (s StringToken) Value() any      { return string(s) }

type NumberToken float64

func (n NumberToken) Type() TokenType { return TokenTypeNumber }
func (n NumberToken) Value() any      { return float64(n) }

type BoolToken bool

func (b BoolToken) Type() TokenType { return TokenTypeBool }
func (b BoolToken) Value() any      { return bool(b) }

type RefToken string

func (r RefToken) Type() TokenType { return TokenTypeRef }
func (r RefToken) Value() any      { return string(r) }

type FuncToken string

func (f FuncToken) Type() TokenType { return TokenTypeFun }
func (f FuncToken) Value() any      { return string(f) }
