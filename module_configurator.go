package gpfx

import (
	"github.com/lonverce/gpfx/config"
	"github.com/lonverce/gpfx/service"
	"reflect"
)

// ModuleConfigurator 模块配置器
type ModuleConfigurator struct {
	services       service.Registry
	configuration  config.Provider
	optionBuilders map[reflect.Type]optionBuilder
}

func (ctx *ModuleConfigurator) GetConfig() config.Provider {
	return ctx.configuration
}

func (ctx *ModuleConfigurator) GetRegistry() service.Registry {
	return ctx.services
}

func ConfigOptions[TOption any](ctx *ModuleConfigurator, configAction func(option *TOption)) {
	if configAction == nil {
		panic("configAction == nil")
	}
	optionType := service.Typeof[TOption]()
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

type ModulePostConfigurator struct {
	services      service.Registry
	configuration config.Provider
	options       map[reflect.Type]any
}

func (ctx *ModulePostConfigurator) GetConfig() config.Provider {
	return ctx.configuration
}

func (ctx *ModulePostConfigurator) GetRegistry() service.Registry {
	return ctx.services
}

func GetOption[TOption any](ctx *ModulePostConfigurator) *TOption {
	t := service.Typeof[config.Option[TOption]]()
	o, exist := ctx.options[t].(config.Option[TOption])

	if !exist {
		panic("option not found")
	}
	return o.Value()
}
