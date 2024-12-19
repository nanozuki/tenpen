package tenpen

import "encoding/json"

func Eval(rule string, envs ...string) (string, error) {
	// TODO: Implement
	return "", nil
}

type Result struct {
	Ret string
	Err error
}

func (r Result) Unmarshal(dest interface{}) error {
	if r.Err != nil {
		return r.Err
	}
	return json.Unmarshal([]byte(r.Ret), dest)
}

func EvalResult(rule string, envs ...string) Result {
	ret, err := Eval(rule, envs...)
	return Result{Ret: ret, Err: err}
}
