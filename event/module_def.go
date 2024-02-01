package event

import (
	"github.com/lonverce/gpfx"
	"github.com/lonverce/gpfx/service"
	"github.com/lonverce/gpfx/uow"
	"reflect"
)

var ModuleDef = &gpfx.Module{
	Name:     "gpfx.event",
	DependOn: []*gpfx.Module{uow.ModuleDef},
	ConfigureServices: func(ctx *gpfx.ModuleConfigurator) {
		registry := ctx.GetRegistry()
		gpfx.ConfigOptions[Option](ctx, func(option *Option) {
			option.localEventTypeMap = make(map[reflect.Type]reflect.Type)
			option.localEventRegActions = make([]func(registry service.Registry), 0)
		})

		service.AddTransient[DefaultLocalEventBus](registry, service.Typeof[*DefaultLocalEventBus](), service.Typeof[LocalEventBus]())
	},
	PostConfigureServices: func(ctx *gpfx.ModulePostConfigurator) {
		opt := gpfx.GetOption[Option](ctx)
		registry := ctx.GetRegistry()

		for _, regAction := range opt.localEventRegActions {
			regAction(registry)
		}
	},
}
