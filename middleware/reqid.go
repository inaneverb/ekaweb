package ekaweb_middleware

import (
	"net/http"

	"github.com/inaneverb/ekaweb"
	"github.com/inaneverb/ekaweb/private"
)

func ForwardRequestID(fallback func(r *http.Request) string) ekaweb_private.Middleware {
	return ekaweb_private.MiddlewareFunc(func(next ekaweb_private.Handler) ekaweb_private.Handler {
		return ekaweb_private.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			var reqID = r.Header.Get(ekaweb.HeaderXRequestID)
			if reqID == "" && fallback != nil {
				reqID = fallback(r)
			}

			if reqID != "" {
				w.Header().Set(ekaweb.HeaderXRequestID, reqID)
			}

			next.ServeHTTP(w, r)
		})
	})
}
