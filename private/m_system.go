package ekaweb_private

import (
	"context"
	"net/http"
)

// MiddlewareCoreInitialize is an HTTP middleware that prepares
// http.Request's context for being used by the user.
//
// It injects a thread-safety lock-free key-value storage inside
// (to be able get-n-store user values without context.WithValue() calls)
// and ensures context.Context will be cancelled when the middleware is done.
//
// This is a part of automatically enabled middlewares,
// that is shadowy integrated at the router level to absolutely each handler.
type MiddlewareCoreInitialize struct{}

func (*MiddlewareCoreInitialize) CheckErrorBefore() bool {
	return false
}

func (*MiddlewareCoreInitialize) Callback(next Handler) Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancelFunc := context.WithCancel(UkvsInit(r.Context()))
		defer cancelFunc()
		defer UkvsDestroy(ctx)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func NewCoreInitializerMiddleware() Middleware {
	return &MiddlewareCoreInitialize{}
}

////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

type MiddlewareJSONEncodeDecode struct {
	Option *RouterOptionCustomJSON
}

func (*MiddlewareJSONEncodeDecode) CheckBeforeError() bool {
	return false
}

func (m *MiddlewareJSONEncodeDecode) Callback(next Handler) Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		UkvsInsertJSONEncoderDecoder(r.Context(), m.Option)
		next.ServeHTTP(w, r)
	})
}

func NewJSONEncodeDecodeMiddleware(option *RouterOptionCustomJSON) Middleware {
	return &MiddlewareJSONEncodeDecode{option}
}

////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

type MiddlewareCustomResponseHeaders struct {
	CustomHeaders http.Header
}

func (*MiddlewareCustomResponseHeaders) CheckErrorBefore() bool {
	return false
}

func (m *MiddlewareCustomResponseHeaders) Callback(next Handler) Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		UkvsInsertResponseCustomHeaders(r.Context(), m.CustomHeaders)
		next.ServeHTTP(w, r)
	})
}

func NewCustomResponseHeadersMiddleware(headers http.Header) Middleware {
	return &MiddlewareCustomResponseHeaders{headers}
}
