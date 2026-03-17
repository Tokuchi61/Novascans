package account

import (
	accesshttp "github.com/Tokuchi61/Novascans/internal/modules/identity/access/http"
	accountapp "github.com/Tokuchi61/Novascans/internal/modules/user/account/app"
	accounthttp "github.com/Tokuchi61/Novascans/internal/modules/user/account/http"
	"github.com/Tokuchi61/Novascans/internal/modules/user/account/store"
	moduleshared "github.com/Tokuchi61/Novascans/internal/platform/module"
)

type Module struct {
	service *accountapp.Service
	handler *accounthttp.Handler
	guard   *accesshttp.Guard
}

func NewModule(deps moduleshared.Dependencies, guard *accesshttp.Guard) *Module {
	var (
		repo accountapp.Repository
		uow  accountapp.UnitOfWork
	)

	if deps.DB != nil {
		postgresRepository := store.NewPostgresRepository(deps.DB)
		repo = postgresRepository
		uow = store.NewPostgresUnitOfWork(postgresRepository, deps.TxManager)
	} else {
		memoryRepository := store.NewMemoryRepository()
		repo = memoryRepository
		uow = store.NewMemoryUnitOfWork(memoryRepository)
	}

	service := accountapp.NewService(repo, uow)
	handler := accounthttp.NewHandler(deps.Logger, deps.Validator, service)

	return &Module{
		service: service,
		handler: handler,
		guard:   guard,
	}
}

func (m *Module) Key() string {
	return "user.account"
}

func (m *Module) Service() *accountapp.Service {
	return m.service
}
