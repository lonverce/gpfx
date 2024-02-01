package service

// transientServiceMaintainer 用于创建临时生命周期的服务
type transientServiceMaintainer struct {
	*maintainerBase
}

func newTransientMaintainer(reg *Registration) *transientServiceMaintainer {
	tp := &transientServiceMaintainer{
		maintainerBase: &maintainerBase{
			Constructor: reg.Constructor,
			Injector:    reg.Injector,
		},
	}
	return tp
}

func (p *transientServiceMaintainer) CreateServiceInstance(provider *interimProvider) (instance any, needInject bool) {
	// 先调用构造器创建服务实例
	v := p.Constructor(provider)

	if p.Injector == nil {
		return v, false
	}

	return v, true
}

func (p *transientServiceMaintainer) InjectForInstance(v any, provider *interimProvider) {
	// 对服务实例执行注入操作
	p.Injector(v, provider)
}

func (p *transientServiceMaintainer) Fork(*defaultProvider) maintainer {
	tp := &transientServiceMaintainer{
		maintainerBase: p.maintainerBase,
	}
	return tp
}
