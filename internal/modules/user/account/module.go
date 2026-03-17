package account

import (
	authapp "github.com/Tokuchi61/Novascans/internal/modules/identity/auth/app"
	accountapp "github.com/Tokuchi61/Novascans/internal/modules/user/account/app"
	accounthttp "github.com/Tokuchi61/Novascans/internal/modules/user/account/http"
	"github.com/Tokuchi61/Novascans/internal/modules/user/account/store"
	moduleshared "github.com/Tokuchi61/Novascans/internal/platform/module"
)

type Module struct {
	service *accountapp.Service
	handler *accounthttp.Handler
	auth    *accounthttp.AuthMiddleware
}

func NewModule(deps moduleshared.Dependencies, authenticator authapp.Authenticator) *Module {
	var (
		repo accountapp.Repository
		uow  accountapp.UnitOfWork
	)

	if deps.DB != nil {
		repo = store.NewPostgresRepository(deps.DB)
		uow = store.NewPostgresUnitOfWork(deps.DB)
	} else {
		memoryRepository := store.NewMemoryRepository()
		repo = memoryRepository
		uow = store.NewMemoryUnitOfWork(memoryRepository)
	}

	service := accountapp.NewService(repo, uow)
	handler := accounthttp.NewHandler(deps.Logger, deps.Validator, service)
	authMiddleware := accounthttp.NewAuthMiddleware(authenticator)

	return &Module{
		service: service,
		handler: handler,
		auth:    authMiddleware,
	}
}

func (m *Module) Key() string {
	return "user.account"
}

func (m *Module) Service() *accountapp.Service {
	return m.service
}
