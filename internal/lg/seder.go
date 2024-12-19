package lg

import (
	"encoding/json"
	"strconv"
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
			return parseRef(jv)
		}
		if strings.HasPrefix(jv, "$") && !strings.HasPrefix(jv, "$$") {
			return parseFnName(jv)
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
		if len(arr) > 0 && arr[0].Type() == ExprFn {
			if arr[0].(Fn)[0] == "def" {
				return parseFnDef(arr)
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

func parseRef(s string) (Ref, error) {
	if len(s) < 2 {
		return nil, tperr.InvalidRefError()
	}
	stepStrs := strings.Split(s[1:], ".")
	steps := make([]Step, 0, len(stepStrs))
	for _, s := range stepStrs {
		switch {
		case s == "":
			return nil, tperr.InvalidRefError()
		case s[0] >= '0' && s[0] <= '9':
			n, err := strconv.Atoi(s)
			if err != nil {
				return nil, tperr.InvalidRefError()
			}
			steps = append(steps, NumberStep(n))
		default:
			steps = append(steps, StringStep(s))
		}
	}
	return steps, nil
}

func parseFnName(s string) (Fn, error) {
	if len(s) < 2 {
		return nil, tperr.InvalidFnNameError()
	}
	steps := strings.Split(s[1:], ".")
	for _, s := range steps {
		if s == "" {
			return nil, tperr.InvalidFnNameError()
		}
	}
	return Fn(steps), nil
}

func parseFnCall(arr Array) (FnCall, error) {
	// arr[0] is name of function, arr[1:] are arguments
	if len(arr) < 2 {
		return FnCall{}, tperr.InvalidFnCallError()
	}
	return FnCall{
		Fn:   arr[0].(Fn),
		Args: arr[1:],
	}, nil
}

func parseFnDef(arr Array) (FnDef, error) {
	// arr[0] is function name "def", arr[1] is string arguments, arr[2] is body
	if len(arr) != 3 || arr[1].Type() != ExprArray {
		return FnDef{}, tperr.InvalidFnDefError()
	}
	args := make([]String, 0, len(arr[1].(Array)))
	for _, arg := range arr[1].(Array) {
		if arg.Type() != ExprString {
			return FnDef{}, tperr.InvalidFnDefError()
		}
		args = append(args, arg.(String))
	}
	return FnDef{
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
	case Ref:
		return expr.String()
	case Fn:
		return expr.String()
	case FnCall:
		arr := make([]any, 0, len(expr.Args)+1)
		arr = append(arr, exprToJSONValue(expr.Fn))
		for _, arg := range expr.Args {
			arr = append(arr, exprToJSONValue(arg))
		}
		return arr
	case FnDef:
		args := make([]string, 0, len(expr.Args))
		for _, arg := range expr.Args {
			args = append(args, string(arg))
		}
		return []any{
			"#def",
			args,
			exprToJSONValue(expr.Body),
		}
	default:
		panic("unreachable")
	}
}
