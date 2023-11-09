package ekaweb_middleware

import (
	"fmt"
	"net/http"

	"github.com/inaneverb/ekaweb/v2"
	"github.com/inaneverb/ekaweb/v2/private"
)

// Recover returns a new HTTP middleware, that prevents panics from 'next'
// HTTP handler.
// If panic is occurred during executing next handler, it will be recovered,
// transformed into error and saved to http.Request's context.Context.
func Recover() ekaweb.Middleware {

	var m = func(next ekaweb.Handler) ekaweb.Handler {
		return ekaweb.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			defer func() {
				if recovered := recover(); recovered != nil {
					var err = fmt.Errorf("PANIC RECOVERED: %+v", recovered)
					ekaweb_private.UkvsInsertUserError(r.Context(), err)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}

	return ekaweb.MiddlewareFuncNoErrorCheck(m)
}
