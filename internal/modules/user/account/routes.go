package account

import "github.com/go-chi/chi/v5"

func (m *Module) RegisterRoutes(r chi.Router) {
	r.Route("/account", func(r chi.Router) {
		r.With(m.guard.Authenticate).Get("/me", m.handler.Me)
		r.With(m.guard.Authenticate).Get("/profile", m.handler.GetOwnProfile)
		r.With(m.guard.Authenticate).Patch("/profile", m.handler.UpdateProfile)
		r.With(m.guard.Authenticate).Get("/settings", m.handler.GetSettings)
		r.With(m.guard.Authenticate).Patch("/settings", m.handler.UpdateSettings)
		r.With(m.guard.Authenticate).Get("/privacy", m.handler.GetPrivacy)
		r.With(m.guard.Authenticate).Patch("/privacy", m.handler.UpdatePrivacy)
		r.With(m.guard.Resolve).Get("/profile/{username}", m.handler.GetPublicProfile)
	})
}
