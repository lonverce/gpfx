package injection

import (
	"github.com/lonverce/gpfx"
	"sync"
)

type scopedServiceMaintainer struct {
	maintainerBase
	// 由于实例仅能构建一次, 这里需要用onceConstructor来保证并发安全
	onceConstructor func() any
	// 由于实例最多仅能注入一次, 这里需要一个互斥对象来保证并发安全
	injectMutex sync.Mutex
	// 是否已经执行完注入. 如果此服务对象不需要注入, 则在创建后立即将此值置位true
	injected bool
	// 已新建的服务实例, 不保证已完成注入, 请优先确保injected为true,然后再使用此实例
	instance any
}

func newScopedMaintainer(reg *gpfx.RegistrationItem) serviceMaintainer {
	tp := new(scopedServiceMaintainer)
	tp.constructor = reg.Constructor
	tp.onceConstructor = sync.OnceValue(tp.createOnce)
	tp.injector = reg.Injector
	return tp
}

func (p *scopedServiceMaintainer) Clear() {
	p.onceConstructor = nil
	p.injected = false
	p.instance = nil
}

func (p *scopedServiceMaintainer) ReUse() {
	p.onceConstructor = sync.OnceValue(p.createOnce)
}

// createOnce 此方法受sync.OnceValue保护, 确保并发安全且仅执行一次
func (p *scopedServiceMaintainer) createOnce() any {
	// 先创建服务实例
	p.instance = p.constructor()

	if p.injector == nil {
		// 由于此方法受到sync.OnceValue保护,
		// 且必然在InjectForInstance前调用,
		// 所以可以直接操作injected
		p.injected = true
	} else {
		p.injected = false
	}
	return p.instance
}

func (p *scopedServiceMaintainer) CreateServiceInstance() (instance any, needInject bool) {
	if p.injected {
		return p.instance, false
	}

	v := p.onceConstructor()
	return v, !p.injected
}

func (p *scopedServiceMaintainer) InjectForInstance(v any, provider *interimServiceContext) {
	if p.injected {
		return
	}

	// 确保并发安全
	p.injectMutex.Lock()
	defer p.injectMutex.Unlock()

	if p.injected {
		return
	}

	// 执行注入
	p.injector(v, provider)
	p.injected = true
}

func (p *scopedServiceMaintainer) Fork(scope *defaultServiceContext) serviceMaintainer {
	child := new(scopedServiceMaintainer)
	child.maintainerBase = p.maintainerBase
	child.onceConstructor = sync.OnceValue(child.createOnce)
	child.injected = false
	child.instance = nil
	return child
}

type reuseServiceMaintainer struct {
	scopedServiceMaintainer
	reuseDestructor gpfx.ServiceDestructor
}

func newReuseMaintainer(reg *gpfx.RegistrationItem) serviceMaintainer {
	tp := new(reuseServiceMaintainer)
	tp.constructor = reg.Constructor
	tp.onceConstructor = sync.OnceValue(tp.createOnce)
	tp.injector = reg.Injector
	tp.reuseDestructor = reg.ScopedReuseHandler
	return tp
}

func (p *reuseServiceMaintainer) Fork(scope *defaultServiceContext) serviceMaintainer {
	child := new(reuseServiceMaintainer)
	child.maintainerBase = p.maintainerBase
	child.onceConstructor = sync.OnceValue(child.createOnce)
	child.injected = false
	child.instance = nil
	child.reuseDestructor = p.reuseDestructor
	return child
}

func (p *reuseServiceMaintainer) Clear() {
	if p.injected {
		p.reuseDestructor(p.instance)
	}
	p.injected = false
}

func (p *reuseServiceMaintainer) ReUse() {
}
