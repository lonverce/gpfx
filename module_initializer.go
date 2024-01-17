package gpfx

// ModuleInitializer 模块初始化器
type ModuleInitializer struct {
	ctx ServiceContext
}

func (ctx *ModuleInitializer) GetProvider() ServiceContext {
	return ctx.ctx
}
