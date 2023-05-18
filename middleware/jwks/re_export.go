package ekaweb_jwks

import (
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
)

type (
	Key   = jwk.Key
	Set   = jwk.Set
	Token = jwt.Token
)

var (
	RegisterCustomField = jwt.RegisterCustomField
)
