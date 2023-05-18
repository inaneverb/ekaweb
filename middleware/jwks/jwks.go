package ekaweb_jwks

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/inaneverb/ekacore/ekaunsafe/v4"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"

	"github.com/inaneverb/ekaweb"
	"github.com/inaneverb/ekaweb/private"
)

type Middleware interface {
	ekaweb.Middleware
	Manual(ctx context.Context, token string)
}

// Middleware is a control center of the JWKS middleware.
// It's an object that is created by New() constructor.
// Later, its method Middleware.Callback() must be passed to the middleware registrar.
//
// Must not be instantiated manually and must not be public to the end user.
// It's a public just because it's a part of "internal" package.
type middleware struct {
	source Source

	skipErrorCheck bool

	fromOptions struct {
		globalCtx       context.Context
		log             ekaweb.Logger
		errSink         chan jwk.AutoRefreshError
		tokenExtractor  TokenExtractor
		tokenValidator  TokenAdditionalValidator
		refreshDelay    time.Duration
		fillErrorDetail bool
	}
}

const (
	HeaderPrefixDefault = "Bearer "
)

var (
	// ErrTokenNotFound is returned when an HTTP request doesn't have
	// 'Authorization' header or is malformed.
	ErrTokenNotFound = errors.New("middleware.Jwks: Token not found")

	// ErrTokenIncorrect is returned when JWT token is missing, malformed,
	// doesn't pass sign verification or incorrect encoded.
	ErrTokenIncorrect = errors.New("middleware.Jwks: Token is empty or incorrect")

	// ErrTokenInvalid is returned when JWT token is expired, not active yet, etc.
	ErrTokenInvalid = errors.New("middleware.Jwks: Token is expired or invalid")

	// ErrSourceUnavailable is returned when the source of Jwks
	// is unavailable to retrieve and JWT token validation is impossible.
	ErrSourceUnavailable = errors.New("middleware.Jwks: JWKS source is unavailable")

	// ----------

	errInitReqOptIsInvalid = errors.New("middleware.Jwks: Required option is invalid")
	errInitURLIsInvalid    = errors.New("middleware.Jwks: URL is empty or invalid")
	errInitContentIsEmpty  = errors.New("middleware.Jwks: Initialization content is empty")
	errInitKeyIsInvalid    = errors.New("middleware.Jwks: Initialization JWK is invalid")
	errInitSetIsInvalid    = errors.New("middleware.Jwks: Initialization JWK set is empty or invalid")
)

func (m *middleware) CheckErrorBefore() bool {
	return !m.skipErrorCheck
}

func (m *middleware) Callback(next ekaweb.Handler) ekaweb.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var token = m.fromOptions.tokenExtractor(r)
		m.Manual(r.Context(), token)

		next.ServeHTTP(w, r)
	})
}

func (m *middleware) Manual(ctx context.Context, token string) {

	if len(token) == 0 {
		ekaweb_private.UkvsInsertUserError(ctx, ErrTokenNotFound)
		return
	}

	// If raw token was found, save it, no matter valid it or not.
	// User may want to validate that token other way. Help them with that.

	SetTokenRawByContext(ctx, token)

	var tokenTyped, err = m.manual(ctx, token)
	if err != nil {
		ekaweb_private.UkvsInsertUserError(ctx, err)
		return
	}

	// Ok, token is parsed and validated.
	// It means the sign is correct and nbf, iat, exp fields are ok too.

	SetTokenByContext(ctx, tokenTyped)
}

