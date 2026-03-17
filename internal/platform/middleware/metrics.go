package middleware

import (
	"net/http"
	"time"

	"github.com/Tokuchi61/Novascans/internal/platform/metrics"
)

func Metrics(registry *metrics.Registry) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if registry == nil {
				next.ServeHTTP(w, r)
				return
			}

			startedAt := time.Now()
			writer := &loggingResponseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			next.ServeHTTP(writer, r)

			registry.ObserveRequest(writer.statusCode, time.Since(startedAt))
		})
	}
}
