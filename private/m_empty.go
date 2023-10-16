package ekaweb_private

import (
	"net/http"
)

// EmptyMiddleware is an empty Middleware. It just gives the control
// to the 'next' middleware. Has disabled error check before.
type EmptyMiddleware struct{}

func (*EmptyMiddleware) Callback(next Handler) Handler { return next }

func (*EmptyMiddleware) CheckErrorBefore() bool { return false }
func NewEmptyMiddleware() Middleware            { return &EmptyMiddleware{} }

////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

// EmptyHandler is just an empty Handler. It does nothing.
// Has disabled error check before.
type EmptyHandler struct{}

func (*EmptyHandler) ServeHTTP(http.ResponseWriter, *http.Request) {}

func (*EmptyHandler) CheckErrorBefore() bool { return false }
func NewEmptyHandler() Handler               { return &EmptyHandler{} }
