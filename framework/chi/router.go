package ekaweb_chi

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/inaneverb/ekacore/ekaunsafe/v4"
	"github.com/inaneverb/ekaweb"
	"github.com/inaneverb/ekaweb/private"
)

type Router struct {
	origin    *chi.Mux
	manifests []childManifest
}

type childManifest struct {
	prefix string
	child  *Router
}

// _ChiMuxBindFunc is an alias for func signature of chi.Mux Get, Post, Patch,
// and other complex HTTP handler registration functions.
type _ChiMuxBindFunc func(prefix string, handler http.HandlerFunc)

// _ChiMuxBindFunc2 is an alias for func signature of chi.Mux NotFound,
// MethodNotAllowed and other easy HTTP handler registration functions.
type _ChiMuxBindFunc2 func(handler http.HandlerFunc)

////////////////////////////////////////////////////////////////////////////////
///// Router interface implementation checker //////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

var _ ekaweb.Router = (*Router)(nil)

////////////////////////////////////////////////////////////////////////////////
///// Router interface implementation //////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func (r *Router) Use(components ...any) ekaweb.Router {

	checkError := ekaweb_private.NewCheckErrorMiddleware()
	middlewares, _ := ekaweb_private.BuildHandlerOut(components, checkError, true)
	middlewaresRawFuncs := ekaweb_private.ConvertMiddlewaresToRawFuncs(middlewares)

	r.origin.Use(middlewaresRawFuncs...)
	return r
}

func (r *Router) Group(prefix string, middlewares ...any) ekaweb.Router {

	child := newEmptyRouter(chi.NewRouter())
	child.Use(middlewares...)

	r.manifests = append(r.manifests, newChildManifest(prefix, child))
	return child
}

func (r *Router) Get(path string, middlewaresAndHandler ...any) ekaweb.Router {
	return r.reg(r.origin.Get, path, middlewaresAndHandler)
}

func (r *Router) Head(path string, middlewaresAndHandler ...any) ekaweb.Router {
	return r.reg(r.origin.Head, path, middlewaresAndHandler)
}

func (r *Router) Post(path string, middlewaresAndHandler ...any) ekaweb.Router {
	return r.reg(r.origin.Post, path, middlewaresAndHandler)
}

func (r *Router) Put(path string, middlewaresAndHandler ...any) ekaweb.Router {
	return r.reg(r.origin.Put, path, middlewaresAndHandler)
}

func (r *Router) Delete(path string, middlewaresAndHandler ...any) ekaweb.Router {
	return r.reg(r.origin.Delete, path, middlewaresAndHandler)
}

func (r *Router) Connect(path string, middlewaresAndHandler ...any) ekaweb.Router {
	return r.reg(r.origin.Connect, path, middlewaresAndHandler)
}

func (r *Router) Options(path string, middlewaresAndHandler ...any) ekaweb.Router {
	return r.reg(r.origin.Options, path, middlewaresAndHandler)
}

func (r *Router) Trace(path string, middlewaresAndHandler ...any) ekaweb.Router {
	return r.reg(r.origin.Trace, path, middlewaresAndHandler)
}

func (r *Router) Patch(path string, middlewaresAndHandler ...any) ekaweb.Router {
	return r.reg(r.origin.Patch, path, middlewaresAndHandler)
}

func (r *Router) NotFound(handler any) ekaweb.Router {
	return r.reg2(r.origin.NotFound, handler)
}

func (r *Router) MethodNotAllowed(handler any) ekaweb.Router {
	return r.reg2(r.origin.MethodNotAllowed, handler)
}

////////////////////////////////////////////////////////////////////////////////
///// Router build section (as a part of Router interface) /////////////////////
////////////////////////////////////////////////////////////////////////////////

func (r *Router) prepare() {
	for i, n := 0, len(r.manifests); i < n; i++ {
		r.manifests[i].child.prepare()
		r.origin.Mount(r.manifests[i].prefix, r.manifests[i].child.origin)
	}
}

func (r *Router) Build() ekaweb.Handler {
	r.prepare()
	return r.origin
}

////////////////////////////////////////////////////////////////////////////////
///// Chi router bridge functions //////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

