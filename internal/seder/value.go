package seder

import (
	"encoding/json"

	"github.com/nanozuki/tenpen/internal/ast"
)

func ValueUnmarshal(s string) (ast.Value, error) {
	var jv interface{}
	if err := json.Unmarshal([]byte(s), &jv); err != nil {
		return nil, err // TODO: wrap error
	}
	return parseValue(jv), nil
}

func parseValue(jv interface{}) ast.Value {
	switch jv := jv.(type) {
	case nil:
		return ast.NullValue{}
	case string:
		return ast.StringValue(jv)
	case float64:
		return ast.NumberValue(jv)
	case bool:
		return ast.BoolValue(jv)
	case []interface{}:
		values := make([]ast.Value, 0, len(jv))
		for _, v := range jv {
			values = append(values, parseValue(v))
		}
		return ast.ArrayValue(values)
	case map[string]interface{}:
		values := make(ast.ObjectValue, len(jv))
		for key, v := range jv {
			values[key] = parseValue(v)
		}
		return values
	default:
		panic("unexpected value type")
	}
}

func ValueMarshal(v ast.Value) (string, error) {
	b, err := json.Marshal(v)
	return string(b), err
}
