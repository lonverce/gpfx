package gpfx

import (
	"fmt"
	"reflect"
	"sync"
)

type ScopeStorageKey struct {
	id string
}

func NewScopeStorageKey() *ScopeStorageKey {
	key := &ScopeStorageKey{}
	key.id = fmt.Sprintf("%p", key)
	return key
}

type ScopeStorage interface {
	Get(key *ScopeStorageKey) (any, bool)
	Set(key *ScopeStorageKey, newVal any)
	Delete(key *ScopeStorageKey)
}

type storageItemOption struct {
	key      string
	callback func(val any) any
	valType  reflect.Type
}

type ScopeStorageOption struct {
	defines map[*ScopeStorageKey]*storageItemOption
}

func (o *ScopeStorageOption) Define(key *ScopeStorageKey, valType reflect.Type, optionalCallback func(val any) any) {
	opt, exist := o.defines[key]
	if exist {
		panic("redefined storage key")
	}

	opt = &storageItemOption{
		key:      key.id,
		valType:  valType,
		callback: optionalCallback,
	}
	o.defines[key] = opt
}

type internalScopeStorage struct {
	data    sync.Map
	options *ScopeStorageOption
}

func (storage *internalScopeStorage) validateKey(key *ScopeStorageKey) *storageItemOption {
	if key == nil {
		panic("scoped local key is null")
	}
	opt, exist := storage.options.defines[key]
	if !exist {
		panic("undefined key")
	}
	return opt
}

func (storage *internalScopeStorage) Get(key *ScopeStorageKey) (any, bool) {
	storage.validateKey(key)
	return storage.data.Load(key)
}

func (storage *internalScopeStorage) Set(key *ScopeStorageKey, newVal any) {
	opt := storage.validateKey(key)

	if !reflect.TypeOf(newVal).AssignableTo(opt.valType) {
		panic("The given value do not match to pre-defined value type")
	}

	storage.data.Store(key, newVal)
}

func (storage *internalScopeStorage) Delete(key *ScopeStorageKey) {
	storage.data.Delete(key)
}

func (storage *internalScopeStorage) CopyTo(dest *internalScopeStorage) {
	storage.data.Range(func(key, value any) bool {
		opt, ok := storage.options.defines[key.(*ScopeStorageKey)]
		if !ok {
			panic("found undefined key while copying")
		}
		if opt.callback != nil {
			newVal := opt.callback(value)
			if !reflect.TypeOf(newVal).AssignableTo(opt.valType) {
				panic("callback function returns a new value which do not match to pre-defined value type")
			}
			value = newVal
		}
		dest.data.Store(key, value)
		return true
	})
}
