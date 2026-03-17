package auth

import "github.com/go-chi/chi/v5"

func (m *Module) RegisterRoutes(r chi.Router) {
	r.Route("/auth", func(r chi.Router) {
		r.Get("/ping", m.handler.Ping)
		r.Post("/users", m.handler.CreateUser)
		r.Get("/users/{id}", m.handler.GetUserByID)
		r.Post("/sessions", m.handler.CreateSession)
		r.Delete("/sessions/{id}", m.handler.RevokeSession)
	})
}
