package ekaweb_jwks

import (
	"fmt"
	"time"

	"github.com/lestrrat-go/jwx/jwk"
)

// RequiredOption is a special option
// that is able to set a Source that is used in the Middleware.
type RequiredOption func(s *middleware) error

// WithPEMContent returns a new RequiredOption
// that will offer a SourceCommon jwk.Set provider to be used in Middleware.
// The created jwk.Set will contain exactly 1 key and it will be parsed
// using provided PEM data. If it's private key, the public one will be used.
func WithPEMContent(content []byte) RequiredOption {

	if len(content) == 0 {
		return withErr(errInitContentIsEmpty)
	}

	if k, err := jwk.ParseKey(content, jwk.WithPEM(true)); err == nil {
		return WithJwk(k)
	} else {
		return withErr(err)
	}
}

// WithJwkContent returns a new RequiredOption
// that will offer a SourceCommon jwk.Set provider to be used in Middleware.
// The created jwk.Set will contain exactly 1 key and it will be parsed
// using provided data of that key.
// It must be JSON. If it's private key, the public one will be used.
func WithJwkContent(content []byte) RequiredOption {

	if len(content) == 0 {
		return withErr(errInitContentIsEmpty)
	}

	if k, err := jwk.ParseKey(content); err == nil {
		return WithJwk(k)

	} else if s, err2 := jwk.Parse(content); err2 == nil {
		return WithJwks(s)

	} else {
		return withErr(err)
	}
}

// WithJwk returns a new RequiredOption
// that will offer a SourceCommon jwk.Set provider to be used in Middleware.
// The created jwk.Set will contain exactly 1 key - the provided one by argument.
// IIf it's private key, the public one will be used.
func WithJwk(k jwk.Key) RequiredOption {

	if k == nil {
		return withErr(errInitKeyIsInvalid)
	}

	k, err := k.PublicKey()
	if err != nil {
		return withErr(err)
	}

	s := jwk.NewSet()
	s.Add(k)

	return WithJwks(s)
}

// WithJwks returns a new RequiredOption
// that will offer a SourceCommon jwk.Set provider to be used in Middleware.
// It will provide the jwk.Set that is passed to this function.
func WithJwks(s jwk.Set) RequiredOption {

	if s == nil || s.Len() == 0 {
		return withErr(errInitSetIsInvalid)
	}

	approved := jwk.NewSet()
	for i, n := 0, s.Len(); i < n; i++ {
		keyToBeApproved, _ := s.Get(i)

		keyToBeApproved, err := keyToBeApproved.PublicKey()
		if err != nil {
			err = fmt.Errorf("failed to get a public key for JWK[%d]: %w", i, err)
			return withErr(err)
		}

		algo, err := InferAlgorithmForKey(keyToBeApproved)
		if err != nil {
			err = fmt.Errorf("failed to infer JWA for JWK[%d]: %w", i, err)
			return withErr(err)
		}

		_ = keyToBeApproved.Set(jwk.AlgorithmKey, algo)
		approved.Add(keyToBeApproved)
	}

	source, err := NewSourceCommon(approved)
	if err != nil {
		return withErr(err)
	}

	return func(s *middleware) error { s.source = source; return nil }
}

// WithAutoRefresh returns a new RequiredOption
// that will offer a SourceAutoRefresh jwk.Set provider to be used in Middleware.
// It will save provided context and URL to be used later but also constructs
// a new jwk.AutoRefresh object based on provided URL.
func WithAutoRefresh(url string) RequiredOption {
	return func(s *middleware) error {

		if url == "" {
			return errInitURLIsInvalid
		}

		opts := make([]jwk.AutoRefreshOption, 0, 10)
		if s.fromOptions.refreshDelay >= 10*time.Second {
			opts = append(opts, jwk.WithMinRefreshInterval(s.fromOptions.refreshDelay))
		}

		source, err := NewSourceAutoRefresh(
			s.fromOptions.globalCtx, url, s.fromOptions.errSink, opts...)

		if err != nil {
			return err
		}

		s.source = source
		return nil
	}
}

// withErr is just a helper function that returns RequiredOption that will return
// provided error.
func withErr(err error) RequiredOption {
	return func(_ *middleware) error { return err }
}
