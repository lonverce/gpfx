package gpfx

import (
	"strings"
)

// Module 模块
type Module struct {
	Name                    string
	DependOn                []*Module
	ConfigureServices       func(ctx *ModuleConfigurator)
	PostConfigureServices   func(ctx *ModulePostConfigurator)
	OnApplicationInitialize func(ctx *ModuleInitializer)
	OnApplicationShutdown   func()
}

func (g *Module) GetString() string {
	n := strings.TrimSpace(g.Name)
	if strings.Compare(n, "") == 0 {
		panic("未命名模块")
	}
	return n
}
