package http

import (
	"errors"
	"net/http"
	"strings"

	platformhttp "github.com/Tokuchi61/Novascans/internal/platform/http"
)

func bearerToken(r *http.Request) (string, error) {
	value := strings.TrimSpace(r.Header.Get("Authorization"))
	if value == "" {
		return "", platformhttp.Unauthorized("missing authorization header")
	}

	parts := strings.SplitN(value, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", errors.New("invalid authorization header")
	}

	token := strings.TrimSpace(parts[1])
	if token == "" {
		return "", errors.New("missing bearer token")
	}

	return token, nil
}
