package event

import (
	"github.com/lonverce/gpfx/uow"
	"reflect"
)

// LocalEventBus 本地事件总线
type LocalEventBus interface {
	// Publish 立即发布事件
	Publish(event any, eventType reflect.Type)

	// PublishOnCommit 等到当前工作单元 uow.Unit 提交成功后再发布事件
	PublishOnCommit(event any, eventType reflect.Type)
}

type LocalEventListener[T any] interface {
	HandleLocalEvent(event *T)
}

type DefaultLocalEventBus struct {
	UnitOfWorkManager uow.Manager                         `gpfx.inject:""`
	ListenerManagers  []AbstractLocalEventListenerManager `gpfx.inject:""`
}

func (bus *DefaultLocalEventBus) Publish(event any, eventType reflect.Type) {

	for _, listenerManger := range bus.ListenerManagers {
		if listenerManger.SupportEventType() == eventType {
			listenerManger.DispatchEvent(event)
			return
		}
	}

	panic("事件类型未注册")
}

func (bus *DefaultLocalEventBus) PublishOnCommit(event any, eventType reflect.Type) {
	unit, exist := bus.UnitOfWorkManager.GetCurrentUnit()
	if !exist {
		bus.Publish(event, eventType)
	}

	unit.AddCommitEventListener(func() {
		bus.Publish(event, eventType)
	})
}
