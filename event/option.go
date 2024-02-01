package event

import (
	"github.com/lonverce/gpfx/service"
	"reflect"
)

type Option struct {
	// event type => LocalEventListenerManager type
	localEventTypeMap    map[reflect.Type]reflect.Type
	localEventRegActions []func(registry service.Registry)
}

// RegisterLocalEvent 注册本地事件
func RegisterLocalEvent[TEvent any](option *Option) {
	eventType := service.Typeof[TEvent]()
	_, exist := option.localEventTypeMap[eventType]
	if exist {
		return
	}

	option.localEventTypeMap[eventType] = service.Typeof[LocalEventListenerManager[TEvent]]()
	option.localEventRegActions = append(option.localEventRegActions, AddLocalEventManager[TEvent])
}

func AddLocalEventManager[TEvent any](registry service.Registry) {
	service.AddSingleton[LocalEventListenerManager[TEvent]](registry,
		service.Typeof[*LocalEventListenerManager[TEvent]](),
		service.Typeof[AbstractLocalEventListenerManager]())
}
