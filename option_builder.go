package gpfx

import (
	"github.com/lonverce/gpfx/config"
	"github.com/lonverce/gpfx/service"
	"reflect"
	"sync"
)

type optionBuilder interface {
	Build() (any, reflect.Type)
}

type typedOptionBuilder[TOption any] struct {
	actions []func(option *TOption)
}

type defaultOption[TOption any] struct {
	actions       []func(option *TOption)
	onceValueFunc func() *TOption
}

func (d *defaultOption[TOption]) OnceValue() *TOption {
	return d.onceValueFunc()
}

func (d *defaultOption[TOption]) Value() *TOption {
	v := new(TOption)
	for _, action := range d.actions {
		action(v)
	}
	return v
}

func (d *typedOptionBuilder[TOption]) Build() (any, reflect.Type) {

	o := &defaultOption[TOption]{
		actions: make([]func(option *TOption), len(d.actions)),
	}

	copy(o.actions, d.actions)

	o.onceValueFunc = sync.OnceValue[*TOption](o.Value)

	return o, service.Typeof[config.Option[TOption]]()
}
