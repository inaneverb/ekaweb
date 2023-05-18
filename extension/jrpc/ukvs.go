package ekaweb_jrpc

import (
	"context"
	"net/http"

	"github.com/inaneverb/ekaweb/private"
)

type _ukvsJRpcEnabled struct{}
type _ukvsJRpcRequestId struct{}
type _ukvsJRpcRequestMethod struct{}

func markJRpc(ctx context.Context) {
	ekaweb_private.UkvsInsert(ctx, _ukvsJRpcEnabled{}, _ukvsJRpcEnabled{})
}

func IsJRPCByContext(ctx context.Context) bool {
	var _, found = ekaweb_private.UkvsLookup(ctx, _ukvsJRpcEnabled{})
	return found
}

func IsJRPC(r *http.Request) bool {
	return IsJRPCByContext(r.Context())
}

func requestIdSave(ctx context.Context, id string) {
	if id != "" {
		ekaweb_private.UkvsInsert(ctx, _ukvsJRpcRequestId{}, id)
	}
}

func RequestIDByContext(ctx context.Context) string {
	return ekaweb_private.UkvsGetOrDefault(ctx, _ukvsJRpcRequestId{}, "").(string)
}

func RequestID(r *http.Request) string {
	return RequestIDByContext(r.Context())
}

func requestMethodSave(ctx context.Context, method string) {
	if method != "" {
		ekaweb_private.UkvsInsert(ctx, _ukvsJRpcRequestMethod{}, method)
	}
}

func RequestMethod(r *http.Request) string {
	return RequestMethodByContext(r.Context())
}

func RequestMethodByContext(ctx context.Context) string {
	return ekaweb_private.UkvsGetOrDefault(ctx, _ukvsJRpcRequestMethod{}, "").(string)
}
