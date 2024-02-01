package security

import (
	"github.com/lonverce/gpfx"
	"github.com/lonverce/gpfx/service"
)

var ModuleDef = &gpfx.Module{
	Name:     "gpfx.security",
	DependOn: []*gpfx.Module{gpfx.ModuleDef},
	ConfigureServices: func(ctx *gpfx.ModuleConfigurator) {
		gpfx.ConfigOptions[gpfx.ScopeStorageOption](ctx, func(option *gpfx.ScopeStorageOption) {
			option.Define(identityStorageKey, service.Typeof[*ClaimIdentity](), nil)
		})

		registry := ctx.GetRegistry()
		service.AddTransient[internalIdentityAccessor](registry, service.Typeof[*internalIdentityAccessor](), service.Typeof[IdentityAccessor]())
		service.AddTransient[DefaultCurrentUser](registry, service.Typeof[CurrentUser]())
	},
}
