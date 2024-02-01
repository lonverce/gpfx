package authority

import (
	"github.com/lonverce/gpfx"
	"github.com/lonverce/gpfx/security"
	"github.com/lonverce/gpfx/service"
)

var ModuleDef = &gpfx.Module{
	Name:     "gpfx.authority",
	DependOn: []*gpfx.Module{security.ModuleDef},
	ConfigureServices: func(ctx *gpfx.ModuleConfigurator) {
		registry := ctx.GetRegistry()

		gpfx.ConfigOptions[Option](ctx, func(option *Option) {
			option.items = make(map[string]*policyDeclare)
		})

		service.AddSingleton[DefaultPolicyManager](registry, service.Typeof[PolicyManager]())
		service.AddTransient[DefaultService](registry, service.Typeof[Service](), service.Typeof[*DefaultService]())
	},
}
