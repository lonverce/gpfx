package service

type proxyServiceMaintainer struct {
	target maintainer
}

func (p *proxyServiceMaintainer) CreateServiceInstance(provider *interimProvider) (instance any, needInject bool) {
	return p.target.CreateServiceInstance(provider)
}

func (p *proxyServiceMaintainer) InjectForInstance(v any, provider *interimProvider) {
	p.target.InjectForInstance(v, provider)
}

func (p *proxyServiceMaintainer) Fork(scope *defaultProvider) maintainer {
	return p.target.Fork(scope)
}

func (p *proxyServiceMaintainer) Clear() {
}

func (p *proxyServiceMaintainer) ReUse() {
	if t, ok := p.target.(*singletonServiceMaintainer); ok {
		if obj, ok := t.TryGetInstance(); ok {
			p.target = &instanceServiceMaintainer{
				instance: obj,
			}
		}
	}
}
