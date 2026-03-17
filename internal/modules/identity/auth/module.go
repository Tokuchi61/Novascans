package auth

import (
	moduleshared "github.com/Tokuchi61/Novascans/internal/platform/module"

	"github.com/Tokuchi61/Novascans/internal/modules/identity/auth/store"
)

type Module struct {
	handler *Handler
}

func NewModule(deps moduleshared.Dependencies) *Module {
	var repo store.Repository
	if deps.DB != nil {
		repo = store.NewPostgresRepository(deps.DB)
	} else {
		repo = store.NewMemoryRepository()
	}

	service := NewService(repo, deps.TxManager, deps.Events)
	handler := NewHandler(deps.Logger, deps.Validator, service)

	return &Module{
		handler: handler,
	}
}

func (m *Module) Key() string {
	return "identity.auth"
}
