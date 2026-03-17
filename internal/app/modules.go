package app

import (
	authmodule "github.com/Tokuchi61/Novascans/internal/modules/identity/auth"
	moduleshared "github.com/Tokuchi61/Novascans/internal/platform/module"
)

func buildModules(deps moduleshared.Dependencies) []moduleshared.Module {
	return []moduleshared.Module{
		authmodule.NewModule(deps),
	}
}
