package clock

import (
	"github.com/lonverce/gpfx"
	"github.com/lonverce/gpfx/service"
)

var ModuleDef = &gpfx.Module{
	Name:     "gpfx.clock",
	DependOn: []*gpfx.Module{gpfx.ModuleDef},
	ConfigureServices: func(ctx *gpfx.ModuleConfigurator) {
		registry := ctx.GetRegistry()

		gpfx.ConfigOptions[Option](ctx, func(option *Option) {
			option.UseUTC = true
		})

		service.AddSingleton[DefaultService](registry, service.Typeof[Service]())
	},
}
