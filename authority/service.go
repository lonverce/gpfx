package authority

import (
	"github.com/lonverce/gpfx/security"
	"github.com/lonverce/gpfx/service"
)

type Service interface {
	Check(policyName string) bool
}

type DefaultService struct {
	PolicyMgr PolicyManager             `gpfx.inject:""`
	Identity  security.IdentityAccessor `gpfx.inject:""`
	Provider  service.Provider          `gpfx.inject:""`
}

func (srv *DefaultService) Check(policyName string) bool {
	identity, ok := srv.Identity.GetIdentity()

	if !ok {
		return false
	}

	validators := srv.PolicyMgr.GetAllValidatorForPolicy(policyName)

	if len(validators) == 0 {
		return false
	}

	ctx := &defaultValidationContext{
		serviceProvider: srv.Provider,
	}

	approvedCount := 0

	for _, validator := range validators {
		validator.Validate(identity, ctx)

		if ctx.rejected {
			return false
		}

		if ctx.approved {
			approvedCount++
			ctx.approved = false
		}
	}

	if approvedCount == 0 {
		return false
	}

	return true
}

type defaultValidationContext struct {
	serviceProvider service.Provider
	rejected        bool
	approved        bool
	detail          string
}

func (d *defaultValidationContext) Provider() service.Provider {
	return d.serviceProvider
}

func (d *defaultValidationContext) Approve() {
	d.approved = true
}

func (d *defaultValidationContext) Reject(detail string) {
	d.rejected = true
	d.detail = detail
}
