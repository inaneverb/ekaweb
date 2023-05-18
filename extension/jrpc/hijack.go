package ekaweb_jrpc

import (
	"io"

	"github.com/goccy/go-json"
)

type rcStolenBody struct {
	origin io.Closer
	msg    json.RawMessage
}

func (r *rcStolenBody) Read(p []byte) (n int, err error) {
	switch {

	case len(r.msg) == 0:
		err = io.EOF

	default:
		n = copy(p, r.msg)
		r.msg = r.msg[n:]
	}

	return n, err
}

func (r *rcStolenBody) Close() error {
	return r.origin.Close()
}

func newStolenBody(origin io.Closer, jrm json.RawMessage) io.ReadCloser {
	return &rcStolenBody{origin, jrm}
}
