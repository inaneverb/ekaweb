package ekaweb_jrpc

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/inaneverb/ekaweb/v2/private"
)

type (
	// _JRpcContext is a struct that holds info about current HTTP API route
	// jRPC context. It's request ID, occurred jRPC error (if any) and state.
	//
	// It should be initialized by the main jRPC route's controller
	// and is stored inside UKVS context.Context with _JRpcContextKey key.
	_JRpcContext struct {
		RequestID json.RawMessage
		Method    string
	}

	// _JRpcContextKey is a type that is used as a key in UKVS
	// to store _JRpcContext object.
	_JRpcContextKey struct{}
)

// UkvsGetMeta returns jRPC context's metadata. It includes jRPC `id` of request
// and jRPC `method` that is requested by the client.
// If it's called not inside jRPC context, an empty `method` is returned.
func UkvsGetMeta(r *http.Request) (id json.RawMessage, method string) {
	return UkvsGetMetaByContext(r.Context())
}

// UkvsGetMetaByContext is the same as just UkvsGetMeta() but works directly
// with context.Context, instead of wrapped http.Request.
func UkvsGetMetaByContext(
	ctx context.Context) (id json.RawMessage, method string) {

	var jCtx = getCtx(ctx)
	if jCtx == nil {
		return nil, "" // early exit: not inside jRPC context
	}

	return jCtx.RequestID, jCtx.Method
}

// getCtx extracts and returns _JRpcContext from UVKS stored inside
// given context.Context.
func getCtx(ctx context.Context) *_JRpcContext {
	var v = ekaweb_private.UkvsGet(ctx, (*_JRpcContextKey)(nil))
	var jCtx, _ = v.(*_JRpcContext) // avoid panic if nil is returned
	return jCtx
}

// setCtx saves given _JRpcContext to the UKVS, that is located inside
// provided http.Request's context.Context.
func setCtx(r *http.Request, ctx *_JRpcContext) {
	ekaweb_private.UkvsInsert(r.Context(), (*_JRpcContextKey)(nil), ctx)
}
