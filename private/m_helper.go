package ekaweb_private

import (
	"net/http"
)

// AsMiddleware normalizes provided middleware
// (or the component that may become) and returns it.
// So, if it's not a middleware you will get a nil object (passes nil check).
func AsMiddleware(v any) Middleware {
	switch v := v.(type) {
	case Middleware:
		return v

	case func(next Handler) Handler: // also handles MiddlewareFunc
		return MiddlewareFunc(v)

	default:
		return nil
	}
}

// AsHandler normalizes provided handler
// (or the component that may become) and returns it.
// So, if it's not a handler you will get a nil object.
func AsHandler(v any) Handler {
	switch v := v.(type) {
	case Handler:
		return v

	case func(w http.ResponseWriter, r *http.Request):
		return http.HandlerFunc(v)

	default:
		return nil
	}
}
