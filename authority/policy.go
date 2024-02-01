package authority

import (
	"github.com/lonverce/gpfx/config"
	"github.com/lonverce/gpfx/security"
	"github.com/lonverce/gpfx/service"
)

type policyDeclare struct {
	validators []PolicyValidator
}

type Option struct {
	items map[string]*policyDeclare
}

func (p *Option) AddPolicy(policyName string, validators ...PolicyValidator) {
	if policyName == "" {
		panic("policyName can not be empty")
	}
	d := &policyDeclare{
		validators: validators,
	}
	p.items[policyName] = d
}

type ValidationContext interface {
	Provider() service.Provider
	Approve()
	Reject(detail string)
}

type PolicyValidator interface {
	Validate(identity *security.ClaimIdentity, context ValidationContext)
}

type PolicyManager interface {
	GetAllValidatorForPolicy(policyName string) []PolicyValidator
}

type DefaultPolicyManager struct {
	Option config.Option[Option] `gpfx.inject:""`
}

func (mgr *DefaultPolicyManager) GetAllValidatorForPolicy(policyName string) []PolicyValidator {
	declare, ok := mgr.Option.OnceValue().items[policyName]
	if ok {
		return declare.validators
	}
	return nil
}
