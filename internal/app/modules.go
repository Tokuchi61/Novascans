package app

import (
	accessmodule "github.com/Tokuchi61/Novascans/internal/modules/identity/access"
	authmodule "github.com/Tokuchi61/Novascans/internal/modules/identity/auth"
	moduleshared "github.com/Tokuchi61/Novascans/internal/platform/module"
)

func buildModules(deps moduleshared.Dependencies) []moduleshared.Module {
	auth := authmodule.NewModule(deps)
	access := accessmodule.NewModule(deps, auth.Service())

	return []moduleshared.Module{
		auth,
		access,
	}
}
