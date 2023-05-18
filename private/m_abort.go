package ekaweb_private

import (
	"net/http"
)

// AbortWith generates a middleware that is always returns passed error
// failing the whole process of handling HTTP request.
func AbortWith(err error) Middleware {
	return MiddlewareFuncNoErrorCheck(func(next Handler) Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			UkvsInsertUserError(r.Context(), err)
		})
	})
}
