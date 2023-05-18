package ekaweb_jrpc

import (
	"errors"
	"io"
	"net/http"

	"github.com/goccy/go-json"

	"github.com/inaneverb/ekacore/ekaunsafe/v4"
	"github.com/inaneverb/ekaweb"
	"github.com/inaneverb/ekaweb/private"
)

type jRpcRouter struct {
	routes      map[string]ekaweb.Handler
	middlewares []ekaweb.Middleware
}

var (
	ErrMethodNotFound = errors.New("jRPC: Method not found")
	ErrNoData         = errors.New("jRPC: No data")
)

func (j *jRpcRouter) Reg(path string, middlewaresAndHandler ...any) ekaweb.RouterSimple {
	var checkError = ekaweb_private.NewCheckErrorMiddleware()
	var middlewares, handler = ekaweb_private.BuildHandlerOut(middlewaresAndHandler, checkError, false)
	handler = ekaweb_private.MergeMiddlewares(middlewares, handler)
	j.routes[path] = handler
	return j
}

func (j *jRpcRouter) Build() ekaweb.Handler {

	var routes = make(map[string]ekaweb.Handler)
	for k, v := range j.routes {
		routes[k] = v
	}

	var handler ekaweb.Handler = ekaweb.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ctx = r.Context()

		markJRpc(ctx)

		var req jRpcRequest
		var err = json.NewDecoder(r.Body).Decode(&req)

		switch {
		case errors.Is(err, io.EOF):
			ekaweb_private.UkvsInsertUserError(ctx, ErrNoData)
			return

		case err != nil:
			ekaweb_private.UkvsInsertUserError(ctx, err)
			return
		}

		requestIdSave(ctx, req.ID)
		requestMethodSave(ctx, req.Method)

		var handler = routes[req.Method]
		if handler == nil {
			ekaweb_private.UkvsInsertUserError(ctx, ErrMethodNotFound)
			return
		}

		r.Body = newStolenBody(r.Body, req.Params)
		handler.ServeHTTP(w, r)
	})

	if len(j.middlewares) > 0 {
		handler = ekaweb_private.MergeMiddlewares(j.middlewares, handler)
	}

	return handler
}

func NewRouter(options ...ekaweb.RouterOption) ekaweb.RouterSimple {

	var r jRpcRouter
	r.routes = make(map[string]ekaweb.Handler)
	r.middlewares = []ekaweb.Middleware{ekaweb.MiddlewareFunc(encDecMiddleware)}

	var doCoreInit = true

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
			r.middlewares = append(r.middlewares, middleware)

		case *ekaweb_private.RouterOptionErrorHandler:
			if option.Handler != nil {
				var middleware = ekaweb_private.NewErrorHandlerMiddleware(option.Handler)
				r.middlewares = append(r.middlewares, middleware)
			}
		}
	}

	if doCoreInit {
		var middleware = ekaweb_private.NewCoreInitializerMiddleware()
		r.middlewares = append([]ekaweb.Middleware{middleware}, r.middlewares...)
	}

	return &r
}

var _ ekaweb.RouterSimple = (*jRpcRouter)(nil)
