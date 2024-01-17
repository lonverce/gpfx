package gpfx

import (
	"reflect"
)

// ModuleConfigurator 模块配置器
type ModuleConfigurator struct {
	services       ServiceRegistry
	configuration  Configuration
	optionBuilders map[reflect.Type]optionBuilder
}

func (ctx *ModuleConfigurator) GetConfig() Configuration {
	return ctx.configuration
}

func (ctx *ModuleConfigurator) GetRegistry() ServiceRegistry {
	return ctx.services
}

func ConfigOptions[TOption any](ctx *ModuleConfigurator, configAction func(option *TOption)) {
	if configAction == nil {
		panic("configAction == nil")
	}
	optionType := Typeof[TOption]()
	regList, ok := ctx.optionBuilders[optionType]

	if !ok {
		regList = &typedOptionBuilder[TOption]{
			actions: make([]func(option *TOption), 0),
		}
		ctx.optionBuilders[optionType] = regList
	}

	b := regList.(*typedOptionBuilder[TOption])
	b.actions = append(b.actions, configAction)
}
