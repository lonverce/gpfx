package gpfx

var ModuleDef = &Module{
	Name: "gpfx.gpfx",
	ConfigureServices: func(ctx *ModuleConfigurator) {
		services := ctx.GetRegistry()
		ConfigOptions[ScopeStorageOption](ctx, func(option *ScopeStorageOption) {
			option.defines = make(map[*ScopeStorageKey]*storageItemOption)
		})
		AddScoped[internalScopeStorage](services, Typeof[*internalScopeStorage](), Typeof[ScopeStorage]())

		services.AddLifetimeCreatedEventHandler(func(parent, child ServiceContext) {
			parentStorage := LoadService[*internalScopeStorage](parent)
			childStorage := LoadService[*internalScopeStorage](child)
			childStorage.options = parentStorage.options
			childStorage.CopyFrom(parentStorage)
		})
	},

	OnApplicationInitialize: func(ctx *ModuleInitializer) {
		services := ctx.GetProvider()
		rootStorage := LoadService[*internalScopeStorage](services)
		rootStorage.options = LoadService[Option[ScopeStorageOption]](services).Value()
	},
}
