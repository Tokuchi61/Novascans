package auth

import "github.com/go-chi/chi/v5"

func (m *Module) RegisterRoutes(r chi.Router) {
	r.Route("/auth", func(r chi.Router) {
		r.Get("/ping", m.handler.Ping)
		r.Post("/register", m.handler.Register)
		r.Post("/login", m.handler.Login)
		r.Post("/refresh", m.handler.Refresh)
		r.Get("/me", m.handler.Me)
		r.Post("/logout", m.handler.LogoutCurrentSession)
		r.Post("/logout-all", m.handler.LogoutAllSessions)
		r.Post("/email/verify-request", m.handler.RequestEmailVerification)
		r.Post("/email/verify", m.handler.VerifyEmail)
		r.Post("/password/forgot", m.handler.ForgotPassword)
		r.Post("/password/reset", m.handler.ResetPassword)
	})
}
