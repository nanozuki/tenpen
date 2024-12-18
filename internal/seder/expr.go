package seder

import (
	"encoding/json"
	"strings"

	"github.com/nanozuki/tenpen/internal/ast"
)

func ExprUnmarshal(s string) (ast.Expr, error) {
	var jv interface{}
	if err := json.Unmarshal([]byte(s), &jv); err != nil {
		return nil, err // TODO: wrap error
	}
	return parseExpr(jv), nil
}

func parseExpr(jv interface{}) ast.Expr {
	switch jv := jv.(type) {
	case nil:
		return ast.Null{}
	case string:
		if strings.HasPrefix(jv, "#") && !strings.HasPrefix(jv, "##") {
			return ast.Ref(jv)
		}
		if strings.HasPrefix(jv, "$") && !strings.HasPrefix(jv, "$$") {
			return ast.Fn(jv)
		}
		return ast.String(jv)
	case float64:
		return ast.Number(jv)
	case bool:
		return ast.Bool(jv)
	case []any:
		isValue := true
		exprs := make([]ast.Expr, 0, len(jv))
		for _, v := range jv {
			expr := parseExpr(v)
			exprs = append(exprs, expr)
			isValue = isValue && expr.ExprType() == ast.ExprTypeValue
		}
		if len(exprs) > 0 && exprs[0].ExprType() == ast.ExprTypeValue {
			value := exprs[0].(ast.Value)
			if value.ValueType() == ast.TypeFn {
				return ast.FunCallExpr{
					Name: value.(ast.Fn),
					Args: exprs[1:],
				}
			}
		}
		if isValue {
			values := make([]ast.Value, 0, len(exprs))
			for _, expr := range exprs {
				values = append(values, expr.(ast.Value))
			}
			return ast.Array(values)
		}
		return ast.ArrayExpr(exprs)
	case map[string]any:
		exprs := make(ast.ObjectExpr, len(jv))
		isValue := true
		for key, v := range jv {
			exprs[key] = parseExpr(v)
			isValue = isValue && exprs[key].ExprType() == ast.ExprTypeValue
		}
		if isValue {
			values := make(ast.Object, len(exprs))
			for key, expr := range exprs {
				values[key] = expr.(ast.Value)
			}
			return values
		}
		return exprs
	default:
		panic("unexpected value type")
	}
}
