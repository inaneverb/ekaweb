package ekaweb_jwks

import (
	"context"
	"time"

	"github.com/inaneverb/ekaweb/v2"
)

// Option is a callback that allows to modify Middleware under its construction.
type Option func(m *middleware)

// WithContext returns an Option that will be useful if Middleware is created
// with WithAutoRefresh() option. Otherwise it's useless.
// Allows to specify a context, cancelling which will lead to stop refreshing jwk.Set.
//
// If this option is not specified, default context.Background() is used,
// meaning that the process of refreshing jwk.Set will be stopped only at the os.Exit(0),
// and w/o any gracefully.
func WithContext(ctx context.Context) Option {
	return func(m *middleware) {
		if ctx != nil {
			m.fromOptions.globalCtx = ctx
		}
	}
}

// WithRefreshDelay returns an Option that will be useful if Middleware is created
// with WithAutoRefresh() option. Otherwise it's useless.
// Allows to specify a delay after which previous fetched jwk.Set is considered
// outdated and must be refreshed.
//
// If delay is not provided, a jwk.AutoRefresh has an excellent heuristic mechanism
// that allows refresh jwk.Set on-demand based on HTTP headers of response
// from jwk.Set server. Read more about that in docs of jwk.AutoRefresh.
// In 99% cases it's OK if you do not provide delay.
//
// If delay is not provided and HTTP response from the jwk.Set generator
// does not contains necessary headers, it will fallback to default value: 1h.
func WithRefreshDelay(delay time.Duration) Option {
	return func(m *middleware) {
		if delay >= 10*time.Second {
			m.fromOptions.refreshDelay = delay
		}
	}
}

// WithLogger returns an Option that allows you to register a logger
// that will be used to write an error or fatal messages produced by the middleware.
//
// If logger is not specified, error messages are ignored, fatal are replaced by panic.
func WithLogger(log ekaweb.Logger) Option {
	return func(m *middleware) {
		if log != nil {
			m.fromOptions.log = log
		}
	}
}

// WithTokenExtractors returns an Option that allows you to register a many
// TokenExtractor callbacks which will be used to extract a token
// from the http.Request one-by-one until they're succeeded.
func WithTokenExtractors(extractors ...TokenExtractor) Option {
	return func(m *middleware) {
		m.fromOptions.tokenExtractor = tokenExtractorMerge(extractors)
	}
}

// WithTokenExtractorFromHeader returns an Option that allows you to register
// TokenExtractor callback that will extract a token from the http.Request header
// with the provided `key`.
func WithTokenExtractorFromHeader(key string) Option {
	return WithTokenExtractorFromHeaderExtended(key, HeaderPrefixDefault)
}

// WithTokenExtractorFromHeaderExtended is the same as just WithTokenExtractorFromHeader
// but also allows you to specify `skipPrefix`. It allows you to ignore a part of
// the header's value.
func WithTokenExtractorFromHeaderExtended(key, skipPrefix string) Option {
	return func(m *middleware) {
		m.fromOptions.tokenExtractor = tokenExtractorFromHeader(key, skipPrefix)
	}
}

// WithTokenExtractorFromQuery returns an Option that allows you to register
// TokenExtractor callback that will extract a token from the http.Request URL
// values with the provided `key`.
func WithTokenExtractorFromQuery(key string) Option {
	return WithTokenExtractorFromQueryExtended(key, "")
}

// WithTokenExtractorFromQueryExtended is the same as just WithTokenExtractorFromQuery
// but also allows you to specify `skipPrefix`. It allows you to ignore a part of
// the URL value.
func WithTokenExtractorFromQueryExtended(key, skipPrefix string) Option {
	return func(m *middleware) {
		m.fromOptions.tokenExtractor = tokenExtractorFromQuery(key, skipPrefix)
	}
}

// WithSkipErrorCheckBefore affects an HTTP middleware call flow at all
// not a JWKS middleware. A JWKS middleware will tell a "call manager"
// to not wrap it by the error check.
func WithSkipErrorCheckBefore() Option {
	return func(m *middleware) {
		m.skipErrorCheck = true
	}
}

// WithTokenAdditionalValidator returns an Option that allows you to register
// TokenAdditionalValidator - a custom callback, that allows you to perform
// additional token checks right before a middleware job is finished.
//
// If your callback reports that token is invalid, ErrTokenInvalid is applied.
func WithTokenAdditionalValidator(cb TokenAdditionalValidator) Option {
	return func(m *middleware) {
		if cb != nil {
			m.fromOptions.tokenValidator = cb
		} else {
			m.fromOptions.tokenValidator = tokenAdditionalValidatorEmpty
		}
	}
}

// WithUseUnderlyingErrorAsErrorDetail returns an Option that allows you
// to fill error detail (accessed using http.Request's context getters)
// by an original error.
//
// E.g:
// You may get ErrTokenInvalid, ErrTokenIncorrect as an error.
// But you don't know what exactly happened. If you will specify this option,
// a text of underlying error will be applied as an error detail.
func WithUseUnderlyingErrorAsErrorDetail() Option {
	return func(m *middleware) {
		m.fromOptions.fillErrorDetail = true
	}
}
