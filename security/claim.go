package security

type Claim struct {
	name  string
	value string
}

func NewClaim(name, value string) *Claim {
	return &Claim{
		name:  name,
		value: value,
	}
}

func (c *Claim) Name() string {
	return c.name
}

func (c *Claim) Value() string {
	return c.value
}

type ClaimIdentity struct {
	provider string
	claims   []*Claim
}

func NewClaimIdentity(provider string, claims []*Claim) *ClaimIdentity {
	return &ClaimIdentity{
		provider: provider,
		claims:   claims[:],
	}
}

func (c *ClaimIdentity) Provider() string {
	return c.provider
}

func (c *ClaimIdentity) Claims() []*Claim {
	return c.claims[:]
}
