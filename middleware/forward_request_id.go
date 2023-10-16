package ekaweb_middleware

import (
	"net/http"

	"github.com/inaneverb/ekaweb"
)

// ForwardRequestID returns a new HTTP middleware, that will copy header
// "X-Request-ID" from http.Request to the http.Response.
//
// If http.Request doesn't have such header, 'fallback' will be called
// to generate header's value and this header will be also applied to response.
//
// If 'fallback' is nil and no "X-Request-ID" header found in http.Request,
// then it will do nothing about copying headers.
//
// Calls 'next' handler in any case.
func ForwardRequestID(fallback func(r *http.Request) string) ekaweb.Middleware {

	var m = func(next ekaweb.Handler) ekaweb.Handler {
		return ekaweb.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			var reqID = r.Header.Get(ekaweb.HeaderXRequestID)
			if reqID == "" && fallback != nil {
				reqID = fallback(r)
			}

			if reqID != "" {
				w.Header().Set(ekaweb.HeaderXRequestID, reqID)
			}

			next.ServeHTTP(w, r)
		})
	}

	return ekaweb.MiddlewareFuncNoErrorCheck(m)
}
