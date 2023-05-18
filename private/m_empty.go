package ekaweb_private

import (
	"net/http"
)

// EmptyMiddleware is an empty middleware. It does nothing but returns
// a passed types.Handler thus granting flow control to that.
type EmptyMiddleware struct{}

func (*EmptyMiddleware) Callback(next Handler) Handler {
	return next
}

func NewEmptyMiddleware() Middleware {
	return &EmptyMiddleware{}
}

////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

// EmptyHandler is an empty types.Handler. It does nothing.
type EmptyHandler struct{}

func (*EmptyHandler) ServeHTTP(_ http.ResponseWriter, _ *http.Request) {}

func NewEmptyHandler() Handler {
	return &EmptyHandler{}
}
