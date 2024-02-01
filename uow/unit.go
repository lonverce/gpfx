package uow

import (
	"errors"
	"log"
)

type Unit interface {
	IsTransactional() bool
	LoadValue(key string) UnitValue
	Commit() error
	Rollback() error
	Close()
	AddCommitEventListener(func())
	AddRollbackEventListener(func())
}

type internalUnit struct {
	values                map[string]UnitValue
	valuesBuilder         map[string]UnitValueBuilder
	owner                 *DefaultManager
	parent                *internalUnit
	transactional         bool
	commitEventListener   []func()
	rollbackEventListener []func()
}

func (i *internalUnit) IsTransactional() bool {
	return i.transactional
}

func (i *internalUnit) LoadValue(key string) UnitValue {
	v, exist := i.values[key]
	if exist {
		return v
	}

	builder, exist := i.valuesBuilder[key]
	if !exist {
		panic("undefined key")
	}

	v, err := builder.Build(i.transactional, i.owner.Provider)
	if err != nil {
		panic(err.Error())
	}

	i.values[key] = v
	return v
}

func (i *internalUnit) Commit() error {
	var commitErr error

	for _, value := range i.values {
		commitErr = value.Commit()
		if commitErr != nil {
			return commitErr
		}
	}

	for _, listener := range i.commitEventListener {
		defer i.catchListenerError()
		listener()
	}

	return nil
}

func (i *internalUnit) catchListenerError() {
	if r := recover(); r != nil {
		log.Printf("%v\n", r)
	}
}

func (i *internalUnit) Rollback() error {
	var errList = make([]error, 0)

	for _, value := range i.values {
		err := value.Rollback()
		if err != nil {
			errList = append(errList, err)
		}
	}

	if len(errList) > 0 {
		return errors.Join(errList...)
	}

	for _, listener := range i.rollbackEventListener {
		defer i.catchListenerError()
		listener()
	}
	return nil
}

func (i *internalUnit) Close() {
	defer func() {
		i.owner.OnUnitClosed(i)
	}()

	for _, value := range i.values {
		defer i.catchListenerError()
		value.Close()
	}
}

func (i *internalUnit) AddCommitEventListener(f func()) {
	if f == nil {
		panic("listener is nil")
	}
	i.commitEventListener = append(i.commitEventListener, f)
}

func (i *internalUnit) AddRollbackEventListener(f func()) {
	if f == nil {
		panic("listener is nil")
	}
	i.rollbackEventListener = append(i.rollbackEventListener, f)
}
