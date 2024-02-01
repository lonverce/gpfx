package security

type CurrentUser interface {
	IsAuthenticated() bool
	Id() string
	Name() string
	ClientId() string
}

type DefaultCurrentUser struct {
	IdentityAccessor IdentityAccessor `gpfx.inject:""`
}

func (user *DefaultCurrentUser) ClientId() string {
	v, found := user.FindClaim("client_id")
	if !found {
		return ""
	}

	return v
}

func (user *DefaultCurrentUser) Id() string {
	v, found := user.FindClaim("sid")
	if !found {
		return ""
	}

	return v
}

func (user *DefaultCurrentUser) Name() string {
	v, found := user.FindClaim("name")
	if !found {
		return ""
	}

	return v
}

func (user *DefaultCurrentUser) Scopes() []string {
	return user.FindAllClaims("scope")
}

func (user *DefaultCurrentUser) IsAuthenticated() bool {
	_, ok := user.IdentityAccessor.GetIdentity()
	return ok
}

func (user *DefaultCurrentUser) FindClaim(name string) (string, bool) {
	identity, ok := user.IdentityAccessor.GetIdentity()
	if !ok {
		return "", false
	}

	for _, claim := range identity.Claims() {
		if claim.Name() == name {
			return claim.Value(), true
		}
	}

	return "", false
}

func (user *DefaultCurrentUser) FindAllClaims(name string) []string {
	identity, ok := user.IdentityAccessor.GetIdentity()
	arr := make([]string, 0)

	if ok {
		for _, claim := range identity.Claims() {
			if claim.Name() == name {
				arr = append(arr, claim.Value())
			}
		}
	}

	return arr
}
