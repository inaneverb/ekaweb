package ekaweb_jrpc

import (
	"context"
	"net/http"

	"github.com/goccy/go-json"

	"github.com/inaneverb/ekaweb"
	"github.com/inaneverb/ekaweb/private"
)

type jRpcRequest struct {
	ID     json.RawMessage `json:"id"`
	Method string          `json:"method"`
	Params json.RawMessage `json:"params"`
}

type jRpcResponse struct {
	Header string          `json:"jsonrpc"`
	ID     json.RawMessage `json:"id,omitempty"`
	Result any             `json:"result,omitempty"`
	Error  any             `json:"error,omitempty"`
}

func encDecMiddleware(next ekaweb.Handler) ekaweb.Handler {
	return ekaweb.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ctx = r.Context()
		var opt = genEncDecJson(ctx)
		ekaweb_private.UkvsInsertJSONEncoderDecoder(ctx, opt)
		next.ServeHTTP(w, r)
	})
}

func genEncDecJson(ctx context.Context) *ekaweb_private.RouterOptionCustomJSON {

	var opt = ekaweb_private.UkvsGetJSONEncoderDecoder(ctx)
	if opt == nil {
		opt = new(ekaweb_private.RouterOptionCustomJSON)
	}
	if opt.Encoder == nil {
		opt.Encoder = json.Marshal
	}
	if opt.Decoder == nil {
		opt.Decoder = json.Unmarshal
	}

	var orig = opt.Encoder
	opt.Encoder = func(v any) ([]byte, error) {
		var resp = jRpcResponse{"2.0", RequestIDByContext(ctx), v, nil}
		if ekaweb_private.UkvsGetUserError(ctx) != nil {
			resp.Error, resp.Result = resp.Result, nil
		}
		return orig(resp)
	}

	return opt
}

var _ ekaweb.Middleware = ekaweb.MiddlewareFunc(encDecMiddleware)
