package service

import (
	"container/list"
	"reflect"
)

// Lifetime 服务生命周期
type Lifetime int

const (
	Transient Lifetime = iota
	Singleton
	Scoped
	ExternalInstance
	ExcludedMaxLifetime
)

type Constructor func(provider InterimProvider) any
type Injector func(instance any, provider InterimProvider)

// Registration 服务注册项
type Registration struct {
	Constructor Constructor
	Injector    Injector
	Lifetime    Lifetime

	// 当 Lifetime 为 gpfx.ExternalInstance 时，必须设置此项
	Instance any
}

type LifetimeCreatedEventHandler func(parent, child Provider)

// Registry 服务注册表
type Registry interface {
	AddService(item Registration, types ...reflect.Type)
	AddLifetimeCreatedEventHandler(handler LifetimeCreatedEventHandler)
	Build() LifetimeScope
}

// NewRegistry 创建服务注册表
func NewRegistry() Registry {
	services := new(defaultRegistry)
	services.items = list.New()
	services.serviceTable = make(map[reflect.Type]*list.List)
	return services
}
