package ekaweb_otel

import (
	"net/http"

	"github.com/inaneverb/ekaweb/v2"
)

// SpanNameFromRoutePath returns span name for OpenTelemetry,
// based on current API route.
// It has the signature that suits OpenTelemetry option's functions,
// allowing you to pass this function to some option registration callbacks.
func SpanNameFromRoutePath(op string, r *http.Request) string {
	var path = ekaweb.RoutePath(r)
	if op == "" {
		return path
	}
	return op + " " + path
}