func (m *middleware) manual(
	ctx context.Context, token string) (jwt.Token, error) {

	var jwkSet, err = m.source.GetJwks()
	switch {
	case err != nil && jwkSet != nil:
		// Ignore error if jwk.Set is present.
		// It's possible only if source is SourceAutoRefresh
		// and fetching the new version of jwk.Set is failed.
		// Thus, the old JWK set is returned along with the error.
		err = nil

	case err != nil:
		return nil, m.errorApply(ctx, err, ErrSourceUnavailable)
	}

	// We're dropping err of getting JWKS above because existing
	// the special option that allows caller to get all of these errors.

	var typedToken jwt.Token
	typedToken, err = jwt.Parse(
		ekaunsafe.StringToBytes(token),
		jwt.WithKeySet(jwkSet), jwt.UseDefaultKey(true),
	)
	if err != nil {
		return nil, m.errorApply(ctx, err, ErrTokenIncorrect)
	}

	err = jwt.Validate(typedToken)
	if err == nil {
		err = m.fromOptions.tokenValidator(ctx, typedToken)
	}

	if err != nil {
		return nil, m.errorApply(ctx, err, ErrTokenInvalid)
	}

	return typedToken, nil
}

// errorApply applies origErr as an error detail, if it's required.
// It also returns retErr.
func (m *middleware) errorApply(
	ctx context.Context, origErr, retErr error) error {

	// TODO: Log orig err

	if m.fromOptions.fillErrorDetail {
		ekaweb_private.UkvsInsertUserErrorDetail(ctx, origErr.Error())
	}
	return retErr
}

var _ ekaweb.Middleware = (*middleware)(nil)

////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

// New creates a new Middleware object based on provided options.
// The RequiredOption must be present, fatal error otherwise*. If RequiredOption
// is failed to apply it is also treated as fatal error*.
//
// After Middleware object is created it's ready to be used as a middleware
// providing JWKS based authorization and authentication.
// Pass Middleware.Callback() method to middleware registrar.
func New(reqOpt RequiredOption, opts []Option) Middleware {
	const s = "Failed to initialize JWKS middleware: %s"

	var m middleware

	WithContext(context.Background())(&m)
	WithTokenExtractorFromHeader(ekaweb.HeaderAuthorization)(&m)
	WithTokenAdditionalValidator(nil)(&m)

	for i, n := 0, len(opts); i < n; i++ {
		if opts[i] != nil {
			opts[i](&m)
		}
	}

	var err error

	if reqOpt != nil {
		err = reqOpt(&m)
	} else {
		err = errInitReqOptIsInvalid
	}

	if err != nil {
		if m.fromOptions.log != nil {
			m.fromOptions.log.Emerg(s, err.Error())
		}
		panic(err)
	}

	return &m
}

// InferAlgorithmForKey tries to figure out what jwa.SignatureAlgorithm
// should be used for the provided jwk.Key. Only EC algorithms are supported
// for now. Returns an empty string if inferring is failed.
func InferAlgorithmForKey(k jwk.Key) (jwa.SignatureAlgorithm, error) {

	// Maybe algorithm already present?

	if existedAlgorithm := k.Algorithm(); existedAlgorithm != "" {
		return jwa.SignatureAlgorithm(existedAlgorithm), nil
	}

	// Get the public key.
	// We will parse jwk.Key to the underlying public key structure.

	var err error
	k, err = k.PublicKey()
	if err != nil {
		return "", fmt.Errorf("failed to get public key: %w", err)
	}

	switch keyType := k.KeyType(); keyType {
	case jwa.EC:
		var rawPublicKey ecdsa.PublicKey
		if err = k.Raw(&rawPublicKey); err != nil {
			return "", fmt.Errorf("failed to extract raw part of EC public key: %w", err)
		}

		curveParams := rawPublicKey.Curve.Params()
		if curveParams == nil {
			return "", fmt.Errorf("failed to extract EC params from public key: %w", err)
		}

		switch bitSize := curveParams.BitSize; bitSize {
		case elliptic.P256().Params().BitSize:
			return jwa.ES256, nil

		case elliptic.P384().Params().BitSize:
			return jwa.ES384, nil

		case elliptic.P521().Params().BitSize:
			return jwa.ES512, nil

		default:
			return "", fmt.Errorf("unknown or unsupported EC bitsize: %d", bitSize)
		}

	default:
		return "", fmt.Errorf("unsupported key type: %v", keyType)
	}
}
