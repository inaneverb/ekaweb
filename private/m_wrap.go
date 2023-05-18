package ekaweb_private

import (
	"net/http"
)

func HandlerBeforeNext(handler Handler) Middleware {
	return MiddlewareFunc(func(next Handler) Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handler.ServeHTTP(w, r)
			next.ServeHTTP(w, r)
		})
	})
}

func HandlerAfterNext(handler Handler) Middleware {
	return MiddlewareFunc(func(next Handler) Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
			handler.ServeHTTP(w, r)
		})
	})
}
