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
	interimPool                  *sync.Pool
	isNew                        bool
}

func (r *defaultServiceContext) LoadService(serviceType reflect.Type) any {
	return r.getRequiredServiceByMaintainer(r.getMaintainerByServiceType(serviceType))
}

func (r *defaultServiceContext) LoadAllServices(serviceType reflect.Type) []any {
	maintainers := r.getAllMaintainersByServiceType(serviceType)
	return r.getAllServicesByMaintainers(maintainers)
}

func clearInterimProviderAndPut(interimProvider *interimServiceContext, pool *sync.Pool) {
	clear(interimProvider.instanceMap)
	interimProvider.owner = nil
	pool.Put(interimProvider)
}

func (r *defaultServiceContext) getRequiredServiceByMaintainer(maintainer serviceMaintainer) any {
	interimProvider := r.interimPool.Get().(*interimServiceContext)
	defer clearInterimProviderAndPut(interimProvider, r.interimPool)

	interimProvider.owner = r
	res := interimProvider.GetRequiredServiceByMaintainer(maintainer)
	return res
}

func (r *defaultServiceContext) getAllServicesByMaintainers(maintainers []serviceMaintainer) []any {
	interimProvider := r.interimPool.Get().(*interimServiceContext)
	defer clearInterimProviderAndPut(interimProvider, r.interimPool)

	interimProvider.owner = r
	results := interimProvider.getAllServicesByMaintainers(maintainers)
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

	if child.isNew == true {
		for i, maintainer := range r.maintainers {
			child.maintainers[i] = maintainer.Fork(child)
		}
		child.isNew = false
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
