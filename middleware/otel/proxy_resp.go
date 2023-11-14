package ekaweb_otel

import (
	"bytes"
	"net/http"
)

// _ProxyResp is a special type that wraps original http.ResponseWriter.
// It caches all data that is written to the original http.ResponseWriter,
// allowing to use it later as an OpenTelemetry attributes.
type _ProxyResp struct {
	orig http.ResponseWriter
	buf  *bytes.Buffer
}

func (r _ProxyResp) Header() http.Header {
	return r.orig.Header()
}

func (r _ProxyResp) Write(bytes []byte) (int, error) {
	var n, err = r.orig.Write(bytes)
	_, _ = r.buf.Write(bytes)
	return n, err
}

func (r _ProxyResp) WriteHeader(statusCode int) {
	r.orig.WriteHeader(statusCode)
}

////////////////////////////////////////////////////////////////////////////////
///// PRIVATE FUNCTIONS ////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

// wrapResponse returns _ProxyResp as http.ResponseWriter, wrapping original
// http.ResponseWriter, caching all data that is written to the HTTP response.
func wrapResponse(orig http.ResponseWriter, buf *bytes.Buffer) http.ResponseWriter {
	return _ProxyResp{orig, buf}
}
