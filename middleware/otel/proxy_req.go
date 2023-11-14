package ekaweb_otel

import (
	"bytes"
	"errors"
	"io"
)

// _ProxyReq is a special type that wraps original io.ReadCloser,
// assuming that it's from http.Request's Body field.
// It caches all data during read operations, allowing to use it later
// as an OpenTelemetry attributes.
type _ProxyReq struct {
	orig io.ReadCloser
	buf  *bytes.Buffer
}

func (r _ProxyReq) Read(p []byte) (n int, err error) {
	if n, err = r.orig.Read(p); err != nil && !errors.Is(err, io.EOF) {
		return n, err
	}
	_, _ = r.buf.Write(p[:n])
	return n, err
}

func (r _ProxyReq) Close() error {
	return r.orig.Close()
}

////////////////////////////////////////////////////////////////////////////////
///// PRIVATE FUNCTIONS ////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

// wrapRequestBody returns _ProxyReq as io.ReadCloser, wrapping original
// io.ReadCloser, assuming that it's a body of from http.Request,
// caching all read that during read operations (no reads - no data).
func wrapRequestBody(orig io.ReadCloser, buf *bytes.Buffer) io.ReadCloser {
	return _ProxyReq{orig, buf}
}
