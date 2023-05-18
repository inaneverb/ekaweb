package ekaweb_jwks

import (
	"bytes"
	"context"
	"encoding/base64"

	"github.com/goccy/go-json"
	"github.com/inaneverb/ekacore/ekaunsafe/v4"
	"github.com/lestrrat-go/jwx/jwt"

	"github.com/inaneverb/ekaweb/private"
)

type (
	ukvsKeyTokenRaw   struct{} // Ukvs key to store token in RAW format (string)
	ukvsKeyTokenKid   struct{} // Ukvs key to store token's KID (cache)
	ukvsKeyTokenTyped struct{} // Ukvs key to store typed token (jwks format)
)

// GetTokenRawByContext retrieves the raw token from the user's scope
// of context.Context. If raw token has not been set the empty string is returned.
//
// GetTokenRawByContext always returns a token in the same format
// as it was being a part of HTTP 'Authorization' header
// but only if JWKS middleware has been passed successfully.
func GetTokenRawByContext(ctx context.Context) string {
	return ekaweb_private.UkvsGetOrDefault(ctx, ukvsKeyTokenRaw{}, "").(string)
}

// SetTokenRawByContext saves the raw token into the user's scope
// of context.Context. If token was set already it will be overwritten.
func SetTokenRawByContext(ctx context.Context, tokenRaw string) {
	ekaweb_private.UkvsInsert(ctx, ukvsKeyTokenRaw{}, tokenRaw)
}

// GetTokenByContext retrieves the token from the user's scope
// of the context.Context. If token has not been set the nil is returned.
//
// GetTokenByContext always returns a token
// if the JWKS middleware has been completed successfully.
func GetTokenByContext(ctx context.Context) jwt.Token {

	// We can't use UkvsGetOrDefault here, since jwt.Token is an interface.
	// And nil typed interface is the same as nil non-typed interface,
	// so types assertion will fail.

	if typedToken, ok := ekaweb_private.UkvsLookup(ctx, ukvsKeyTokenTyped{}); ok {
		return typedToken.(jwt.Token)
	} else {
		return nil
	}
}

// GetTokenKidByContext retrieves the token from the user's scope
// of the context.Context and then extracts "kid" field from the JWT token headers.
// If token has not been set or "kid" is empty, the empty string is returned.
func GetTokenKidByContext(ctx context.Context) string {

	// Maybe we're already did extract it some time before?

	var kidCachedUntyped, isKidCached = ekaweb_private.UkvsLookup(ctx, ukvsKeyTokenKid{})
	if isKidCached {
		return kidCachedUntyped.(string)
	}

	// Ok, slow case.

	var rawToken = GetTokenRawByContext(ctx)
	if rawToken == "" {
		return ""
	}

	var kidReal = GetTokenKidByRawToken(rawToken)
	ekaweb_private.UkvsInsert(ctx, ukvsKeyTokenKid{}, kidReal)

	return kidReal
}

// GetTokenKidByRawToken extracts "kid" field's value from the RAW JWT token.
// If either token or "kid" is empty, the empty string is returned.
//
// It's preferred to call GetTokenKidByContext() instead, since it caches
// the token kid to the context. Call this method only if you have no token
// in your context, and you cannot place it there.
func GetTokenKidByRawToken(rawToken string) string {

	// https://github.com/lestrrat-go/jwx/discussions/547

	type JwtHeader struct {
		KeyID string `json:"kid"`
	}

	var headerBytes = ekaunsafe.StringToBytes(rawToken)
	var separatorIdx = bytes.IndexByte(headerBytes, '.')

	if separatorIdx > -1 {
		headerBytes = headerBytes[:separatorIdx]
	}

	// Decode Base64.

	var decodedHeaders = make([]byte, base64.RawStdEncoding.DecodedLen(separatorIdx))
	_, _ = base64.StdEncoding.Decode(decodedHeaders, headerBytes)

	// And JSON.

	var headers JwtHeader

	_ = json.Unmarshal(decodedHeaders, &headers)
	return headers.KeyID
}

// SetTokenByContext saves the token into the user's scope of context.Context.
// If token was set already it will be overwritten.
func SetTokenByContext(ctx context.Context, token jwt.Token) {
	ekaweb_private.UkvsInsert(ctx, ukvsKeyTokenTyped{}, token)
}
