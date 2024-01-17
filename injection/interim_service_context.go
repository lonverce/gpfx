package injection

import (
	"github.com/lonverce/gpfx"
	"reflect"
)

type pendingItem struct {
	maintainer serviceMaintainer
	instance   any
}

type interimServiceContext struct {
	owner             *defaultServiceContext
	instanceMap       map[serviceMaintainer]any
	pendingInjectList []*pendingItem
}

func (p *interimServiceContext) GetOwner() gpfx.ServiceContext {
	return p.owner
}

func (p *interimServiceContext) LoadService(serviceType reflect.Type) any {
	return p.GetRequiredServiceByMaintainer(p.owner.getMaintainerByServiceType(serviceType))
}

func (p *interimServiceContext) LoadAllServices(serviceType reflect.Type) []any {
	maintainers := p.owner.getAllMaintainersByServiceType(serviceType)
	return p.getAllServicesByMaintainers(maintainers)
}

func (p *interimServiceContext) findCacheInstance(maintainer serviceMaintainer) (any, bool) {
	if p.instanceMap == nil {
		return nil, false
	}

	v, ok := p.instanceMap[maintainer]
	return v, ok
}

func (p *interimServiceContext) getAllServicesByMaintainers(maintainers []serviceMaintainer) []any {
	results := make([]any, len(maintainers))

	for i, maintainer := range maintainers {
		results[i] = p.GetRequiredServiceByMaintainer(maintainer)
	}

	return results
}

func (p *interimServiceContext) GetRequiredServiceByMaintainer(maintainer serviceMaintainer) any {
	cacheInstance, ok := p.findCacheInstance(maintainer)
	if ok {
		return cacheInstance
	}

	newInstance, needInject := maintainer.CreateServiceInstance()

	if needInject {
		if p.pendingInjectList == nil {
			p.pendingInjectList = make([]*pendingItem, 0, 8)
			if p.instanceMap == nil {
				p.instanceMap = make(map[serviceMaintainer]any)
			}
		}

		p.pendingInjectList = append(p.pendingInjectList, &pendingItem{
			maintainer: maintainer,
			instance:   newInstance,
		})
	}

	if p.instanceMap != nil {
		p.instanceMap[maintainer] = newInstance
	}

	return newInstance
}

func (p *interimServiceContext) BeginInjections() {

	if p.pendingInjectList == nil {
		return
	}

	for len(p.pendingInjectList) > 0 {
		item := p.pendingInjectList[0]
		p.pendingInjectList = p.pendingInjectList[1:]
		item.maintainer.InjectForInstance(item.instance, p)
	}

	clear(p.instanceMap)
	clear(p.pendingInjectList)
}
