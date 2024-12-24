package lg

import (
	"strconv"
	"strings"

	"github.com/nanozuki/tenpen/tperr"
)

type Step interface {
	StepType() StepType
	String() string
}
type StepType int

const (
	StepTypeString StepType = iota
	StepTypeNumber
)

type StringStep string

func (s StringStep) StepType() StepType { return StepTypeString }
func (s StringStep) String() string     { return string(s) }

type NumberStep int

func (n NumberStep) StepType() StepType { return StepTypeNumber }
func (n NumberStep) String() string     { return strconv.Itoa(int(n)) }

type Path []Step // Path is a list of steps, for example: a.b.c.0.d

func ParsePath(s string) (Path, error) {
	stepStrs := strings.Split(s, ".")
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

func (r Path) String() string {
	var b strings.Builder
	for i, s := range r {
		if i > 0 {
			b.WriteRune('.')
		}
		b.WriteString(s.String())
	}
	return b.String()
}

func (r Path) GetFrom(target Expr) (Expr, error) {
	switch {
	case target.Type() == ExprObject && len(r) > 0 && r[0].StepType() == StepTypeString:
		obj := target.(Object)
		key := string(r[0].(StringStep))
		if _, ok := obj[key]; !ok {
			return Null{}, nil
		}
		if len(r) > 1 {
			return r[1:].GetFrom(obj[key])
		}
		return obj[key], nil
	case target.Type() == ExprArray && len(r) > 0 && r[0].StepType() == StepTypeNumber:
		arr := target.(Array)
		idx := int(r[0].(NumberStep))
		if idx < 0 || idx >= len(arr) {
			return Null{}, nil
		}
		if len(r) > 1 {
			return r[1:].GetFrom(arr[idx])
		}
		return arr[idx], nil
	default:
		return nil, tperr.NoRefError()
	}
}

func (r Path) SetTo(target Expr, value Expr) error {
	switch {
	case target.Type() == ExprObject && len(r) > 0 && r[0].StepType() == StepTypeString:
		obj := target.(Object)
		key := string(r[0].(StringStep))
		if len(r) > 1 {
			if v, ok := obj[key]; !ok || v.Type() == ExprNull {
				if r[1].StepType() == StepTypeString {
					obj[key] = Object{}
				} else {
					obj[key] = Array{}
				}
			}
			return r[1:].SetTo(obj[key], value)
		}
		obj[key] = value
		return nil
	case target.Type() == ExprArray && len(r) > 0 && r[0].StepType() == StepTypeNumber:
		arr := target.(Array)
		idx := int(r[0].(NumberStep))
		if idx < 0 {
			return tperr.NoRefError()
		}
		for i := len(arr); i <= idx; i++ {
			arr = append(arr, Null{})
		}
		if len(r) > 1 {
			if arr[idx].Type() == ExprNull {
				if r[1].StepType() == StepTypeString {
					arr[idx] = Object{}
				} else {
					arr[idx] = Array{}
				}
			}
			return r[1:].SetTo(arr[idx], value)
		}
		arr[idx] = value
		return nil
	default:
		return tperr.NoRefError()
	}
}

func (p Path) IsChildOf(other Path) bool {
	if len(p) <= len(other) {
		return false
	}
	for i, s := range other {
		if p[i] != s {
			return false
		}
	}
	return true
}
