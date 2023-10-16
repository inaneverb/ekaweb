package ekaweb_middleware

import (
	"net/http"

	"github.com/inaneverb/ekaweb"
	"github.com/inaneverb/ekaweb/private"
)

// AbortWith returns a new HTTP middleware, that always fails processing
// incoming HTTP request with given error.
// This error will be stored to http.Request's context.Context.
func AbortWith(err error) ekaweb.Middleware {

	var m = func(next ekaweb.Handler) ekaweb.Handler {
		return ekaweb.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ekaweb_private.UkvsInsertUserError(r.Context(), err)
		})
	}

	return ekaweb.MiddlewareFuncNoErrorCheck(m)
}

// AbortIf returns a new HTTP middleware, that fails processing incoming HTTP
// request if given 'cb' returns not-nil error.
// This error will be stored to http.Request's context.Context.
func AbortIf(cb func(r *http.Request) error) ekaweb.Middleware {

	var m = func(next ekaweb.Handler) ekaweb.Handler {
		return ekaweb.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if err := cb(r); err != nil {
				ekaweb_private.UkvsInsertUserError(r.Context(), err)
			} else {
				next.ServeHTTP(w, r)
			}
		})
	}

	return ekaweb.MiddlewareFuncNoErrorCheck(m)
}
