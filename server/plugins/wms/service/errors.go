package service

import (
	"errors"
	"fmt"
	"strings"
)

type ErrorKind string

const (
	ErrorKindUnknown   ErrorKind = "unknown"
	ErrorKindInvalid   ErrorKind = "invalid"
	ErrorKindNotFound  ErrorKind = "not_found"
	ErrorKindConflict  ErrorKind = "conflict"
	ErrorKindForbidden ErrorKind = "forbidden"
)

type BizError struct {
	Kind  ErrorKind
	Msg   string
	Cause error
}

func (e *BizError) Error() string {
	if e == nil {
		return ""
	}
	if e.Msg != "" {
		return e.Msg
	}
	if e.Cause != nil {
		return e.Cause.Error()
	}
	return string(e.Kind)
}

func (e *BizError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Cause
}

func InvalidError(msg string) error {
	return &BizError{Kind: ErrorKindInvalid, Msg: msg}
}

func NotFoundError(msg string) error {
	return &BizError{Kind: ErrorKindNotFound, Msg: msg}
}

func ConflictError(msg string) error {
	return &BizError{Kind: ErrorKindConflict, Msg: msg}
}

func ForbiddenError(msg string) error {
	return &BizError{Kind: ErrorKindForbidden, Msg: msg}
}

func WrapDBError(action string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", action, err)
}

func ErrorKindOf(err error) ErrorKind {
	var bizErr *BizError
	if errors.As(err, &bizErr) && bizErr != nil {
		return bizErr.Kind
	}
	return ErrorKindUnknown
}

func IsUniqueConstraintError(err error) bool {
	if err == nil {
		return false
	}
	raw := strings.ToLower(err.Error())
	patterns := []string{
		"error 1062",
		"duplicate entry",
		"duplicated key",
		"unique constraint failed",
		"violates unique constraint",
	}
	for _, p := range patterns {
		if strings.Contains(raw, p) {
			return true
		}
	}
	return false
}
