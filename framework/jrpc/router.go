package ekaweb_jrpc

import (
	"encoding/json"
	"net/http"

	"github.com/inaneverb/ekacore/ekaunsafe/v4"
	"github.com/inaneverb/ekaweb/v2"
	"github.com/inaneverb/ekaweb/v2/middleware"
	"github.com/inaneverb/ekaweb/v2/private"
)

// _JRpcRouter is a main thing of jRPC routing. It implements
// ekaweb.RouterSimple, thus allowing you to register jRPC methods
// and their handlers. After all handlers are registered you have to call
// Build() to get main final handler that will serve incoming jRPC requests.
type _JRpcRouter struct {
	routes      map[string]ekaweb.Handler
	middlewares []ekaweb.Middleware
}

// Reg registers new jRPC route with given 'method' and a set of middlewares
// and handler(s). They will be invoked when jRPC request with associated
// method is received. Returns current jRPC router.
//
// WARNING! If the route with the same method is already registered,
// it will be overwritten.
func (j *_JRpcRouter) Reg(
	method string, middlewaresAndHandler ...any) ekaweb.RouterSimple {

	var checkError = ekaweb_private.NewCheckErrorMiddleware()
	var middlewares, handler = ekaweb_private.BuildHandlerOut(
		middlewaresAndHandler, checkError, false)

	handler = ekaweb_private.MergeMiddlewares(middlewares, handler)
	j.routes[method] = handler
	return j
}

// Build builds final handler that you should register in your server,
// or other router.
func (j *_JRpcRouter) Build() ekaweb.Handler {

	var routes = make(map[string]ekaweb.Handler)
	for k, v := range j.routes {
		routes[k] = v
	}

	var h = ekaweb.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ctx = r.Context()

		// Step 1.
		// Initialize jRPC context completely, even before doing any checks.
		// It will get a caller additional
		// information about jRPC context during error handling.

		var jCtx _JRpcContext
		ekaweb_private.UkvsInsert(r.Context(), (*_JRpcContextKey)(nil), &jCtx)

		var codec = ekaweb_private.UkvsGetCodec(ctx)
		codec.EncoderGetter =
			newConnectedEncodeGetter(ctx, &jCtx, codec.EncoderGetter)

		ekaweb_private.UkvsInsertCodec(ctx, codec)

		// Step 2.
		// Try to parse body of incoming request, considering it's jRPC request.
		// It extracts jRPC ID, jRPC method, jRPC params. Params are saved back
		// as the real payload of jRPC request, but id & method are returned.

		var err error
		if jCtx.RequestID, jCtx.Method, err = parseJRpcRequest(r); err != nil {
			ekaweb_private.UkvsInsertUserError(ctx, err)
			return // early exit: malformed jRPC request
		}

		// Step 3.
		// Check if jRPC method is provided. Lookup for requested jRPC method.

		if jCtx.Method == "" {
			ekaweb_private.UkvsInsertUserError(ctx, ErrRequestMalformed)
			return // early exit: no method is provided
		}

		var h = routes[jCtx.Method]
		if h == nil {
			ekaweb_private.UkvsInsertUserError(ctx, ErrMethodNotRegistered)
			return // early exit: no such jRPC method found
		}

		// Step 4.
		// Execute handler.

		h.ServeHTTP(w, r)
	})

	return ekaweb_private.MergeMiddlewares(j.middlewares, h)
}

// NewRouter initializes and returns a new jRPC router.
//
// WARNING! IF YOU PLAN TO USE RETURNED ROUTER AS A SUB ROUTER OF OTHER
// ROUTER, THAT IS ALSO PART OF EKAWEB, YOU SHOULD PASS OPTION THAT DISABLES
// CORE INITIALIZATION: ekaweb.WithCoreInit(false).
func NewRouter(options ...ekaweb.RouterOption) ekaweb.RouterSimple {

	var r _JRpcRouter
	r.routes = make(map[string]ekaweb.Handler)

	var doCoreInit = true
	var customResponseHeaders = http.Header{}

	var ukvsManager *ekaweb_private.UkvsManager
	var optCodec *ekaweb_private.RouterOptionCodec

	for i, n := 0, len(options); i < n; i++ {
		if ekaunsafe.UnpackInterface(options[i]).Word == nil {
			continue
		}

		switch option := options[i].(type) {

		case *ekaweb_private.RouterOptionCoreInit:
			doCoreInit = option.Enable

		case *ekaweb_private.RouterOptionCodec:
			optCodec = option

		case *ekaweb_private.RouterOptionServerName:
			if option.ServerName != "" {
				customResponseHeaders.Set(ekaweb.HeaderServer, option.ServerName)
			}

		case *ekaweb_private.RouterOptionErrorHandler:
			if option.Handler != nil {
				var middleware = ekaweb_private.NewErrorHandlerMiddleware(option.Handler)
				r.middlewares = append(r.middlewares, middleware)
			}
		}
	}

	if doCoreInit {
		if ukvsManager == nil {
			if optCodec == nil {
				type T = ekaweb_private.RouterOptionCodec
				optCodec = ekaweb.WithCodec(json.NewEncoder, json.NewDecoder).(*T)
			}
			var g = ekaweb_private.NewUkvsMapGeneratorSlice()
			ukvsManager = ekaweb_private.NewUkvsManager(g, *optCodec)
		}
		var mCoreInit = ekaweb_private.NewUkvsManagerMiddleware(ukvsManager)
		r.middlewares = append([]ekaweb.Middleware{mCoreInit}, r.middlewares...)
	}

	if len(customResponseHeaders) > 0 {
		var mCustomHeaders = ekaweb_middleware.CustomHeaders(customResponseHeaders)
		r.middlewares = append(r.middlewares, mCustomHeaders)
	}

	return &r
}

var _ ekaweb.RouterSimple = (*_JRpcRouter)(nil)
