package access

import "github.com/go-chi/chi/v5"

func (m *Module) RegisterRoutes(r chi.Router) {
	r.Route("/access", func(r chi.Router) {
		r.With(m.guard.Authenticate).Get("/me", m.handler.Me)
		r.With(m.guard.RequireBaseRoles("admin")).Get("/permissions", m.handler.ListPermissions)
		r.With(m.guard.RequireBaseRoles("admin")).Get("/sub-roles", m.handler.ListSubRoles)
		r.With(m.guard.RequireBaseRoles("admin")).Post("/sub-roles", m.handler.CreateSubRole)
		r.With(m.guard.RequireBaseRoles("admin")).Put("/users/{userID}/base-role", m.handler.UpdateBaseRole)
		r.With(m.guard.RequireBaseRoles("admin")).Post("/users/{userID}/sub-roles/{subRoleID}", m.handler.AssignSubRole)
		r.With(m.guard.RequireBaseRoles("admin")).Delete("/users/{userID}/sub-roles/{subRoleID}", m.handler.RemoveSubRole)
	})
}
