package ekaweb_middleware

import (
	"net/http"

	"github.com/inaneverb/ekaweb"
	"github.com/inaneverb/ekaweb/private"
)

// CustomHeaders returns a new HTTP middleware, that applies given 'headers'
// to each http.Response. Then calls 'next' handler. Has no error check before.
func CustomHeaders(headers http.Header) ekaweb.Middleware {

	var m = func(next ekaweb.Handler) ekaweb.Handler {

		if len(headers) == 0 {
			return next
		}

		return ekaweb.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ekaweb.HeadersMerge(w.Header(), headers, true)
			next.ServeHTTP(w, r)
		})
	}

	return ekaweb.MiddlewareFuncNoErrorCheck(m)
}

// CustomHeadersIf is almost the same as just CustomHeaders(), but allows you
// to implement difficult logic with conditional applying custom headers
// depend on given http.Request. Has no error check before.
func CustomHeadersIf(cb func(r *http.Request) http.Header) ekaweb.Middleware {

	if cb == nil {
		return ekaweb_private.NewEmptyMiddleware()
	}

	var m = func(next ekaweb.Handler) ekaweb.Handler {
		return ekaweb.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			var headers = cb(r)
			if len(headers) != 0 {
				ekaweb.HeadersMerge(w.Header(), headers, true)
			}

			next.ServeHTTP(w, r)
		})
	}

	return ekaweb.MiddlewareFuncNoErrorCheck(m)
}
