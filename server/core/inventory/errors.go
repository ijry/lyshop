package inventory

import (
	"errors"
	"fmt"
)

type ErrorKind string

const (
	ErrorKindInvalid  ErrorKind = "invalid"
	ErrorKindConflict ErrorKind = "conflict"
	ErrorKindNotFound ErrorKind = "not_found"
	ErrorKindInternal ErrorKind = "internal"
)

type Error struct {
	Kind    ErrorKind
	Message string
}

func (e *Error) Error() string {
	return e.Message
}

func newError(kind ErrorKind, msg string) error {
	return &Error{Kind: kind, Message: msg}
}

func Invalid(msg string) error {
	return newError(ErrorKindInvalid, msg)
}

func Conflict(msg string) error {
	return newError(ErrorKindConflict, msg)
}

func NotFound(msg string) error {
	return newError(ErrorKindNotFound, msg)
}

func Internal(msg string, err error) error {
	if err == nil {
		return newError(ErrorKindInternal, msg)
	}
	return newError(ErrorKindInternal, fmt.Sprintf("%s: %v", msg, err))
}

var ErrInventoryBusy = errors.New("inventory busy")
