package uow

import (
	"github.com/lonverce/gpfx"
	"github.com/lonverce/gpfx/service"
)

var ModuleDef = &gpfx.Module{
	Name:     "gpfx.uow",
	DependOn: []*gpfx.Module{gpfx.ModuleDef},
	ConfigureServices: func(ctx *gpfx.ModuleConfigurator) {
		registry := ctx.GetRegistry()

		gpfx.ConfigOptions[gpfx.ScopeStorageOption](ctx, func(option *gpfx.ScopeStorageOption) {
			option.Define(uowStorageKey, service.Typeof[*internalUnit](), nil)
		})

		gpfx.ConfigOptions[Option](ctx, func(option *Option) {
			option.builderMap = make(map[string]UnitValueBuilder)
		})

		service.AddSingleton[DefaultValueBuilderManager](registry)
		service.AddTransient[DefaultManager](registry, service.Typeof[Manager](), service.Typeof[*DefaultManager]())
	},
}
