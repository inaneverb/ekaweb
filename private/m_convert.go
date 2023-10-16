package ekaweb_private

import (
	"unsafe"
)

// ConvertMiddlewareToHandler converts given 'middleware' to Handler.
// It just asking middleware to get its Handler,
// giving EmptyHandler as a 'next' handler to its middleware.
func ConvertMiddlewareToHandler(middleware Middleware) Handler {
	return middleware.Callback(NewEmptyHandler())
}

// ConvertHandlerToMiddleware converts given 'handler' to Middleware.
// Handler that is obtained by returned Middleware
// will call given 'handler' followed by 'next' Handler of this Middleware.
func ConvertHandlerToMiddleware(handler Handler) Middleware {

	return MiddlewareFunc(func(next Handler) Handler {
		return mergeTwoHandlers(handler, next)
	})
}

// ConvertMiddlewaresToItsFuncs converts given 'middlewares' to list of
// MiddlewareFunc. It just builds array of Middleware.Callback functions.
func ConvertMiddlewaresToItsFuncs(middlewares []Middleware) []MiddlewareFunc {

	if len(middlewares) == 0 {
		return nil
	}

	var middlewareFunc = make([]MiddlewareFunc, 0, len(middlewares))
	for _, middleware := range middlewares {
		middlewareFunc = append(middlewareFunc, middleware.Callback)
	}

	return middlewareFunc
}

// ConvertMiddlewaresToRawFuncs converts given 'middlewares' to list of
// functions, each of them is implicitly MiddlewareFunc,
// but explicitly is just a function.
func ConvertMiddlewaresToRawFuncs(
	middlewares []Middleware) []func(next Handler) Handler {

	var middlewareFuncs = ConvertMiddlewaresToItsFuncs(middlewares)
	return ConvertMiddlewareFuncsToRawFuncs(middlewareFuncs)
}

// ConvertMiddlewareFuncsToRawFuncs converts given list of MiddlewareFunc
// to its underlying type (MiddlewareFunc is alias).
//
// Trivia:
// Despite the fact that having alias A for func F, it's allowed
// to convert these objects vice-versa, you cannot convert their slices
func ConvertMiddlewareFuncsToRawFuncs(
	middlewareFuncs []MiddlewareFunc) []func(next Handler) Handler {

	// SAFETY:
	// Unsafe here is safe until check below is OK.

	var _ func(next Handler) Handler = MiddlewareFunc(nil)
	return *(*[]func(next Handler) Handler)(unsafe.Pointer(&middlewareFuncs))
}
