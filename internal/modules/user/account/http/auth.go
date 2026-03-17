package http

import (
	"context"
	nethttp "net/http"
	"strings"

	authapp "github.com/Tokuchi61/Novascans/internal/modules/identity/auth/app"
	platformhttp "github.com/Tokuchi61/Novascans/internal/platform/http"
)

type currentUserContextKey struct{}

type AuthMiddleware struct {
	authenticator authapp.Authenticator
}

func NewAuthMiddleware(authenticator authapp.Authenticator) *AuthMiddleware {
	return &AuthMiddleware{authenticator: authenticator}
}

func (middleware *AuthMiddleware) Resolve(next nethttp.Handler) nethttp.Handler {
	return nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		if strings.TrimSpace(r.Header.Get("Authorization")) == "" {
			next.ServeHTTP(w, r)
			return
		}

		token, err := bearerToken(r)
		if err != nil {
			platformhttp.WriteError(w, platformhttp.Unauthorized("invalid authorization header"))
			return
		}

		currentUser, err := middleware.authenticator.AuthenticateAccessToken(r.Context(), token)
		if err != nil {
			writeServiceError(w, err)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), currentUserContextKey{}, currentUser)))
	})
}

func (middleware *AuthMiddleware) Authenticate(next nethttp.Handler) nethttp.Handler {
	return middleware.Resolve(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		if _, ok := CurrentUser(r.Context()); !ok {
			platformhttp.WriteError(w, platformhttp.Unauthorized("authentication required"))
			return
		}

		next.ServeHTTP(w, r)
	}))
}

func CurrentUser(ctx context.Context) (authapp.AuthenticatedUser, bool) {
	currentUser, ok := ctx.Value(currentUserContextKey{}).(authapp.AuthenticatedUser)
	return currentUser, ok
}
