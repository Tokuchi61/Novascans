package middleware

import (
	"net/http"

	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

func RealIP(next http.Handler) http.Handler {
	return chimiddleware.RealIP(next)
}
