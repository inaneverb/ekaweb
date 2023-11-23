package ekaweb_otel

// Option is a callback that allows to modify Middleware under its construction.
type Option func(m *middleware)

// WithClientIP adds a http.Request's RemoteAddr to the OpenTelemetry
// attributes as a client's IP address.
//
// NOTE.
// Maybe you also want to register middleware that will overwrite RemoteAddr
// field in http.Request from some headers if the request went through proxies.
func WithClientIP(enable ...bool) Option {
	return func(m *middleware) {
		m.addClientIP = !(len(enable) > 0 && !enable[0])
	}
}

// WithRequestData adds http.Request's data (headers, body) to the OpenTelemetry
// attributes if correspondent boolean flags are set.
// WARNING! It may significantly increase the time of processing HTTP request.
func WithRequestData(addHeaders, addBody bool) Option {
	return func(m *middleware) {
		m.addRequestHeaders, m.addRequestBody = addHeaders, addBody
	}
}

// WithResponseData adds HTTP response's data (headers, body)
// to the OpenTelemetry attributes if correspondent boolean flags are set.
// WARNING! It may significantly increase the time of processing HTTP request.
func WithResponseData(addHeaders, addBody bool) Option {
	return func(m *middleware) {
		m.addResponseHeaders, m.addResponseBody = addHeaders, addBody
	}
}

// WithRecheckMethodPath adds additional recheck that HTTP method or/and
// path was changed during call of next ekaweb.Handler and if it so,
// changes span name to the new one.
func WithRecheckMethodPath(enable ...bool) Option {
	return func(m *middleware) {
		m.recheckMethodPath = !(len(enable) > 0 && !enable[0])
	}
}
