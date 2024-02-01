package domain

import (
	"github.com/lonverce/gpfx"
	"github.com/lonverce/gpfx/clock"
	"github.com/lonverce/gpfx/ddd/shared"
	"github.com/lonverce/gpfx/service"
)

var ModuleDef = &gpfx.Module{
	Name:     "gpfx.ddd.domain",
	DependOn: []*gpfx.Module{shared.ModuleDef, gpfx.ModuleDef, clock.ModuleDef},
	ConfigureServices: func(ctx *gpfx.ModuleConfigurator) {
		registry := ctx.GetRegistry()
		service.AddTransient[Service](registry, service.Typeof[*Service]())
	},
}
