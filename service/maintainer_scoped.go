package service

import (
	"sync"
)

type scopedServiceMaintainer struct {
	*maintainerBase

	createMutex      sync.Mutex
	created          bool
	creationProvider *interimProvider

	// 是否已经执行完注入. 如果此服务对象不需要注入, 则在创建后立即将此值置位true
	injected bool
	// 已新建的服务实例, 不保证已完成注入, 请优先确保injected为true,然后再使用此实例
	instance any
}

func newScopedMaintainer(reg *Registration) *scopedServiceMaintainer {
	tp := &scopedServiceMaintainer{
		maintainerBase: &maintainerBase{
			Constructor: reg.Constructor,
			Injector:    reg.Injector,
		},
		injected: false,
	}
	return tp
}

func (p *scopedServiceMaintainer) Clear() {
	p.createMutex.Lock()
	defer p.createMutex.Unlock()

	p.created = false
	p.injected = false
	p.instance = nil
	p.creationProvider = nil
}

func (p *scopedServiceMaintainer) TryGetCreatedInstance() (any, bool, bool) {
	if p.created {
		return p.instance, true, p.injected
	} else {
		return nil, false, false
	}
}

func (p *scopedServiceMaintainer) CreateServiceInstance(provider *interimProvider) (any, bool) {
	if obj, created, injected := p.TryGetCreatedInstance(); created {
		return obj, !injected
	}

	p.createMutex.Lock()
	defer p.createMutex.Unlock()

	if !p.created {
		p.instance = p.Constructor(provider)
		p.created = true

		if p.Injector == nil {
			p.injected = true
		} else {
			p.injected = false
			p.creationProvider = provider
		}
	}

	return p.instance, !p.injected
}

func (p *scopedServiceMaintainer) InjectForInstance(v any, provider *interimProvider) {
	if _, created, injected := p.TryGetCreatedInstance(); created && injected {
		return
	}

	p.createMutex.Lock()
	defer p.createMutex.Unlock()

	if p.injected {
		return
	}

	if !p.created || p.creationProvider == nil {
		panic("invalid state")
	}

	// 执行注入
	p.Injector(v, p.creationProvider)
	p.creationProvider = nil
	p.injected = true
}

func (p *scopedServiceMaintainer) Fork(scope *defaultProvider) maintainer {
	child := &scopedServiceMaintainer{
		maintainerBase: p.maintainerBase,
	}

	return child
}
