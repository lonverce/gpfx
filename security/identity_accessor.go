package security

import (
	"github.com/lonverce/gpfx"
)

type IdentityAccessor interface {
	GetIdentity() (identity *ClaimIdentity, exist bool)
	Change(identity *ClaimIdentity) (disposer func())
}

var identityStorageKey = gpfx.NewScopeStorageKey()

type internalIdentityAccessor struct {
	Storage gpfx.ScopeStorage `gpfx.inject:""`
}

func (i *internalIdentityAccessor) GetIdentity() (identity *ClaimIdentity, exist bool) {
	v, ok := i.Storage.Get(identityStorageKey)
	if !ok {
		return nil, false
	}
	return v.(*ClaimIdentity), true
}

func (i *internalIdentityAccessor) Change(identity *ClaimIdentity) (disposer func()) {
	v, ok := i.Storage.Get(identityStorageKey)

	i.Storage.Set(identityStorageKey, identity)

	if ok {
		oldIdentity := v.(*ClaimIdentity)
		return func() {
			i.Storage.Set(identityStorageKey, oldIdentity)
		}
	} else {
		return func() {
			i.Storage.Delete(identityStorageKey)
		}
	}
}
