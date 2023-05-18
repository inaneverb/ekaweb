package ekaweb_jwks

import (
	"net/http"
)

/*
JWKS is a middleware that allows to perform authentication and authorization
using JSON Web Tokens and related standards: RFC 7515-7519.
Thus there's also provided a way to specify JWK or even JWKS as a source
of private keys. Moreover you can enable auto-refreshing JWKS over time.

First of all you should
*/

////////////////////////////////////////////////////////////////////////////////
///// Middleware generators ////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

// NewByPEM returns a new JWKS middleware that uses provided byte array
// as a PEM encoded certificate JWK public key.
//
// Since JwksMiddleware implements types.Middleware,
// you may pass it to types.Router as is directly (without any conversion).
func NewByPEM(content []byte, options ...Option) Middleware {
	return New(WithPEMContent(content), options)
}

// NewbyJwkContent returns a new JWKS middleware that uses provided byte array
// as a raw data of JWK or its set (thus it must be JSON encoded data).
// Read more about format: https://datatracker.ietf.org/doc/html/rfc7517#section-4
//
// Since JwksMiddleware implements types.Middleware,
// you may pass it to types.Router as is directly (without any conversion).
func NewbyJwkContent(content []byte, options ...Option) Middleware {
	return New(WithJwkContent(content), options)
}

// NewByJwk returns a new JWKS middleware that allow manually provide
// a typed JWK from" the https://github.com/lestrrat-go/jwx package.
//
// Since JwksMiddleware implements types.Middleware,
// you may pass it to types.Router as is directly (without any conversion).
func NewByJwk(k Key, options ...Option) Middleware {
	return New(WithJwk(k), options)
}

// NewByJwks returns a new JWKS middleware that allow manually provide
// a typed JWK set from the https://github.com/lestrrat-go/jwx package.
//
// Since JwksMiddleware implements types.Middleware,
// you may pass it to types.Router as is directly (without any conversion).
func NewByJwks(s Set, options ...Option) Middleware {
	return New(WithJwks(s), options)
}

// JwksByAutoRefresh returns a new JWKS middleware that allow to obtain
// a JWKS from the remote HTTP resource and then refresh them automatically
// in the background (with thread-safety) by the some period.
//
// Use custom options to specify interval and behaviour or read internal docs
// to know defaults.
//
// Since JwksMiddleware implements types.Middleware,
// you may pass it to types.Router as is directly (without any conversion).
func JwksByAutoRefresh(url string, options ...Option) Middleware {
	return New(WithAutoRefresh(url), options)
}

////////////////////////////////////////////////////////////////////////////////
///// Extractors from the HTTP Request or context.Context //////////////////////
////////////////////////////////////////////////////////////////////////////////

// GetTokenRaw is an extractor. It allows you to get a JWT raw token
// in the next middleware or handler after JWKS middleware has been passed.
func GetTokenRaw(r *http.Request) string {
	return GetTokenRawByContext(r.Context())
}

// GetToken is an extractor.
// It allows you to get a JWT token in the middleware or handler
// after JWKS middleware has been passed.
func GetToken(r *http.Request) Token {
	return GetTokenByContext(r.Context())
}

// GetTokenKid is an extractor.
// It allows you to get a JWT token's "kid" header in the middleware or handler
// after JWKS middleware has been passed.
func GetTokenKid(r *http.Request) string {
	return GetTokenKidByContext(r.Context())
}
