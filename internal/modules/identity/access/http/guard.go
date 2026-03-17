package http

import (
	"context"
	"errors"
	nethttp "net/http"
	"strings"

	accessapp "github.com/Tokuchi61/Novascans/internal/modules/identity/access/app"
	accessdomain "github.com/Tokuchi61/Novascans/internal/modules/identity/access/domain"
	authapp "github.com/Tokuchi61/Novascans/internal/modules/identity/auth/app"
	platformhttp "github.com/Tokuchi61/Novascans/internal/platform/http"
)

type principalContextKey struct{}

type Guard struct {
	authenticator authapp.Authenticator
	service       *accessapp.Service
}

func NewGuard(authenticator authapp.Authenticator, service *accessapp.Service) *Guard {
	return &Guard{
		authenticator: authenticator,
		service:       service,
	}
}

func (guard *Guard) Resolve(next nethttp.Handler) nethttp.Handler {
	return nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		if strings.TrimSpace(r.Header.Get("Authorization")) == "" {
			guest := guard.service.GuestPrincipal()
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), principalContextKey{}, guest)))
			return
		}

		token, err := bearerToken(r)
		if err != nil {
			var appErr *platformhttp.AppError
			if errors.As(err, &appErr) {
				platformhttp.WriteError(w, appErr)
				return
			}

			platformhttp.WriteError(w, platformhttp.Unauthorized("invalid authorization header"))
			return
		}

		principal, err := guard.service.AuthenticatePrincipal(r.Context(), guard.authenticator, token)
		if err != nil {
			writeServiceError(w, err)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), principalContextKey{}, principal)))
	})
}

func (guard *Guard) Authenticate(next nethttp.Handler) nethttp.Handler {
	return guard.Resolve(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		principal, ok := CurrentPrincipal(r.Context())
		if !ok || !principal.IsAuthenticated() {
			platformhttp.WriteError(w, platformhttp.Unauthorized("authentication required"))
			return
		}

		next.ServeHTTP(w, r)
	}))
}

func (guard *Guard) RequireBaseRoles(roles ...string) func(nethttp.Handler) nethttp.Handler {
	return func(next nethttp.Handler) nethttp.Handler {
		return guard.Resolve(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
			principal, ok := CurrentPrincipal(r.Context())
			if !ok || !principal.HasAnyBaseRole(roles...) {
				platformhttp.WriteError(w, platformhttp.Forbidden("forbidden"))
				return
			}

			next.ServeHTTP(w, r)
		}))
	}
}

func CurrentPrincipal(ctx context.Context) (accessdomain.Principal, bool) {
	principal, ok := ctx.Value(principalContextKey{}).(accessdomain.Principal)
	return principal, ok
}
