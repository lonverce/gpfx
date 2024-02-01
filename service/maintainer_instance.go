package service

type instanceServiceMaintainer struct {
	instance any
}

func (p *instanceServiceMaintainer) CreateServiceInstance(*interimProvider) (instance any, needInject bool) {
	return p.instance, false
}

func (p *instanceServiceMaintainer) InjectForInstance(_ any, _ *interimProvider) {
	panic("instanceServiceMaintainer should never be called for inject")
}

func (p *instanceServiceMaintainer) Fork(*defaultProvider) maintainer {
	child := new(instanceServiceMaintainer)
	child.instance = p.instance
	return child
}

func (p *instanceServiceMaintainer) Clear() {
}

func (p *instanceServiceMaintainer) ReUse() {
}
