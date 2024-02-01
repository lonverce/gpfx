package gpfx

import (
	"github.com/lonverce/gpfx/config"
	"github.com/lonverce/gpfx/service"
	"log"
	"reflect"
)

type HostBuilder struct {
	StartupModule *Module
	Configuration config.Provider
	Registry      service.Registry
}

func (builder *HostBuilder) Build() *Host {

	ctx := &moduleLoader{
		moduleState:   make(map[*Module]moduleStateEnum),
		solvedModules: make([]*Module, 0),
	}

	log.Println("[gpfx] 开始加载应用依赖模块")
	ctx.Load(builder.StartupModule, "")

	cfgCtx := &ModuleConfigurator{
		services:       builder.Registry,
		configuration:  builder.Configuration,
		optionBuilders: make(map[reflect.Type]optionBuilder),
	}

	log.Println("[gpfx] 正在配置依赖模块")
	for _, module := range ctx.solvedModules {

		if module.ConfigureServices == nil {
			continue
		}
		module.ConfigureServices(cfgCtx)
	}

	postCtx := &ModulePostConfigurator{
		services:      builder.Registry,
		configuration: builder.Configuration,
		options:       make(map[reflect.Type]any),
	}
	{
		for _, optionBuilder := range cfgCtx.optionBuilders {
			o, t := optionBuilder.Build()
			postCtx.options[t] = o
		}

		for _, module := range ctx.solvedModules {

			if module.PostConfigureServices == nil {
				continue
			}
			module.PostConfigureServices(postCtx)
		}
	}

	{
		// 补充注册
		service.AddInstanceOnly(cfgCtx.services, builder.Configuration)
		service.AddTransient[HostedServiceManager](builder.Registry, service.Typeof[*HostedServiceManager]())
		for optionType, optionInstance := range postCtx.options {
			cfgCtx.services.AddService(service.Registration{
				Lifetime: service.ExternalInstance,
				Instance: optionInstance,
			}, optionType)
		}
	}

	rootScope := cfgCtx.GetRegistry().Build()

	return &Host{
		rootScope: rootScope,
		modules:   ctx.solvedModules,
	}
}
