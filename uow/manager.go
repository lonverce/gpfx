package uow

import (
	"github.com/lonverce/gpfx"
	"github.com/lonverce/gpfx/service"
)

var uowStorageKey = gpfx.NewScopeStorageKey()

type Manager interface {
	GetCurrentUnit() (Unit, bool)
	Begin(transactional bool) Unit
}

type DefaultManager struct {
	OptionMgr *DefaultValueBuilderManager `gpfx.inject:""`
	Storage   gpfx.ScopeStorage           `gpfx.inject:""`
	Provider  service.Provider            `gpfx.inject:""`
}

func (uowMgr *DefaultManager) GetCurrentUnit() (Unit, bool) {
	val, ok := uowMgr.Storage.Get(uowStorageKey)
	if !ok {
		return nil, false
	}
	unit := val.(*internalUnit)
	return unit, unit != nil
}

func (uowMgr *DefaultManager) Begin(transactional bool) Unit {
	data := uowMgr.OptionMgr.GetMap()
	unit := &internalUnit{
		transactional:         transactional,
		values:                make(map[string]UnitValue),
		valuesBuilder:         data,
		owner:                 uowMgr,
		commitEventListener:   make([]func(), 0),
		rollbackEventListener: make([]func(), 0),
	}

	val, ok := uowMgr.Storage.Get(uowStorageKey)
	if ok {
		unit.parent = val.(*internalUnit)
	}
	uowMgr.Storage.Set(uowStorageKey, unit)
	return unit
}

func (uowMgr *DefaultManager) OnUnitClosed(unit *internalUnit) {
	val, ok := uowMgr.Storage.Get(uowStorageKey)
	if !ok {
		panic("The closing unit must be same as current unit")
	}

	currentUnit := val.(*internalUnit)
	if currentUnit != unit {
		panic("The closing unit must be same as current unit")
	}
	if unit.parent != nil {
		uowMgr.Storage.Set(uowStorageKey, unit.parent)
		unit.parent = nil
	} else {
		uowMgr.Storage.Delete(uowStorageKey)
	}
}
