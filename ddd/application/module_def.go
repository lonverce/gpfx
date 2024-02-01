package application

import (
	"github.com/lonverce/gpfx"
	"github.com/lonverce/gpfx/authority"
	"github.com/lonverce/gpfx/cache"
	"github.com/lonverce/gpfx/ddd/contract"
	"github.com/lonverce/gpfx/ddd/domain"
	"github.com/lonverce/gpfx/event"
	"github.com/lonverce/gpfx/service"
	"github.com/lonverce/gpfx/uow"
)

var ModuleDef = &gpfx.Module{
	Name: "gpfx.ddd.application",
	DependOn: []*gpfx.Module{contract.ModuleDef, domain.ModuleDef,
		authority.ModuleDef, cache.ModuleDef,
		uow.ModuleDef, event.ModuleDef},
	ConfigureServices: func(ctx *gpfx.ModuleConfigurator) {
		registry := ctx.GetRegistry()
		service.AddTransient[Service](registry, service.Typeof[*Service]())
	},
}
