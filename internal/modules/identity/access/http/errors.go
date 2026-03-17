package http

import (
	"errors"
	nethttp "net/http"

	accessapp "github.com/Tokuchi61/Novascans/internal/modules/identity/access/app"
	platformhttp "github.com/Tokuchi61/Novascans/internal/platform/http"
)

func writeServiceError(w nethttp.ResponseWriter, err error) {
	var appErr *accessapp.Error
	if !errors.As(err, &appErr) {
		platformhttp.WriteError(w, err)
		return
	}

	switch appErr.Code {
	case accessapp.CodeBadRequest:
		platformhttp.WriteError(w, platformhttp.BadRequest(appErr.Message, appErr.Err))
	case accessapp.CodeConflict:
		platformhttp.WriteError(w, platformhttp.Conflict(appErr.Message, appErr.Err))
	case accessapp.CodeForbidden:
		platformhttp.WriteError(w, platformhttp.Forbidden(appErr.Message))
	case accessapp.CodeUnauthorized:
		platformhttp.WriteError(w, platformhttp.Unauthorized(appErr.Message))
	case accessapp.CodeNotFound:
		platformhttp.WriteError(w, platformhttp.NotFound(appErr.Message, appErr.Err))
	default:
		platformhttp.WriteError(w, platformhttp.Internal(appErr.Message, appErr.Err))
	}
}
