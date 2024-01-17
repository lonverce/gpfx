package injection

import (
	"github.com/lonverce/gpfx"
	"sync"
)

type singletonServiceMaintainer struct {
	scopedServiceMaintainer

	// maintainer所属生命周期
	scope *defaultServiceContext
}

func newSingletonMaintainer(reg *gpfx.RegistrationItem, scope *defaultServiceContext) serviceMaintainer {
	tp := new(singletonServiceMaintainer)
	tp.constructor = reg.Constructor
	tp.onceConstructor = sync.OnceValue(tp.createOnce)
	tp.injector = reg.Injector
	tp.scope = scope

	return tp
}

func (p *singletonServiceMaintainer) InjectForInstance(v any, provider *interimServiceContext) {
	if p.injected {
		return
	}

	if provider.owner != p.scope {
		// 由于singleton类型的服务对象始终应保持在root scope中创建,
		// 所以, 如果本次创建请求链最初是由子scope发起的,
		// 就需要切割开这条请求调用链, 单独在root scope中重新构建
		p.scope.getRequiredServiceByMaintainer(p)
		return
	}

	p.scopedServiceMaintainer.InjectForInstance(v, provider)
}

func (p *singletonServiceMaintainer) Fork(*defaultServiceContext) serviceMaintainer {
	if !p.injected {
		v := new(proxyServiceMaintainer)
		v.target = p
		return v
	} else {
		v := new(instanceServiceMaintainer)
		v.instance = p.instance
		return v
	}
}

func (p *singletonServiceMaintainer) Clear() {
}

func (p *singletonServiceMaintainer) ReUse() {
}