// reg is a main function of middlewares and handler registration for specific
// HTTP route. It takes chi.Mux registration function and the set of parameters
// using which a new HTTP route should be created.
func (r *Router) reg(
	originCallback _ChiMuxBindFunc,
	prefix string, components []any) ekaweb.Router {

	prefix = strings.TrimSpace(prefix)
	if prefix == "" || prefix[0] != '/' {
		return r
	}

	componentsBak := components
	components = make([]any, 0, len(componentsBak)+1)

	// ------------------------------------------------------------------ //
	// WARNING! DO NOT REMOVE THIS MIDDLEWARE FROM THIS METHOD
	// UNLESS YOU FULLY UNDERSTAND WHAT ARE YOU DOING.
	// Possible bugs otherwise: Empty chi route context, empty URL variables.
	// This middleware shouldn't be global, it must be path middleware's part.
	components = append(components, NewCleanPathAndVariablesMiddleware())
	// ------------------------------------------------------------------ //

	components = append(components, componentsBak...)

	checkError := ekaweb_private.NewCheckErrorMiddleware()

	middlewares, handler := ekaweb_private.BuildHandlerOut(components, checkError, false)
	handler = ekaweb_private.MergeMiddlewares(middlewares, handler)

	originCallback(prefix, handler.ServeHTTP)
	return r
}

func (r *Router) reg2(
	originCallback _ChiMuxBindFunc2, handler any) ekaweb.Router {

	asHandler := ekaweb_private.AsHandler(handler)
	if asHandler != nil {
		originCallback(asHandler.ServeHTTP)
	}

	return r
}

////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func newEmptyRouter(origin *chi.Mux) *Router {
	return &Router{origin, nil}
}

func newChildManifest(prefix string, child *Router) childManifest {
	return childManifest{prefix, child}
}

////////////////////////////////////////////////////////////////////////////////
///// Router constructors //////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func NewRouter(options ...ekaweb.RouterOption) ekaweb.Router {
	r := newEmptyRouter(chi.NewRouter())

	middlewares := make([]ekaweb.Middleware, 0, 10)
	middlewares = append(middlewares)

	var doCoreInit = true
	var customResponseHeaders = http.Header{}

	for i, n := 0, len(options); i < n; i++ {
		if ekaunsafe.UnpackInterface(options[i]).Word == nil {
			continue
		}

		switch option := options[i].(type) {

		case *ekaweb_private.RouterOptionCoreInit:
			doCoreInit = option.Enable

		case *ekaweb_private.RouterOptionCustomJSON:
			var jsonEncDec ekaweb_private.RouterOptionCustomJSON
			if option.Encoder != nil {
				jsonEncDec.Encoder = option.Encoder
			}
			if option.Decoder != nil {
				jsonEncDec.Decoder = option.Decoder
			}
			var middleware = ekaweb_private.NewJSONEncodeDecodeMiddleware(&jsonEncDec)
			middlewares = append(middlewares, middleware)

		case *ekaweb_private.RouterOptionServerName:
			if option.ServerName != "" {
				customResponseHeaders.Set(ekaweb.HeaderServer, option.ServerName)
			}

		case *ekaweb_private.RouterOptionErrorHandler:
			if option.Handler != nil {
				middleware := ekaweb_private.NewErrorHandlerMiddleware(option.Handler)
				middlewares = append(middlewares, middleware)
			}
		}
	}

	if doCoreInit {
		var mCoreInit = ekaweb_private.NewCoreInitializerMiddleware()
		var mInvalidatePath = NewInvalidatePathMiddleware()
		middlewares = append(
			[]ekaweb.Middleware{mCoreInit, mInvalidatePath}, middlewares...)
	}

	if len(customResponseHeaders) > 0 {
		middleware := ekaweb_private.NewCustomResponseHeadersMiddleware(customResponseHeaders)
		middlewares = append(middlewares, middleware)
	}

	middlewaresRawFuncs := ekaweb_private.ConvertMiddlewaresToRawFuncs(middlewares)
	r.origin.Use(middlewaresRawFuncs...)

	return r
}
