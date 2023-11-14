package ekaweb_otel

import (
	"bytes"
	"net/http"
	"sync"
)

// _ProxyBuf contains buffers where HTTP response headers & body
// will be stored.
type _ProxyBuf struct {
	headers *bytes.Buffer
	body    *bytes.Buffer
}

// poolProxyBuf is a sync pool where all buffers of HTTP request/response
// are stored. Their initial sizes are 4KB.
var poolProxyBuf = sync.Pool{
	New: func() any {
		return _ProxyBuf{
			headers: bytes.NewBuffer(make([]byte, 0, 4096)),
			body:    bytes.NewBuffer(make([]byte, 0, 4096)),
		}
	},
}

// wrapResponse returns _ProxyResp as http.ResponseWriter, wrapping original
// http.ResponseWriter, caching all data that is written to the HTTP response.
// It uses pool to obtain all required buffers.
// An opposite for releaseWrappedResponse().

// acquireBuffer returns _ProxyBuf from its internal pool. Initial sizes
// of all internal buffers are 4KB. An opposite for releaseBuffer().
func acquireBuffer() _ProxyBuf {
	return poolProxyBuf.Get().(_ProxyBuf)
}

// flushHeaders flushes given http.Header to provided _ProxyBuf.
func flushHeaders(to _ProxyBuf, from http.Header) {
	_ = from.Write(to.headers)
}

// releaseBuffer releases (frees) all allocated buffers inside given _ProxyBuf
// preparing it to be reused. An opposite for acquireBuffer().
// WARNING! YOU SHOULD NOT USE _ProxyBuf AFTER CALLING THIS FUNCTION!
func releaseBuffer(p _ProxyBuf) {
	p.body.Reset()
	p.headers.Reset()
	poolProxyBuf.Put(p)
}
