package ekaweb_private

import (
	"unsafe"
)

func ConvertMiddlewareToHandler(middleware Middleware) Handler {
	return middleware.Callback(NewEmptyHandler())
}

func ConvertHandlerToMiddleware(handler Handler) Middleware {
	return MiddlewareFunc(func(next Handler) Handler {
		return mergeTwoHandlers(handler, next)
	})
}

func ConvertMiddlewaresToItsFunc(middlewares []Middleware) []MiddlewareFunc {

	if len(middlewares) == 0 {
		return nil
	}

	middlewaresFunc := make([]MiddlewareFunc, 0, len(middlewares))
	for i, n := 0, len(middlewares); i < n; i++ {
		middlewaresFunc = append(middlewaresFunc, middlewares[i].Callback)
	}

	return middlewaresFunc
}

func ConvertMiddlewaresToRawFuncs(
	middlewares []Middleware) []func(next Handler) Handler {

	middlewareFuncs := ConvertMiddlewaresToItsFunc(middlewares)
	return ConvertMiddlewareFuncsToRawFuncs(middlewareFuncs)
}

func ConvertMiddlewareFuncsToRawFuncs(
	middlewareFuncs []MiddlewareFunc) []func(next Handler) Handler {

	// SAFETY:
	// Unsafe here is safe until check below is failed.

	var _ func(next Handler) Handler = MiddlewareFunc(nil)
	return *(*[]func(next Handler) Handler)(unsafe.Pointer(&middlewareFuncs))
}
