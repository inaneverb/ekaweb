package ekaweb_jrpc

import (
	"context"
	"encoding/json"
	"io"

	"github.com/inaneverb/ekaweb/v2/private"
)

type (
	// _JRpcRequest is a representation of jRPC request. It exists only for
	// read & decode original jRPC http.Request's payload.
	_JRpcRequest struct {
		ID     json.RawMessage `json:"id"`
		Method string          `json:"method"`
		Params json.RawMessage `json:"params"`
	}

	// _JRpcResponse is a representation of jRPC response. Any object, that
	// you will send to the client will be wrapped by this sentence
	// where your data is stored as a Result or Error field.
	_JRpcResponse struct {
		Header string          `json:"jsonrpc"`
		ID     json.RawMessage `json:"id,omitempty"`
		Result any             `json:"result,omitempty"`
		Error  any             `json:"error,omitempty"`
	}

	// _JRpcEncoder is a special thing that changes internal behaviour
	// of sending jRPC response w/o requiring user to change an API.
	//
	// By jRPC router's controller it will be used as a replacement
	// for original Encoder. But it stores parts from it.
	// See Encode() for more details.
	_JRpcEncoder struct {
		ctx  context.Context              // http.Request's related context
		jCtx *_JRpcContext                // jRPC context
		w    io.Writer                    // where to write encoded data
		eg   ekaweb_private.EncoderGetter // original encode getter
	}
)

var (
	gJsonNullValue = []byte("null") // const value JSON representation
)

// newConnectedEncodeGetter creates and returns a new
// ekaweb_private.EncoderGetter that is closures given context.Context,
// jRPC request's ID and original ekaweb_private.EncoderGetter.
// They all will be used during construction of final jRPC response.
func newConnectedEncodeGetter(
	ctx context.Context, jCtx *_JRpcContext,
	origEg ekaweb_private.EncoderGetter) ekaweb_private.EncoderGetter {

	return func(w io.Writer) ekaweb_private.Encoder {
		return _JRpcEncoder{ctx, jCtx, w, origEg}
	}
}

// Encode is where the magic happens. User wants to send given 'v' as a jRPC
// response's payload, but jRPC standard requires to send special object
// as a root of response body. So, here we doing it. We're wrapping given
// by user 'v' by _JRpcResponse and sends it instead of 'v', fulfilling
// jRPC requirements.
func (je _JRpcEncoder) Encode(v any) error {
	var resp = _JRpcResponse{"2.0", je.jCtx.RequestID, v, nil}
	if je.jCtx.RequestID == nil {
		je.jCtx.RequestID = gJsonNullValue
	}

	if ekaweb_private.UkvsGetUserError(je.ctx) != nil {
		resp.Error, resp.Result = resp.Result, nil
	}

	return je.eg(je.w).Encode(&resp)
}
