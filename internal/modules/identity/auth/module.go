package auth

import (
	authapp "github.com/Tokuchi61/Novascans/internal/modules/identity/auth/app"
	authhttp "github.com/Tokuchi61/Novascans/internal/modules/identity/auth/http"
	moduleshared "github.com/Tokuchi61/Novascans/internal/platform/module"

	"github.com/Tokuchi61/Novascans/internal/modules/identity/auth/store"
)

type Module struct {
	service *authapp.Service
	handler *authhttp.Handler
}

func NewModule(deps moduleshared.Dependencies) *Module {
	var (
		repo authapp.Repository
		uow  authapp.UnitOfWork
	)
	if deps.DB != nil {
		repo = store.NewPostgresRepository(deps.DB)
		uow = store.NewPostgresUnitOfWork(deps.DB)
	} else {
		memoryRepository := store.NewMemoryRepository()
		repo = memoryRepository
		uow = store.NewMemoryUnitOfWork(memoryRepository)
	}

	service := authapp.NewService(repo, uow, deps.Events, authapp.ServiceConfig{
		AppEnv:                    deps.Config.App.Env,
		AccessTokenSecret:         deps.Config.Auth.AccessTokenSecret,
		AccessTokenTTL:            deps.Config.Auth.AccessTokenTTL,
		RefreshTokenTTL:           deps.Config.Auth.RefreshTokenTTL,
		EmailVerificationTokenTTL: deps.Config.Auth.EmailVerificationTokenTTL,
		PasswordResetTokenTTL:     deps.Config.Auth.PasswordResetTokenTTL,
	}, nil)
	handler := authhttp.NewHandler(deps.Logger, deps.Validator, service)

	return &Module{
		service: service,
		handler: handler,
	}
}

func (m *Module) Key() string {
	return "identity.auth"
}

func (m *Module) Service() *authapp.Service {
	return m.service
}
