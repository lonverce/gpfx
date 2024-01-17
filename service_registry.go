package gpfx

import (
	"reflect"
)

// ServiceLifetime 服务生命周期
type ServiceLifetime int

const (
	Transient ServiceLifetime = iota
	Singleton
	Scoped
	ExcludedMaxLifetime
)

type ServiceConstructor func() any
type ServiceInjector func(instance any, provider InterimServiceContext)
type ServiceDestructor func(v any)

// RegistrationItem 服务注册项
type RegistrationItem struct {
	Constructor ServiceConstructor
	Injector    ServiceInjector
	Lifetime    ServiceLifetime
	Instance    any
	UseInstance bool

	// 仅在Scope类型注册时有效
	ScopedReuseHandler ServiceDestructor
}

type LifetimeCreatedEventHandler func(parent, child ServiceContext)

// ServiceRegistry 服务注册表
type ServiceRegistry interface {
	AddService(item RegistrationItem, types ...reflect.Type)
	AddLifetimeCreatedEventHandler(handler LifetimeCreatedEventHandler)
	Build() LifetimeScope
}

type ISupportScopeReuse interface {
	HandleClearBeforeReuse()
}
