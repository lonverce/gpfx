package gpfx

import (
	"fmt"
	"github.com/lonverce/gpfx/config"
	"github.com/lonverce/gpfx/service"
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

type StorageCopyCallback func(val any, newCtx service.Provider) any

type StorageItemOption struct {
	key      string
	callback StorageCopyCallback
	valType  reflect.Type
}

type ScopeStorageOption struct {
	defines map[*ScopeStorageKey]*StorageItemOption
}

func (o *ScopeStorageOption) Define(key *ScopeStorageKey, valType reflect.Type, optionalCallback StorageCopyCallback) {
	opt, exist := o.defines[key]
	if exist {
		panic("redefined storage key")
	}

	opt = &StorageItemOption{
		key:      key.id,
		valType:  valType,
		callback: optionalCallback,
	}
	o.defines[key] = opt
}

type DefaultScopeStorage struct {
	data    sync.Map
	Options config.Option[ScopeStorageOption] `gpfx.inject:""`
	Ctx     service.Provider                  `gpfx.inject:""`
}

func (storage *DefaultScopeStorage) ValidateKey(key *ScopeStorageKey) *StorageItemOption {
	if key == nil {
		panic("scoped local key is null")
	}
	opt, exist := storage.Options.OnceValue().defines[key]
	if !exist {
		panic("undefined key")
	}
	return opt
}

func (storage *DefaultScopeStorage) Get(key *ScopeStorageKey) (any, bool) {
	storage.ValidateKey(key)
	return storage.data.Load(key)
}

func (storage *DefaultScopeStorage) Set(key *ScopeStorageKey, newVal any) {
	opt := storage.ValidateKey(key)

	if !reflect.TypeOf(newVal).AssignableTo(opt.valType) {
		panic("The given value do not match to pre-defined value type")
	}

	storage.data.Store(key, newVal)
}

func (storage *DefaultScopeStorage) Delete(key *ScopeStorageKey) {
	storage.data.Delete(key)
}

func (storage *DefaultScopeStorage) CopyFrom(src *DefaultScopeStorage) {
	opt := storage.Options.OnceValue()

	src.data.Range(func(key, value any) bool {
		itemOpt, ok := opt.defines[key.(*ScopeStorageKey)]
		if !ok {
			panic("found undefined key while copying")
		}
		if itemOpt.callback != nil {
			newVal := itemOpt.callback(value, storage.Ctx)
			if !reflect.TypeOf(newVal).AssignableTo(itemOpt.valType) {
				panic("callback function returns a new value which do not match to pre-defined value type")
			}
			value = newVal
		}
		storage.data.Store(key, value)
		return true
	})
}
