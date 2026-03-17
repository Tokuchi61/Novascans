package http

import (
	"errors"
	nethttp "net/http"

	authapp "github.com/Tokuchi61/Novascans/internal/modules/identity/auth/app"
	platformhttp "github.com/Tokuchi61/Novascans/internal/platform/http"
)

func writeServiceError(w nethttp.ResponseWriter, err error) {
	var appErr *authapp.Error
	if !errors.As(err, &appErr) {
		platformhttp.WriteError(w, err)
		return
	}

	switch appErr.Code {
	case authapp.CodeBadRequest:
		platformhttp.WriteError(w, platformhttp.BadRequest(appErr.Message, appErr.Err))
	case authapp.CodeConflict:
		platformhttp.WriteError(w, platformhttp.Conflict(appErr.Message, appErr.Err))
	case authapp.CodeForbidden:
		platformhttp.WriteError(w, platformhttp.Forbidden(appErr.Message))
	case authapp.CodeUnauthorized:
		platformhttp.WriteError(w, platformhttp.Unauthorized(appErr.Message))
	case authapp.CodeNotFound:
		platformhttp.WriteError(w, platformhttp.NotFound(appErr.Message, appErr.Err))
	default:
		platformhttp.WriteError(w, platformhttp.Internal(appErr.Message, appErr.Err))
	}
}
