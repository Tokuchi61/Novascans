package middleware

import (
	"fmt"
	"log/slog"
	"net/http"

	platformhttp "github.com/Tokuchi61/Novascans/internal/platform/http"
)

func Recover(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if recovered := recover(); recovered != nil {
					logger.Error(
						"panic recovered",
						"request_id", GetRequestID(r.Context()),
						"method", r.Method,
						"path", r.URL.Path,
						"panic", fmt.Sprintf("%v", recovered),
					)

					platformhttp.WriteError(w, platformhttp.Internal("internal server error", fmt.Errorf("%v", recovered)))
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
