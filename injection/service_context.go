package injection

import (
	"fmt"
	"github.com/lonverce/gpfx"
	"reflect"
	"sync"
)

type defaultServiceContext struct {
	maintainers                  []serviceMaintainer
	registrations                map[reflect.Type][]int
	lifetimeCreatedEventHandlers []gpfx.LifetimeCreatedEventHandler
	sharedPool                   *sync.Pool
}

func (r *defaultServiceContext) LoadService(serviceType reflect.Type) any {
	return r.getRequiredServiceByMaintainer(r.getMaintainerByServiceType(serviceType))
}

func (r *defaultServiceContext) LoadAllServices(serviceType reflect.Type) []any {
	maintainers := r.getAllMaintainersByServiceType(serviceType)
	return r.getAllServicesByMaintainers(maintainers)
}

func (r *defaultServiceContext) getRequiredServiceByMaintainer(maintainer serviceMaintainer) any {
	interimProvider := &interimServiceContext{
		owner: r,
	}
	res := interimProvider.GetRequiredServiceByMaintainer(maintainer)
	interimProvider.BeginInjections()
	return res
}

func (r *defaultServiceContext) getAllServicesByMaintainers(maintainers []serviceMaintainer) []any {
	interimProvider := &interimServiceContext{
		owner:       r,
		instanceMap: make(map[serviceMaintainer]any),
	}

	results := interimProvider.getAllServicesByMaintainers(maintainers)
	interimProvider.BeginInjections()
	return results
}

func (r *defaultServiceContext) getMaintainerByServiceType(serviceType reflect.Type) serviceMaintainer {
	maintainers, ok := r.registrations[serviceType]
	if !ok {
		panic(fmt.Sprintf("未找到服务类型: %s", serviceType.String()))
	}

	return r.maintainers[maintainers[0]]
}
func (r *defaultServiceContext) getAllMaintainersByServiceType(serviceType reflect.Type) []serviceMaintainer {
	regList, ok := r.registrations[serviceType]
	if !ok {
		panic(fmt.Sprintf("未找到服务类型: %s", serviceType.String()))
	}

	cnt := len(regList)
	maintainers := make([]serviceMaintainer, cnt)

	for i, id := range regList {
		// 因为maintainers的顺序是后注册的在前, 所以要倒一下顺序
		maintainers[cnt-1-i] = r.maintainers[id]
	}
	return maintainers
}

func (r *defaultServiceContext) CreateScope() gpfx.LifetimeScope {
	child := r.sharedPool.Get().(*defaultServiceContext)

	if child.sharedPool == nil {
		// created by pool
		child.sharedPool = r.sharedPool

		for i, maintainer := range r.maintainers {
			child.maintainers[i] = maintainer.Fork(child)
		}
	} else {
		for _, maintainer := range child.maintainers {
			maintainer.ReUse()
		}
	}

	for _, handler := range r.lifetimeCreatedEventHandlers {
		handler(r, child)
	}
	return child
}

func (r *defaultServiceContext) GetServiceContext() gpfx.ServiceContext {
	return r
}

func (r *defaultServiceContext) Close() {
	for _, maintainer := range r.maintainers {
		maintainer.Clear()
	}
	r.sharedPool.Put(r)
}
