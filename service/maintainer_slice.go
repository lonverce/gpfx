package service

import "reflect"

type sliceMaintainer struct {
	itemType        reflect.Type
	sliceType       reflect.Type
	itemMaintainers []int
}

func (m *sliceMaintainer) CreateServiceInstance(provider *interimProvider) (instance any, needInject bool) {
	itemLen := len(m.itemMaintainers)
	slice := reflect.MakeSlice(m.sliceType, itemLen, itemLen)
	ctx := provider.owner

	for i, maintainerId := range m.itemMaintainers {
		maintainer := ctx.maintainers[maintainerId]
		item := provider.loadServiceByMaintainer(maintainer)
		itemValue := reflect.ValueOf(item)
		slice.Index(i).Set(itemValue)
	}

	return slice.Interface(), false
}

func (m *sliceMaintainer) InjectForInstance(v any, provider *interimProvider) {
	panic("sliceMaintainer should not be inject")
}

func (m *sliceMaintainer) Fork(scope *defaultProvider) maintainer {
	return m
}

func (m *sliceMaintainer) Clear() {
}

func (m *sliceMaintainer) ReUse() {
}
