package gpfx

import (
	"reflect"
)

type AbstractServiceContext interface {
	LoadService(serviceType reflect.Type) any
	LoadAllServices(serviceType reflect.Type) []any
}

type ServiceContext interface {
	AbstractServiceContext
	CreateScope() LifetimeScope
}

type InterimServiceContext interface {
	AbstractServiceContext
	GetOwner() ServiceContext
}

type LifetimeScope interface {
	GetServiceContext() ServiceContext
	Close()
}

func LoadService[T any](provider AbstractServiceContext) T {
	return provider.LoadService(Typeof[T]()).(T)
}

func LoadAllServices[T any](provider AbstractServiceContext) []T {
	results := provider.LoadAllServices(Typeof[T]())
	typedResults := make([]T, len(results))

	for i, result := range results {
		typedResults[i] = result.(T)
	}

	return typedResults
}
