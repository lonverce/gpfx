package gpfx

import (
	"github.com/lonverce/gpfx/service"
	"reflect"
	"sync"
)

type LazyServiceProvider interface {
	service.AbstractProvider
}

type DefaultLazyServiceProvider struct {
	Provider service.Provider `gpfx.inject:""`
	Cache    sync.Map
}

func (d *DefaultLazyServiceProvider) MustGet(srvType reflect.Type) any {
	createFunc, ok := d.Cache.Load(srvType)

	if !ok {
		createFunc = sync.OnceValue(func() any {
			return d.Provider.MustGet(srvType)
		})

		if v, loaded := d.Cache.LoadOrStore(srvType, createFunc); loaded {
			createFunc = v
		}
	}

	return createFunc.(func() any)()
}
