package ekaweb_private

import (
	"fmt"
	"net/http"
)

// Recover generates a middleware that prevents panics from next HTTP controllers.
// It wraps these calls by the deferring recover() call and if panic is recovered,
// it transforms it to the error and returns it.
func Recover() Middleware {
	return MiddlewareFuncNoErrorCheck(func(next Handler) Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if recovered := recover(); recovered != nil {
					var err = fmt.Errorf("PANIC RECOVERED: %+v", recovered)
					UkvsInsertUserError(r.Context(), err)
				}
			}()

			next.ServeHTTP(w, r)
		})
	})
}
