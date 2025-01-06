package errors

import (
	"fmt"
	"strings"
)

type (
	Error struct {
		code    string
		message string
		wrapErr error
	}
)

func (e Error) Code() string {
	return e.code
}

func (e Error) Message() string {
	return e.message
}

func (e *Error) Error() string {
	return fmt.Sprintf("code:%d message:%s", e.code, e.message)
}

func (e *Error) Is(err error) bool {
	if e == nil && err == nil {
		return true
	}
	switch valErr := err.(type) {
	case *Error:
		if e.code == valErr.code && e.message == valErr.message {
			return true
		}
		return false
	default:
		return false
	}
}

func WithCode(code string, msg ...string) error {
	e := &Error{
		code:    code,
		message: strings.Join(msg, ";"),
	}
	return e
}

func WrapWithCode(err error, code string, msg ...string) error {
	e := &Error{
		code:    code,
		message: strings.Join(msg, ";"),
		wrapErr: err,
	}
	return e
}

func UnWrap(err error) error {
	switch valErr := err.(type) {
	case *Error:
		return valErr.wrapErr
	default:
		return err
	}
}
