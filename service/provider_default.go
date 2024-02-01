package service

import (
	"fmt"
	"reflect"
	"sync"
)

type defaultProvider struct {
	maintainers                  []maintainer
	registrations                map[reflect.Type][]int
	lifetimeCreatedEventHandlers []LifetimeCreatedEventHandler
	sharedPool                   *sync.Pool
	interimPool                  *sync.Pool
	isNew                        bool
}

func (r *defaultProvider) MustGet(srvType reflect.Type) any {
	target := r.getMaintainerIdByServiceType(srvType)
	return r.loadServiceByMaintainerId(target)
}

func clearInterimProviderAndPut(interimProvider *interimProvider, pool *sync.Pool) {
	interimProvider.owner = nil
	interimProvider.instanceList = interimProvider.instanceList[:0]
	pool.Put(interimProvider)
}

func (r *defaultProvider) loadServiceByMaintainer(maintainer maintainer) any {
	interim := r.interimPool.Get().(*interimProvider)
	defer clearInterimProviderAndPut(interim, r.interimPool)

	interim.owner = r
	res := interim.loadServiceByMaintainer(maintainer)
	return res
}

func (r *defaultProvider) loadServiceByMaintainerId(id int) any {
	return r.loadServiceByMaintainer(r.maintainers[id])
}

func (r *defaultProvider) getMaintainerIdByServiceType(serviceType reflect.Type) int {
	if maintainers, ok := r.registrations[serviceType]; ok {
		return maintainers[0]
	} else {
		panic(fmt.Sprintf("service not found: %s", serviceType.String()))
	}
}

func (r *defaultProvider) CreateScope() LifetimeScope {
	child := r.sharedPool.Get().(*defaultProvider)

	if child.isNew == true {
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

func (r *defaultProvider) GetProvider() Provider {
	return r
}

func (r *defaultProvider) Close() {
	for _, maintainer := range r.maintainers {
		maintainer.Clear()
	}
	r.sharedPool.Put(r)
}
