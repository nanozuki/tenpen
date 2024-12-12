package seder

import (
	"encoding/json"

	"github.com/nanozuki/tenpen/internal/ast"
)

func ExprUnmarshal(s string) (ast.Expr, error) {
	var jv interface{}
	if err := json.Unmarshal([]byte(s), &jv); err != nil {
		return nil, err // TODO: wrap error
	}
	panic("not implemented")
}

// tokenizeJson converts json value to ast.Token
// the type of return can be:
// - ast.Token
// - []any
// - map[string]any
func tokenizeJson(jv any) any {
	switch jv := jv.(type) {
	case nil:
		return ast.NullToken{}
	case string:
		return ast.StringToken(jv)
	case float64:
		return ast.NumberToken(jv)
	case bool:
		return ast.BoolToken(jv)
	case []interface{}:
		tokens := make([]any, 0, len(jv))
		for _, v := range jv {
			tokens = append(tokens, tokenizeJson(v))
		}
		return tokens
	case map[string]interface{}:
		tokens := make(map[string]any, len(jv))
		for key, v := range jv {
			tokens[key] = tokenizeJson(v)
		}
		return tokens
	default:
		panic("unexpected value type")
	}
}
