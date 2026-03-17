package app

import "errors"

type ErrorCode string

const (
	CodeBadRequest   ErrorCode = "bad_request"
	CodeConflict     ErrorCode = "conflict"
	CodeForbidden    ErrorCode = "forbidden"
	CodeUnauthorized ErrorCode = "unauthorized"
	CodeNotFound     ErrorCode = "not_found"
	CodeInternal     ErrorCode = "internal"
)

type Error struct {
	Code    ErrorCode
	Message string
	Err     error
}

func (err *Error) Error() string {
	if err.Err != nil {
		return err.Err.Error()
	}

	return err.Message
}

func (err *Error) Unwrap() error {
	return err.Err
}

func BadRequest(message string, inner error) *Error {
	return &Error{Code: CodeBadRequest, Message: message, Err: inner}
}

func Conflict(message string, inner error) *Error {
	return &Error{Code: CodeConflict, Message: message, Err: inner}
}

func Forbidden(message string, inner error) *Error {
	return &Error{Code: CodeForbidden, Message: message, Err: inner}
}

func Unauthorized(message string, inner error) *Error {
	return &Error{Code: CodeUnauthorized, Message: message, Err: inner}
}

func NotFound(message string, inner error) *Error {
	return &Error{Code: CodeNotFound, Message: message, Err: inner}
}

func Internal(message string, inner error) *Error {
	return &Error{Code: CodeInternal, Message: message, Err: inner}
}

func HasCode(err error, code ErrorCode) bool {
	var target *Error
	if !errors.As(err, &target) {
		return false
	}

	return target.Code == code
}
