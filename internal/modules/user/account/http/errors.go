package http

import (
	"errors"
	nethttp "net/http"

	accountapp "github.com/Tokuchi61/Novascans/internal/modules/user/account/app"
	platformhttp "github.com/Tokuchi61/Novascans/internal/platform/http"
)

func writeServiceError(w nethttp.ResponseWriter, err error) {
	var appErr *accountapp.Error
	if !errors.As(err, &appErr) {
		platformhttp.WriteError(w, err)
		return
	}

	switch appErr.Code {
	case accountapp.CodeBadRequest:
		platformhttp.WriteError(w, platformhttp.BadRequest(appErr.Message, appErr.Err))
	case accountapp.CodeConflict:
		platformhttp.WriteError(w, platformhttp.Conflict(appErr.Message, appErr.Err))
	case accountapp.CodeForbidden:
		platformhttp.WriteError(w, platformhttp.Forbidden(appErr.Message))
	case accountapp.CodeUnauthorized:
		platformhttp.WriteError(w, platformhttp.Unauthorized(appErr.Message))
	case accountapp.CodeNotFound:
		platformhttp.WriteError(w, platformhttp.NotFound(appErr.Message, appErr.Err))
	default:
		platformhttp.WriteError(w, platformhttp.Internal(appErr.Message, appErr.Err))
	}
}
