package gpfx

import (
	"log"
	"reflect"
)

type HostBuilder struct {
	StartupModule *Module
	Configuration Configuration
	Registry      ServiceRegistry
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

	// 补充注册
	AddInstanceOnly(cfgCtx.services, builder.Configuration)
	AddService[hostedServiceManager](builder.Registry, Transient, Typeof[*hostedServiceManager]())
	for _, optionBuilder := range cfgCtx.optionBuilders {
		optionBuilder.PublishToRegistry(cfgCtx.services)
	}

	rootScope := cfgCtx.GetRegistry().Build()

	return &Host{
		rootScope: rootScope,
		modules:   ctx.solvedModules,
	}
}
