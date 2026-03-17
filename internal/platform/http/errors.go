package http

import (
	"errors"
	stdhttp "net/http"
)

type AppError struct {
	Code       string
	Message    string
	StatusCode int
	Details    any
	Err        error
}

func (err *AppError) Error() string {
	if err.Err != nil {
		return err.Err.Error()
	}

	return err.Message
}

func (err *AppError) Unwrap() error {
	return err.Err
}

func ValidationError(fields map[string]string) *AppError {
	return &AppError{
		Code:       "validation_error",
		Message:    "request validation failed",
		StatusCode: stdhttp.StatusBadRequest,
		Details: map[string]any{
			"fields": fields,
		},
	}
}

func BadRequest(message string, inner error) *AppError {
	return &AppError{
		Code:       "bad_request",
		Message:    message,
		StatusCode: stdhttp.StatusBadRequest,
		Err:        inner,
	}
}

func MethodNotAllowed(message string, inner error) *AppError {
	return &AppError{
		Code:       "method_not_allowed",
		Message:    message,
		StatusCode: stdhttp.StatusMethodNotAllowed,
		Err:        inner,
	}
}

func Unauthorized(message string) *AppError {
	return &AppError{
		Code:       "unauthorized",
		Message:    message,
		StatusCode: stdhttp.StatusUnauthorized,
	}
}

func Forbidden(message string) *AppError {
	return &AppError{
		Code:       "forbidden",
		Message:    message,
		StatusCode: stdhttp.StatusForbidden,
	}
}

func NotFound(message string, inner error) *AppError {
	return &AppError{
		Code:       "not_found",
		Message:    message,
		StatusCode: stdhttp.StatusNotFound,
		Err:        inner,
	}
}

func Conflict(message string, inner error) *AppError {
	return &AppError{
		Code:       "conflict",
		Message:    message,
		StatusCode: stdhttp.StatusConflict,
		Err:        inner,
	}
}

func Internal(message string, inner error) *AppError {
	return &AppError{
		Code:       "internal_error",
		Message:    message,
		StatusCode: stdhttp.StatusInternalServerError,
		Err:        inner,
	}
}

func ServiceUnavailable(message string, inner error) *AppError {
	return &AppError{
		Code:       "service_unavailable",
		Message:    message,
		StatusCode: stdhttp.StatusServiceUnavailable,
		Err:        inner,
	}
}

func ResolveError(err error) *AppError {
	if err == nil {
		return Internal(stdhttp.StatusText(stdhttp.StatusInternalServerError), nil)
	}

	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}

	return Internal(stdhttp.StatusText(stdhttp.StatusInternalServerError), err)
}
