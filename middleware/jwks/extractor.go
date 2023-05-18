package ekaweb_jwks

import (
	"net/http"
	"strings"
)

// TokenExtractor is a func alias that should take http.Request and return JWT.
// It allows user to overwrite where token should be taken
// (header or somewhere else).
type TokenExtractor = func(r *http.Request) string

// tokenExtractorFromHeader returns a TokenExtractor that will extract a token
// from the http.Request's headers. It takes header value by the provided `key`
// and trims a value by the `skipPrefix` from the left.
func tokenExtractorFromHeader(key, skipPrefix string) TokenExtractor {

	key = strings.TrimSpace(key)
	if key == "" {
		return tokenExtractorEmpty
	}

	return func(r *http.Request) string {
		return tokenSkipPrefix(r.Header.Get(key), skipPrefix)
	}
}

// tokenExtractorFromQuery returns a TokenExtractor that will extract a token
// from the http.Request's URL values. It takes a URL value by the provided
// `key` and trims a value by the `skipPrefix` from the left.
func tokenExtractorFromQuery(key, skipPrefix string) TokenExtractor {

	key = strings.TrimSpace(key)
	if key == "" {
		return tokenExtractorEmpty
	}

	return func(r *http.Request) string {
		return tokenSkipPrefix(r.URL.Query().Get(key), skipPrefix)
	}
}

// tokenExtractorMerge returns a TokenExtractor that will call each
// TokenExtractor from the provided `extractors` until not empty token
// is returned.
func tokenExtractorMerge(extractors []TokenExtractor) TokenExtractor {

	extractorsBak := make([]TokenExtractor, 0, len(extractors))
	extractors, extractorsBak = extractorsBak, extractors

	for i, n := 0, len(extractorsBak); i < n; i++ {
		if extractorsBak[i] != nil {
			extractors = append(extractors, extractorsBak[i])
		}
	}

	switch len(extractors) {
	case 0:
		return tokenExtractorEmpty
	case 1:
		return extractors[0]
	}

	return func(r *http.Request) string {
		token := extractors[0](r)
		for i := 1; i < len(extractors) && token == ""; i++ {
			token = extractors[i](r)
		}
		return token
	}
}

// tokenExtractorEmpty is an empty TokenExtractor. Just returns an empty string.
func tokenExtractorEmpty(_ *http.Request) string {
	return ""
}

// tokenSkipPrefix ensures that provided `token` has a required `skipPrefix`.
// If token don't have this prefix, an empty string is returned.
// An empty prefix leads for passed token being returned as is.
func tokenSkipPrefix(token, skipPrefix string) string {
	if strings.HasPrefix(token, skipPrefix) {
		return token[len(skipPrefix):]
	} else {
		return ""
	}
}
