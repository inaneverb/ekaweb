package ekaweb_private

import (
	"net/http"

	"github.com/inaneverb/ekacore/ekaunsafe/v4"
)

// MiddlewareCheckError is an HTTP middleware that ensures HTTP request
// has not been failed with error yet. It checks an error
// BEFORE executing a next handler. The handler won't be executed
// if error is occurred. Nil safe.
type MiddlewareCheckError struct{}

func (*MiddlewareCheckError) CheckErrorBefore() bool {
	return false
}

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

func NewCheckErrorMiddleware() Middleware {
	return &MiddlewareCheckError{}
}

////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

// MiddlewareErrorHandler is an HTTP middleware
// that calls saved types.ErrorHandlerHTTP if there's an error in http.Request's
// context.
type MiddlewareErrorHandler struct {
	errorHandler ErrorHandlerHTTP
}

func (m *MiddlewareErrorHandler) CheckErrorBefore() bool {
	return false
}

func (m *MiddlewareErrorHandler) Callback(next Handler) Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		next.ServeHTTP(w, r)

		err := UkvsGetUserError(r.Context())
		if err == nil {
			return
		}

		m.errorHandler(w, r, err)
	})
}

func NewErrorHandlerMiddleware(errorHandler ErrorHandlerHTTP) Middleware {
	if errorHandler != nil {
		return &MiddlewareErrorHandler{errorHandler}
	} else {
		return NewEmptyMiddleware()
	}
}
