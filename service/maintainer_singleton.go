package service

type singletonServiceMaintainer struct {
	*scopedServiceMaintainer

	// maintainer所属生命周期
	scope *defaultProvider
}

func newSingletonMaintainer(reg *Registration, scope *defaultProvider) *singletonServiceMaintainer {
	tp := &singletonServiceMaintainer{
		scopedServiceMaintainer: newScopedMaintainer(reg),
		scope:                   scope,
	}
	return tp
}

func (p *singletonServiceMaintainer) CreateServiceInstance(provider *interimProvider) (any, bool) {

	if provider.owner != p.scope {
		if obj, ok := p.TryGetInstance(); ok {
			return obj, false
		}

		return p.scope.loadServiceByMaintainer(p), false
	}

	return p.scopedServiceMaintainer.CreateServiceInstance(provider)
}

func (p *singletonServiceMaintainer) Fork(*defaultProvider) maintainer {
	if obj, ok := p.TryGetInstance(); ok {
		v := new(instanceServiceMaintainer)
		v.instance = obj
		return v
	}

	v := new(proxyServiceMaintainer)
	v.target = p
	return v
}

func (p *singletonServiceMaintainer) TryGetInstance() (any, bool) {
	obj, created, injected := p.TryGetCreatedInstance()
	if created && injected {
		return obj, true
	}

	return nil, false
}
