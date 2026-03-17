package access

import (
	accessapp "github.com/Tokuchi61/Novascans/internal/modules/identity/access/app"
	accesshttp "github.com/Tokuchi61/Novascans/internal/modules/identity/access/http"
	"github.com/Tokuchi61/Novascans/internal/modules/identity/access/store"
	authapp "github.com/Tokuchi61/Novascans/internal/modules/identity/auth/app"
	moduleshared "github.com/Tokuchi61/Novascans/internal/platform/module"
)

type Module struct {
	service *accessapp.Service
	guard   *accesshttp.Guard
	handler *accesshttp.Handler
}

func NewModule(deps moduleshared.Dependencies, authenticator authapp.Authenticator) *Module {
	var (
		repo accessapp.Repository
		uow  accessapp.UnitOfWork
	)

	if deps.DB != nil {
		repo = store.NewPostgresRepository(deps.DB)
		uow = store.NewPostgresUnitOfWork(deps.DB, deps.TxManager)
	} else {
		memoryRepository := store.NewMemoryRepository()
		repo = memoryRepository
		uow = store.NewMemoryUnitOfWork(memoryRepository)
	}

	service := accessapp.NewService(repo, uow)
	guard := accesshttp.NewGuard(authenticator, service)
	handler := accesshttp.NewHandler(deps.Logger, deps.Validator, service)

	return &Module{
		service: service,
		guard:   guard,
		handler: handler,
	}
}

func (m *Module) Key() string {
	return "identity.access"
}

func (m *Module) Service() *accessapp.Service {
	return m.service
}

func (m *Module) Guard() *accesshttp.Guard {
	return m.guard
}
