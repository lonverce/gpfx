package gpfx

import "github.com/lonverce/gpfx/service"

// ModuleInitializer 模块初始化器
type ModuleInitializer struct {
	ctx service.Provider
}

func (ctx *ModuleInitializer) GetProvider() service.Provider {
	return ctx.ctx
}
