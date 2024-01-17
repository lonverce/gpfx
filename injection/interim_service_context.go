package injection

import (
	"github.com/lonverce/gpfx"
	"reflect"
)

type interimServiceContext struct {
	owner       *defaultServiceContext
	instanceMap map[serviceMaintainer]any
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

func (p *interimServiceContext) getAllServicesByMaintainers(maintainers []serviceMaintainer) []any {
	results := make([]any, len(maintainers))

	for i, maintainer := range maintainers {
		results[i] = p.GetRequiredServiceByMaintainer(maintainer)
	}

	return results
}

func (p *interimServiceContext) GetRequiredServiceByMaintainer(maintainer serviceMaintainer) any {
	cacheInstance, ok := p.instanceMap[maintainer]
	if ok {
		return cacheInstance
	}

	newInstance, needInject := maintainer.CreateServiceInstance()
	p.instanceMap[maintainer] = newInstance

	if needInject {
		maintainer.InjectForInstance(newInstance, p)
	}

	return newInstance
}
