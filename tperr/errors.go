package tperr

import "fmt"

type Error struct {
	Location string
	Message  ErrorMessages
}

func (e *Error) Error() string {
	if e.Location == "" {
		return string(e.Message)
	}
	return fmt.Sprintf("[%s] %s", e.Location, e.Message)
}

func (e *Error) Is(err error) bool {
	if err == nil {
		return false
	}
	if e2, ok := err.(*Error); ok {
		return e.Message == e2.Message
	}
	return false
}

func (e *Error) As(target interface{}) bool {
	if target == nil {
		return false
	}
	if e2, ok := target.(*Error); ok {
		*e2 = *e
		return true
	}
	return false
}

type ErrorMessages string

const (
	InvalidJSON   ErrorMessages = "invalid json"
	InvalidRef    ErrorMessages = "invalid reference"
	InvalidFnName ErrorMessages = "invalid function name"
	InvalidFnCall ErrorMessages = "invalid function call"
	InvalidFnDef  ErrorMessages = "invalid function definition"
	NoRef         ErrorMessages = "no reference"
)

func InvalidJSONError() *Error { // TODO: add location
	return &Error{
		Message: InvalidJSON,
	}
}

func InvalidRefError() *Error { // TODO: add location
	return &Error{
		Message: InvalidRef,
	}
}

func InvalidFnNameError() *Error { // TODO: add location
	return &Error{
		Message: InvalidFnName,
	}
}

func InvalidFnCallError() *Error { // TODO: add location
	return &Error{
		Message: InvalidFnCall,
	}
}

func InvalidFnDefError() *Error { // TODO: add location
	return &Error{
		Message: InvalidFnDef,
	}
}

func NoRefError() *Error { // TODO: add location
	return &Error{
		Message: NoRef,
	}
}
