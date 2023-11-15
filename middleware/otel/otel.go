package ekaweb_otel

import (
	"bytes"
	"net/http"

	"github.com/inaneverb/ekaweb/v2"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

// Inspired by:
// - https://bunrouter.uptrace.dev/guide/golang-http-performance.html#bunrouter-opentelemetry-instrumentation
// - https://github.com/uptrace/bunrouter/blob/v1.0.20/extra/bunrouterotel/bunrouterotel.go

type middleware struct {
	addClientIP        bool
	addRequestHeaders  bool
	addRequestBody     bool
	addResponseHeaders bool
	addResponseBody    bool
}

// New creates a new OpenTelemetry middleware.
func New(opts ...Option) ekaweb.Middleware {
	var m middleware

	for i, n := 0, len(opts); i < n; i++ {
		if opts[i] != nil {
			opts[i](&m)
		}
	}

	return &m
}

func (m *middleware) Callback(next ekaweb.Handler) ekaweb.Handler {

	const (
		AttributeKeyHeaders = "Headers"
		AttributeKeyBody    = "Body"
	)

	return ekaweb.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var span = trace.SpanFromContext(r.Context())
		if !span.IsRecording() {
			next.ServeHTTP(w, r)
			return
		}

		var attrs = make([]attribute.KeyValue, 0, 4)
		var routePath = ekaweb.RoutePath(r)

		span.SetName(r.Method + " " + routePath)
		attrs = append(attrs, semconv.HTTPRouteKey.String(routePath))

		if r.URL != nil {
			attrs = append(attrs, semconv.HTTPTargetKey.String(r.URL.RequestURI()))
		} else {
			attrs = append(attrs, semconv.HTTPTargetKey.String(r.RequestURI))
		}

		if m.addClientIP {
			attrs = append(attrs, semconv.HTTPClientIPKey.String(r.RemoteAddr))
		}

		var bufReq, bufResp _ProxyBuf

		if m.addRequestHeaders || m.addRequestBody {
			bufReq = acquireBuffer()
			//defer releaseBuffer(bufReq) // w/o defer placed below

			if m.addRequestBody {
				r.Body = wrapRequestBody(r.Body, bufReq.body) // <- R wrapping
			}
			if m.addRequestHeaders {
				flushHeaders(bufReq, r.Header)
			}
		}

		if m.addResponseHeaders || m.addResponseBody {
			bufResp = acquireBuffer()
			//defer releaseBuffer(bufResp) // w/o defer placed below

			if m.addResponseBody {
				w = wrapResponse(w, bufResp.body) // <- W wrapping
			}
		}

		span.SetAttributes(attrs...)
		next.ServeHTTP(w, r)

		// WARNING!
		// We should flush headers only after executing inner handler,
		// because headers are stored inside http.ResponseWriter and they
		// will be filled by the inner handler.

		if m.addResponseHeaders {
			flushHeaders(bufResp, w.Header())
		}

		//attrs = attrs[:0]

		// In all these calls like `buf.<kind>.String()` makes a copy of []byte,
		// returning it as a string. Thus, it's safe to return these buffers
		// to the pool later.

		// Add request details as an event (log) to the span, if required.

		if m.addRequestHeaders || m.addRequestBody {
			attrs = attrs[:0]

			if m.addRequestHeaders {
				attrs = append(attrs, attribute.String(
					AttributeKeyHeaders, bufReq.headers.String()))
			}
			if m.addRequestBody {
				attrs = append(attrs, attribute.String(
					AttributeKeyBody, bufReq.body.String()))
			}
			span.AddEvent("Request details", trace.WithAttributes(attrs...))
		}

		// Add response details as an event (log) to the span, if required.

		if m.addResponseHeaders || m.addResponseBody {
			attrs = attrs[:0]

			if m.addResponseHeaders {
				attrs = append(attrs, attribute.String(
					AttributeKeyHeaders, bufResp.headers.String()))
			}
			if m.addResponseBody {
				attrs = append(attrs, attribute.String(
					AttributeKeyBody, bufResp.body.String()))
			}
			span.AddEvent("Response details", trace.WithAttributes(attrs...))
		}

		// To avoid multiple checks, we can just check whether buffer pointer
		// is not nil (thus buffers were allocated).

		if bufReq.body != nil {
			releaseBuffer(bufReq)
		}
		if bufResp.body != nil {
			releaseBuffer(bufResp)
		}

		//if len(attrs) > 0 {
		//	span.SetAttributes(attrs...)
		//}

		if err := ekaweb.ErrorGet(r); err != nil {
			span.SetStatus(codes.Error, err.Error())
		}
	})
}

func (m *middleware) CheckErrorBefore() bool { return false }

////////////////////////////////////////////////////////////////////////////////
///// PRIVATE METHODS //////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func (_ *middleware) addReqRespItem(
	to *[]attribute.KeyValue, key string, buf *bytes.Buffer, check bool) {

	// buf.String() here makes a copy of underlying that.
	// Thus, it's safe to reuse this buffer later.

	if check {
		*to = append(*to, attribute.String(key, buf.String()))
	}
}
