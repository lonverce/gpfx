package event

import (
	"errors"
	"fmt"
	"github.com/lonverce/gpfx/service"
	"log"
	"reflect"
)

type AbstractLocalEventListenerManager interface {
	SupportEventType() reflect.Type
	DispatchEvent(eventData any)
}

type LocalEventListenerManager[T any] struct {
	ListenerBuilders []func(provider service.Provider) LocalEventListener[T] `gpfx.inject:""`
	RootProvider     service.Provider                                        `gpfx.inject:""`
}

func (m *LocalEventListenerManager[T]) SupportEventType() reflect.Type {
	return service.Typeof[T]()
}

func (m *LocalEventListenerManager[T]) DispatchEvent(eventData any) {
	cnt := len(m.ListenerBuilders)
	ch := make(chan error, cnt)
	typedEvent := eventData.(*T)

	for i := 0; i < cnt; i++ {
		executeScope := m.RootProvider.CreateScope()
		go m.BeginHandle(i, typedEvent, executeScope, ch)
	}

	for i := 0; i < cnt; i++ {
		err := <-ch
		if err != nil {
			log.Println(err.Error())
		}
	}
}

func (m *LocalEventListenerManager[T]) BeginHandle(handlerIndex int, eventData *T, scope service.LifetimeScope, ch chan error) {
	defer func() {
		scope.Close()
		err := recover()
		if err != nil {
			ch <- errors.New(fmt.Sprintf("%v", err))
		} else {
			ch <- nil
		}
	}()

	listener := m.ListenerBuilders[handlerIndex](scope.GetProvider())
	listener.HandleLocalEvent(eventData)
}
