package ekaweb_private

import (
	"net/http"

	"github.com/inaneverb/ekacore/ekaunsafe/v4"
)

// MiddlewareCheckError is HTTP middleware that ensures HTTP request
// has not been failed with error yet. It checks an error BEFORE executing
// next handler. The handler won't be executed if error is occurred. Nil safe.
type MiddlewareCheckError struct{}

// CheckErrorBefore always returns false to avoid recursive loop lock.
func (*MiddlewareCheckError) CheckErrorBefore() bool { return false }

// Callback is a middleware implementation, that executes 'next' Handler
// only if there's no saved occurred error in http.Request's context.Context.
func (*MiddlewareCheckError) Callback(next Handler) Handler {

	if ekaunsafe.UnpackInterface(next).Word == nil {
		return NewEmptyHandler()
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := UkvsGetUserError(r.Context()); err == nil {
			next.ServeHTTP(w, r)
		}
	})
}

// NewCheckErrorMiddleware returns a Middleware, that is used to check whether
// an error is occurred. So, this Middleware is often used as a 'checkError'
// in BuildHandlerOut() builder.
func NewCheckErrorMiddleware() Middleware { return &MiddlewareCheckError{} }

////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

// MiddlewareErrorHandler is HTTP middleware that calls given ErrorHandlerHTTP
// if there's an error in http.Request's context.Context.
type MiddlewareErrorHandler struct{ errorHandler ErrorHandlerHTTP }

// CheckErrorBefore returns false, because we want to get control flow
// to handle an error. Otherwise, the control flow may be stopped before
// this Middleware will get its process slot.
func (m *MiddlewareErrorHandler) CheckErrorBefore() bool { return false }

// Callback is a middleware implementation, that calls given 'next' Handler
// and then checks, whether an error is occurred. If so, calls stored
// ErrorHandlerHTTP.
func (m *MiddlewareErrorHandler) Callback(next Handler) Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		if err := UkvsGetUserError(r.Context()); err != nil {
			m.errorHandler(w, r, err)
		}
	})
}

// NewErrorHandlerMiddleware returns a Middleware, that will execute given
// ErrorHandlerHTTP after executing chain middleware(s) or/and handler(s)
// and only if there's an error is occurred.
func NewErrorHandlerMiddleware(errorHandler ErrorHandlerHTTP) Middleware {
	if errorHandler != nil {
		return &MiddlewareErrorHandler{errorHandler}
	} else {
		return NewEmptyMiddleware()
	}
}
