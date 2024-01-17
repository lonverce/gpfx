package injection

type proxyServiceMaintainer struct {
	target *singletonServiceMaintainer
}

func (p *proxyServiceMaintainer) CreateServiceInstance() (instance any, needInject bool) {
	return p.target.CreateServiceInstance()
}

func (p *proxyServiceMaintainer) InjectForInstance(v any, provider *interimServiceContext) {
	p.target.InjectForInstance(v, provider)
}

func (p *proxyServiceMaintainer) Fork(scope *defaultServiceContext) serviceMaintainer {
	return p.target.Fork(scope)
}

func (p *proxyServiceMaintainer) Clear() {
}

func (p *proxyServiceMaintainer) ReUse() {
}
