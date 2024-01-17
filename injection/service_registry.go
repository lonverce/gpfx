package injection

import (
	"container/list"
	"github.com/lonverce/gpfx"
	"reflect"
	"sync"
)

type defaultServiceRegistry struct {
	items                        *list.List
	serviceTable                 map[reflect.Type]*list.List
	lifetimeCreatedEventHandlers []gpfx.LifetimeCreatedEventHandler
}

// NewRegistry 创建服务注册表
func NewRegistry() gpfx.ServiceRegistry {
	services := new(defaultServiceRegistry)
	services.items = list.New()
	services.serviceTable = make(map[reflect.Type]*list.List)

	return services
}

func (services *defaultServiceRegistry) AddService(item gpfx.RegistrationItem, types ...reflect.Type) {
	if !item.UseInstance {
		if item.Lifetime < 0 || item.Lifetime > gpfx.ExcludedMaxLifetime {
			panic("unknown Lifetime")
		}

		if item.Constructor == nil {
			panic("initializer is nil")
		}

		if len(types) == 0 {
			panic("lack of publish type")
		}
	}

	id := services.items.Len()
	services.items.PushBack(item)

	for _, serviceType := range types {
		regList, ok := services.serviceTable[serviceType]
		if !ok {
			regList = list.New()
		}

		regList.PushFront(id)
		services.serviceTable[serviceType] = regList
	}
}

func (services *defaultServiceRegistry) AddLifetimeCreatedEventHandler(handler gpfx.LifetimeCreatedEventHandler) {
	if handler == nil {
		return
	}
	services.lifetimeCreatedEventHandlers = append(services.lifetimeCreatedEventHandlers, handler)
}

func (services *defaultServiceRegistry) Build() gpfx.LifetimeScope {
	p := new(defaultServiceContext)

	srvLen := services.items.Len()
	p.maintainers = make([]serviceMaintainer, srvLen+1)

	i := 0
	for pElem := services.items.Front(); pElem != nil; pElem = pElem.Next() {
		reg := pElem.Value.(gpfx.RegistrationItem)
		var maintainer serviceMaintainer

		if reg.UseInstance {
			maintainer = &instanceServiceMaintainer{
				instance: reg.Instance,
			}
		} else {
			switch reg.Lifetime {
			case gpfx.Transient:
				maintainer = newTransientMaintainer(&reg)
			case gpfx.Singleton:
				maintainer = newSingletonMaintainer(&reg, p)
			case gpfx.Scoped:
				if reg.ScopedReuseHandler != nil {
					maintainer = newReuseMaintainer(&reg)
				} else {
					maintainer = newScopedMaintainer(&reg)
				}

			default:
				panic("unhandled default lifetime")
			}
		}

		p.maintainers[i] = maintainer
		i++
	}

	selfMaintainerIndex := i

	// 为ServiceProvider自身提交一个注册
	selfMaintainer := new(serviceProviderMaintainer)
	selfMaintainer.instance = p
	p.maintainers[selfMaintainerIndex] = selfMaintainer

	// 处理服务表
	p.registrations = make(map[reflect.Type][]int)

	for serviceType, regIdList := range services.serviceTable {
		staticIdList := make([]int, regIdList.Len())

		j := 0
		for pElem := regIdList.Front(); pElem != nil; pElem = pElem.Next() {
			staticIdList[j] = pElem.Value.(int)
			j++
		}

		p.registrations[serviceType] = staticIdList
	}

	p.registrations[gpfx.Typeof[gpfx.ServiceContext]()] = []int{selfMaintainerIndex}
	//p.registrations[gpfx.Typeof[gpfx.LifetimeScope]()] = []int{selfMaintainerIndex}
	p.sharedPool = new(sync.Pool)

	regs := p.registrations
	mSize := len(p.maintainers)
	h := services.lifetimeCreatedEventHandlers

	p.sharedPool.New = func() any {
		child := new(defaultServiceContext)
		child.registrations = regs
		child.maintainers = make([]serviceMaintainer, mSize)
		child.lifetimeCreatedEventHandlers = h
		return child
	}

	p.lifetimeCreatedEventHandlers = services.lifetimeCreatedEventHandlers
	return p
}

func (services *defaultServiceRegistry) GetOrAddRegistration(t reflect.Type) *list.List {
	registrations, ok := services.serviceTable[t]
	if !ok {
		registrations = list.New()
		services.serviceTable[t] = registrations
	}
	return registrations
}
