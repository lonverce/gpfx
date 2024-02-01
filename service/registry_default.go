package service

import (
	"container/list"
	"reflect"
	"slices"
	"sync"
)

type defaultRegistry struct {
	items                        *list.List
	serviceTable                 map[reflect.Type]*list.List
	lifetimeCreatedEventHandlers []LifetimeCreatedEventHandler
}

func (services *defaultRegistry) AddService(item Registration, types ...reflect.Type) {
	if item.Lifetime < 0 || item.Lifetime >= ExcludedMaxLifetime {
		panic("unknown Lifetime")
	}

	if item.Lifetime != ExternalInstance {
		if item.Constructor == nil {
			panic("initializer is nil")
		}

		if len(types) == 0 {
			panic("缺少公布的服务类型")
		}

		for _, publishType := range types {
			k := publishType.Kind()

			if k == reflect.Interface {
				continue
			}

			if k == reflect.Pointer {
				elemType := publishType.Elem()
				if elemType.Kind() == reflect.Struct {
					continue
				}
			}

			panic("发布的服务类型必须为结构体指针或接口")
		}
	} else {
		if item.Instance == nil {
			panic("当Lifetime为ExternalInstance时，必须设置Instance")
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

func (services *defaultRegistry) AddLifetimeCreatedEventHandler(handler LifetimeCreatedEventHandler) {
	if handler == nil {
		return
	}
	services.lifetimeCreatedEventHandlers = append(services.lifetimeCreatedEventHandlers, handler)
}

func (services *defaultRegistry) Build() LifetimeScope {
	handlers := services.lifetimeCreatedEventHandlers[:]

	p := &defaultProvider{
		/*
		 * 因为要为每个注册的item再注册一个lazy_maintainer, 所以srvLen*2
		 * 然后为每组服务再注册一个slice_maintainer, 所以再+len(services.serviceTable)*2
		 */
		maintainers:                  make([]maintainer, 0, (services.items.Len())*2+len(services.serviceTable)*2),
		registrations:                make(map[reflect.Type][]int),
		sharedPool:                   new(sync.Pool),
		interimPool:                  new(sync.Pool),
		isNew:                        false,
		lifetimeCreatedEventHandlers: handlers,
	}
	p.sharedPool.New = func() any {
		child := &defaultProvider{
			registrations:                p.registrations,
			lifetimeCreatedEventHandlers: p.lifetimeCreatedEventHandlers,
			sharedPool:                   p.sharedPool,
			interimPool:                  p.interimPool,
			maintainers:                  make([]maintainer, len(p.maintainers)),
			isNew:                        true,
		}

		for i, maintainer := range p.maintainers {
			child.maintainers[i] = maintainer.Fork(child)
		}
		return child
	}

	p.interimPool.New = func() any {
		c := new(interimProvider)
		return c
	}

	for pElem := services.items.Front(); pElem != nil; pElem = pElem.Next() {
		reg := pElem.Value.(Registration)
		var maintainer maintainer

		switch reg.Lifetime {
		case ExternalInstance:
			maintainer = &instanceServiceMaintainer{
				instance: reg.Instance,
			}
		case Transient:
			maintainer = newTransientMaintainer(&reg)
		case Singleton:
			maintainer = newSingletonMaintainer(&reg, p)
		case Scoped:
			maintainer = newScopedMaintainer(&reg)
		default:
			panic("注册中含有无法处理的生命周期")
		}

		p.maintainers = append(p.maintainers, maintainer)
	}

	// 处理服务表
	{
		for serviceType, regIdList := range services.serviceTable {
			staticIdList := make([]int, regIdList.Len())

			j := 0
			for pElem := regIdList.Front(); pElem != nil; pElem = pElem.Next() {
				staticIdList[j] = pElem.Value.(int)
				j++
			}

			p.registrations[serviceType] = staticIdList
		}
	}

	// 为所有的现有服务注册都注册多个Lazy服务
	{
		emptyArg := []reflect.Type{Typeof[Provider]()}
		lazyMap := make(map[reflect.Type][]maintainer)

		for itemType, registrations := range p.registrations {
			funcType := reflect.FuncOf(emptyArg, []reflect.Type{itemType}, false)
			maintainerList := make([]maintainer, len(registrations))

			for i, target := range registrations {
				maintainer := &lazyMaintainer{
					funcType:  funcType,
					valueType: itemType,
					target:    target,
				}
				maintainerList[i] = maintainer
			}
			lazyMap[funcType] = maintainerList
		}

		from := len(p.maintainers)
		for funcType, maintainerList := range lazyMap {
			p.maintainers = append(p.maintainers, maintainerList...)
			to := len(p.maintainers)

			idList := make([]int, 0, len(maintainerList))

			for i := from; i < to; i++ {
				idList = append(idList, i)
			}
			p.registrations[funcType] = idList
			from = to
		}
	}

	// 为所有的现有服务注册都注册一组集合服务
	{
		sliceMap := make(map[reflect.Type]*sliceMaintainer)

		for itemType, registrations := range p.registrations {
			reversed := make([]int, len(registrations))
			copy(reversed, registrations)
			slices.Reverse(reversed)

			sliceType := reflect.SliceOf(itemType)
			m := &sliceMaintainer{
				itemType:        itemType,
				sliceType:       sliceType,
				itemMaintainers: reversed,
			}
			sliceMap[sliceType] = m
		}

		for sliceType, maintainer := range sliceMap {
			p.registrations[sliceType] = []int{len(p.maintainers)}
			p.maintainers = append(p.maintainers, maintainer)
		}
	}

	return p
}
