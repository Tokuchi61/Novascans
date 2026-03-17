package app

import (
	accessmodule "github.com/Tokuchi61/Novascans/internal/modules/identity/access"
	authmodule "github.com/Tokuchi61/Novascans/internal/modules/identity/auth"
	accountmodule "github.com/Tokuchi61/Novascans/internal/modules/user/account"
	moduleshared "github.com/Tokuchi61/Novascans/internal/platform/module"
)

func buildModules(deps moduleshared.Dependencies) []moduleshared.Module {
	auth := authmodule.NewModule(deps)
	access := accessmodule.NewModule(deps, auth.Service())
	account := accountmodule.NewModule(deps, access.Guard())
	auth.Service().SetAccountProvisioner(account.Service())

	return []moduleshared.Module{
		auth,
		access,
		account,
	}
}
