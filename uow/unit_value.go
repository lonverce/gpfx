package uow

import (
	"github.com/lonverce/gpfx/config"
	"github.com/lonverce/gpfx/service"
)

type UnitValue interface {
	Commit() error
	Rollback() error
	Close()
}

type UnitValueBuilder interface {
	Build(transactional bool, provider service.Provider) (UnitValue, error)
}

type Option struct {
	builderMap map[string]UnitValueBuilder
}

func (o *Option) AddValueBuilder(key string, builder UnitValueBuilder) {
	_, exist := o.builderMap[key]
	if exist {
		panic("存在已定义的UnitValueBuilder, key=" + key)
	}

	if builder == nil {
		panic("builder is nil")
	}

	o.builderMap[key] = builder
}

type DefaultValueBuilderManager struct {
	Option config.Option[Option] `gpfx.inject:""`
}

func (mgr *DefaultValueBuilderManager) GetMap() map[string]UnitValueBuilder {
	return mgr.Option.OnceValue().builderMap
}
