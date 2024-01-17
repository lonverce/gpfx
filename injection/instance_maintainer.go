package injection

type instanceServiceMaintainer struct {
	instance any
}

func (p *instanceServiceMaintainer) CreateServiceInstance() (instance any, needInject bool) {
	return p.instance, false
}

func (p *instanceServiceMaintainer) InjectForInstance(_ any, _ *interimServiceContext) {
}

func (p *instanceServiceMaintainer) Fork(*defaultServiceContext) serviceMaintainer {
	child := new(instanceServiceMaintainer)
	child.instance = p.instance
	return child
}

func (p *instanceServiceMaintainer) Clear() {
}

func (p *instanceServiceMaintainer) ReUse() {
}

type serviceProviderMaintainer struct {
	instanceServiceMaintainer
}

func (p *serviceProviderMaintainer) Fork(scope *defaultServiceContext) serviceMaintainer {
	child := new(serviceProviderMaintainer)
	child.instance = scope
	return child
}
