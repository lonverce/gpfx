package gpfx

import (
	"github.com/lonverce/gpfx/service"
)

var ModuleDef = &Module{
	Name: "gpfx",
	ConfigureServices: func(ctx *ModuleConfigurator) {
		services := ctx.GetRegistry()
		ConfigOptions[ScopeStorageOption](ctx, func(option *ScopeStorageOption) {
			option.defines = make(map[*ScopeStorageKey]*StorageItemOption)
		})
		service.AddScoped[DefaultScopeStorage](services, service.Typeof[*DefaultScopeStorage](), service.Typeof[ScopeStorage]())
		service.AddTransient[DefaultLazyServiceProvider](services, service.Typeof[LazyServiceProvider]())
		services.AddLifetimeCreatedEventHandler(func(parent, child service.Provider) {
			var parentStorage, childStorage *DefaultScopeStorage
			service.MustLoad(parent, &parentStorage)
			service.MustLoad(child, &childStorage)
			childStorage.CopyFrom(parentStorage)
		})
	},
}
