package account

import "github.com/go-chi/chi/v5"

func (m *Module) RegisterRoutes(r chi.Router) {
	r.Route("/account", func(r chi.Router) {
		r.With(m.auth.Authenticate).Get("/me", m.handler.Me)
		r.With(m.auth.Authenticate).Get("/profile", m.handler.GetOwnProfile)
		r.With(m.auth.Authenticate).Patch("/profile", m.handler.UpdateProfile)
		r.With(m.auth.Authenticate).Get("/settings", m.handler.GetSettings)
		r.With(m.auth.Authenticate).Patch("/settings", m.handler.UpdateSettings)
		r.With(m.auth.Authenticate).Get("/privacy", m.handler.GetPrivacy)
		r.With(m.auth.Authenticate).Patch("/privacy", m.handler.UpdatePrivacy)
		r.With(m.auth.Resolve).Get("/profile/{username}", m.handler.GetPublicProfile)
	})
}
