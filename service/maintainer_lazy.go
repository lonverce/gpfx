package service

import (
	"reflect"
)

type lazyMaintainer struct {
	valueType reflect.Type
	funcType  reflect.Type
	target    int
}

func (m *lazyMaintainer) CreateServiceInstance(*interimProvider) (instance any, needInject bool) {
	return reflect.MakeFunc(m.funcType, m.internalLoadServiceInstance).Interface(), false
}

func (m *lazyMaintainer) internalLoadServiceInstance(args []reflect.Value) (results []reflect.Value) {
	ctx := args[0].Interface().(*defaultProvider)

	instance := ctx.loadServiceByMaintainerId(m.target)
	return []reflect.Value{reflect.ValueOf(instance).Convert(m.valueType)}
}

func (m *lazyMaintainer) InjectForInstance(any, *interimProvider) {
	panic("lazyMaintainer should never inject")
}

func (m *lazyMaintainer) Fork(scope *defaultProvider) maintainer {
	return m
}

func (m *lazyMaintainer) Clear() {
}

func (m *lazyMaintainer) ReUse() {
}
