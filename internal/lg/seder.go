package lg

import (
	"encoding/json"
	"strings"

	"github.com/nanozuki/tenpen/tperr"
)

func Unmarshal(data []byte) (Expr, error) {
	var jv any
	if err := json.Unmarshal(data, &jv); err != nil {
		return nil, err
	}
	return jsonValueToExpr(jv)
}

func jsonValueToExpr(jv any) (Expr, error) {
	switch jv := jv.(type) {
	case nil:
		return Null{}, nil
	case string:
		if strings.HasPrefix(jv, "#") && !strings.HasPrefix(jv, "##") {
			path, err := ParsePath(jv[1:])
			if err != nil {
				return nil, err
			}
			return ValRef(path), nil
		}
		if strings.HasPrefix(jv, "$") && !strings.HasPrefix(jv, "$$") {
			path, err := ParsePath(jv[1:])
			if err != nil {
				return nil, err
			}
			return FnRef(path), nil
		}
		return String(jv), nil
	case float64:
		return Number(jv), nil
	case bool:
		return Bool(jv), nil
	case []any:
		arr := make(Array, 0, len(jv))
		for _, v := range jv {
			expr, err := jsonValueToExpr(v)
			if err != nil {
				return nil, err
			}
			arr = append(arr, expr)
		}
		if len(arr) > 0 && arr[0].Type() == ExprFnRef {
			if arr[0].(FnRef)[0] == StringStep("def") {
				return parseTenpenFn(arr)
			}
			return parseFnCall(arr)
		}
		return arr, nil
	case map[string]any:
		obj := make(Object, len(jv))
		for k, v := range jv {
			expr, err := jsonValueToExpr(v)
			if err != nil {
				return nil, err
			}
			obj[k] = expr
		}
		return obj, nil
	default:
		panic("unreachable")
	}
}

func parseFnCall(arr Array) (FnCall, error) {
	// arr[0] is name of function, arr[1:] are arguments
	if len(arr) < 2 {
		return FnCall{}, tperr.InvalidFnCallError()
	}
	return FnCall{
		FnRef: arr[0].(FnRef),
		Args:  arr[1:],
	}, nil
}

func parseTenpenFn(arr Array) (TenpenFn, error) {
	// arr[0] is function name "def", arr[1] is string arguments, arr[2] is body
	if len(arr) != 3 || arr[1].Type() != ExprArray {
		return TenpenFn{}, tperr.InvalidFnDefError()
	}
	args := make([]String, 0, len(arr[1].(Array)))
	for _, arg := range arr[1].(Array) {
		if arg.Type() != ExprString {
			return TenpenFn{}, tperr.InvalidFnDefError()
		}
		args = append(args, arg.(String))
	}
	return TenpenFn{
		Args: args,
		Body: arr[2],
	}, nil
}

func Marshal(expr Expr) ([]byte, error) {
	jv := exprToJSONValue(expr)
	return json.Marshal(jv)
}

func exprToJSONValue(expr Expr) any {
	switch expr := expr.(type) {
	case Null:
		return nil
	case String:
		return string(expr)
	case Number:
		return float64(expr)
	case Bool:
		return bool(expr)
	case Array:
		arr := make([]any, 0, len(expr))
		for _, v := range expr {
			arr = append(arr, exprToJSONValue(v))
		}
		return arr
	case Object:
		obj := make(map[string]any, len(expr))
		for k, v := range expr {
			obj[k] = exprToJSONValue(v)
		}
		return obj
	case ValRef:
		return expr.String()
	case FnRef:
		return expr.String()
	case FnCall:
		arr := make([]any, 0, len(expr.Args)+1)
		arr = append(arr, exprToJSONValue(expr.FnRef))
		for _, arg := range expr.Args {
			arr = append(arr, exprToJSONValue(arg))
		}
		return arr
	case TenpenFn:
		args := make([]string, 0, len(expr.Args))
		for _, arg := range expr.Args {
			args = append(args, string(arg))
		}
		return []any{
			"#def",
			args,
			exprToJSONValue(expr.Body),
		}
	case GoFn:
		return "<GoFn>"
	default:
		panic("unreachable")
	}
}
