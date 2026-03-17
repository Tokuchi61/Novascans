package access

import (
	"github.com/go-chi/chi/v5"

	"github.com/Tokuchi61/Novascans/internal/modules/identity/access/domain"
)

func (m *Module) RegisterRoutes(r chi.Router) {
	r.Route("/access", func(r chi.Router) {
		r.With(m.guard.Authenticate).Get("/me", m.handler.Me)
		r.With(m.guard.RequireBaseRoles(domain.BaseRoleAdmin)).Get("/permissions", m.handler.ListPermissions)
		r.With(m.guard.RequireBaseRoles(domain.BaseRoleAdmin)).Get("/sub-roles", m.handler.ListSubRoles)
		r.With(m.guard.RequireBaseRoles(domain.BaseRoleAdmin)).Post("/sub-roles", m.handler.CreateSubRole)
		r.With(m.guard.RequireBaseRoles(domain.BaseRoleAdmin)).Put("/users/{userID}/base-role", m.handler.UpdateBaseRole)
		r.With(m.guard.RequireBaseRoles(domain.BaseRoleAdmin)).Post("/users/{userID}/sub-roles/{subRoleID}", m.handler.AssignSubRole)
		r.With(m.guard.RequireBaseRoles(domain.BaseRoleAdmin)).Delete("/users/{userID}/sub-roles/{subRoleID}", m.handler.RemoveSubRole)
	})
}
