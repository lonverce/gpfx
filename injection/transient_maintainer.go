package injection

import "github.com/lonverce/gpfx"

// transientServiceMaintainer 用于创建临时生命周期的服务
type transientServiceMaintainer struct {
	maintainerBase
}

func newTransientMaintainer(reg *gpfx.RegistrationItem) serviceMaintainer {
	tp := new(transientServiceMaintainer)
	tp.constructor = reg.Constructor
	tp.injector = reg.Injector
	return tp
}

func (p *transientServiceMaintainer) CreateServiceInstance() (instance any, needInject bool) {
	// 先调用构造器创建服务实例
	v := p.constructor()

	if p.injector == nil {
		return v, false
	}

	return v, true
}

func (p *transientServiceMaintainer) InjectForInstance(v any, provider *interimServiceContext) {
	// 对服务实例执行注入操作
	p.injector(v, provider)
}

func (p *transientServiceMaintainer) Fork(*defaultServiceContext) serviceMaintainer {
	tp := new(transientServiceMaintainer)
	tp.maintainerBase = p.maintainerBase
	return tp
}
