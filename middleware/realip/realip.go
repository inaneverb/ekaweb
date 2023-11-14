package ekaweb_realip

import (
	"net"
	"net/http"
	"strings"

	"github.com/inaneverb/ekaweb/v2"
)

type (
	middleware struct {
		next              ekaweb.Handler
		trustedCIDRs      []net.IPNet
		additionalHeaders []string
		realIpSaver       IPSaver
	}

	IPSaver = func(r *http.Request, ip string)
)

// gGenericHeaders contains default
var gGenericHeaders = []string{ekaweb.HeaderXForwardedFor}

// ClientIP tries to extract a real client IP addr, using (at least)
// http.Request's RemoteAddr field or extracting IP from proxy headers,
// like X-Forwarded-For, X-Real-IP, etc. Using these headers is allowed
// only if IP in these headers is from trusted CIDRs.
// Moreover, you can provide additional headers, which will be scanned
// on the same manner.
// Inspired by: https://github.com/gin-gonic/gin/blob/v1.9.1/context.go#L771

func New(trustedCIDRs []net.IPNet, options ...Option) ekaweb.Middleware {

	var m = middleware{nil, trustedCIDRs, nil, defaultIPSaver}

	for _, option := range options {
		if option != nil {
			option(&m)
		}
	}

	return &m
}

func (m middleware) Callback(next ekaweb.Handler) ekaweb.Handler {
	m.next = next
	return ekaweb.HandlerFunc(m.serveHTTP)
}

func (m middleware) CheckErrorBefore() bool { return false }

////////////////////////////////////////////////////////////////////////////////
///// PRIVATE METHODS //////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func defaultIPSaver(r *http.Request, ip string) { r.RemoteAddr = ip }

func (m middleware) serveHTTP(w http.ResponseWriter, r *http.Request) {
	r.RemoteAddr = m.extractIP(r)
	m.next.ServeHTTP(w, r)
}

func (m middleware) extractIP(r *http.Request) string {
	var addr string
	if addr = clientIP(r, m.additionalHeaders, m.trustedCIDRs); addr != "" {
		return addr
	}
	if addr = clientIP(r, gGenericHeaders, m.trustedCIDRs); addr != "" {
		return addr
	}
	if addr, _, _ = net.SplitHostPort(r.RemoteAddr); addr != "" {
		return addr
	}
	return r.RemoteAddr
}

// clientIP is a helper for ClientIP(). It performs whole process only for
// given 'headers', returning a calculated IP addr.
func clientIP(r *http.Request, headers []string, trustedCIDRs []net.IPNet) string {

	for _, header := range headers {
		for addr := r.Header.Get(header); addr != ""; {

			var i = strings.LastIndexByte(addr, ',') + 1
			var x = i

			for ; i < len(addr) && addr[i] <= ' '; i++ {
			}

			if i < len(addr) {
				if len(trustedCIDRs) == 0 {
					return addr[i:]
				}

				var ip = net.ParseIP(addr[i:])
				if ip == nil {
					return ""
				}

				var found bool
				for _, trustedCIDR := range trustedCIDRs {
					if found = trustedCIDR.Contains(ip); found {
						break
					}
				}

				if !found {
					return addr[i:]
				}
			}

			if x == 0 {
				addr = ""
			} else {
				addr = addr[:x-1]
			}
		}
	}

	return ""
}
