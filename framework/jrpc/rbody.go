package ekaweb_jrpc

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type (

	// _JRpcRequestBody is a new http.Request's Body value. It implements
	// io.ReadCloser interface, and provides cached data to Read() op.
	// Close() invokes original underlying io.ReadCloser's Close() method.
	//
	// The main reason why this sentence is exist, is that we have to read
	// and parse original jRPC request to extract ID & Method in
	// _JRpcRequest object to perform routing and use ID later in jRPC response.
	// So, the real payload is stored here to provide it to user.
	_JRpcRequestBody struct {
		Data json.RawMessage
		Orig io.Closer
	}
)

// parseJRpcRequest is kinda entry point to jRPC request. It parses given's
// http.Request's Body, extracts jRPC request data from it, returns its
// `id` and `method` as 1st and 2nd returned arguments and saves JSON content
// of `params` field as the new http.Request's Body field using _JRpcRequestBody
// as a new io.ReadCloser.
//
// Returns an error as 3rd arg if any error is occurred
// during initialization of jRPC context.
//
// WARNING!
// No checks whether `method` is empty or not! So, logically even if nil
// as an error is returned, jRPC context may still be invalid.
func parseJRpcRequest(
	r *http.Request) (id json.RawMessage, method string, err error) {

	var req _JRpcRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		err = fmt.Errorf("%w: %s", ErrRequestMalformed, err.Error())
		return nil, "", err
	}

	r.Body = &_JRpcRequestBody{Data: req.Params, Orig: r.Body}
	return req.ID, req.Method, nil
}

// Read implements io.Reader interface and copies cached jRPC params from
// original io.ReadCloser http.Request's Body field.
func (b *_JRpcRequestBody) Read(p []byte) (n int, err error) {
	switch {

	case len(b.Data) == 0:
		err = io.EOF

	default:
		n = copy(p, b.Data)
		b.Data = b.Data[n:]
	}

	return n, err
}

// Close implements io.Closer interface and calls the same method of underlying
// object that is original io.ReadCloser from http.Request's Body field.
func (b *_JRpcRequestBody) Close() error {
	return b.Orig.Close()
}
