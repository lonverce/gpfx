package service

import "reflect"

type AbstractProvider interface {
	MustGet(srvType reflect.Type) any
}

type Provider interface {
	AbstractProvider
	CreateScope() LifetimeScope
}

type InterimProvider interface {
	AbstractProvider
	GetOwner() Provider
}

type LifetimeScope interface {
	GetProvider() Provider
	Close()
}

func MustLoad[T any](provider AbstractProvider, ptr *T) {
	*ptr = provider.MustGet(reflect.TypeOf(ptr).Elem()).(T)
}

func MustGet[T any](provider AbstractProvider) T {
	return provider.MustGet(Typeof[T]()).(T)
}
