package ekaweb_zerolog

import (
	"net/http"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/inaneverb/ekaweb/v2"
	"github.com/inaneverb/ekaweb/v2/private"
)

type (
	// middleware implements ekaweb.Middleware, giving an ability to create
	// a new logging middleware, backed by https://github.com/rs/zerolog .
	middleware struct {
		log zerolog.Logger

		fromOptions struct {
			extStrOnSuccess []extractorStr
			extAnyOnSuccess []extractorAny

			extStrOnFail []extractorStr
			extAnyOnFail []extractorAny
		}
	}

	// extractorStr is just a pair of http.Request's string value extraction
	// callback, and a key with which the extracted value will be logged.
	extractorStr struct {
		key string
		ext CallbackExtractorString
	}

	// extractorAny is the same as extractorString but allows to extract
	// any values from http.Request, not only strings.
	extractorAny struct {
		key string
		ext CallbackExtractorAny
	}

	// CallbackExtractorString is a function you may define, to order
	// a logging middleware log some valuable parameter from http.Request.
	//
	// You should return a value, you want to log. The key, this value
	// will be logged with, is registered by Option's constructors.
	CallbackExtractorString = func(r *http.Request) string

	// CallbackExtractorAny is the same as just CallbackExtractorString,
	// but allows you to extract any value, not only strings.
	CallbackExtractorAny = func(r *http.Request) any
)

// Callback implements ekaweb.Middleware. Returns an ekaweb.Handler,
// that does the main job, of logging HTTP request, wrapping call of 'next' one.
func (m *middleware) Callback(next ekaweb.Handler) ekaweb.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		const (
			TailHijacked = " hijacked with no error"           // 23 chars
			TailNotFound = " not found or not allowed"         // 25 chars
			TailSuccess  = " finished with no error"           // 23 chars
			TailFail     = " finished (or aborted) with error" // 33 chars
		)

		var start = time.Now()

		next.ServeHTTP(w, r)
		var err = ekaweb.ErrorGet(r)

		var execTime = time.Since(start)

		// Message will contain: HTTP method (max=7, OPTIONS, CONNECT),
		// original path (dynamic value), space between HTTP method and path,
		// a tail of result (max=32, see TailFail).
		// So: 7+1+32+N = 40+N (+8 stock).

		var mb strings.Builder
		mb.Grow(48 + len(r.URL.Path))
		mb.WriteString(r.Method)
		mb.WriteByte(' ')
		mb.WriteString(ekaweb.RoutePath(r))

		var ev *zerolog.Event

		var extString []extractorStr
		var extAny []extractorAny

		if completedWithNoError := err == nil; completedWithNoError {
			var ctx = r.Context()

			switch {
			case ekaweb_private.UkvsIsConnectionHijacked(ctx):
				mb.WriteString(TailHijacked)

			case ekaweb_private.UkvsIsPathNotFoundOrNotAllowed(ctx):
				mb.WriteString(TailNotFound)

			default:
				mb.WriteString(TailSuccess)
			}

			ev = m.log.Debug()

			extString = m.fromOptions.extStrOnSuccess
			extAny = m.fromOptions.extAnyOnSuccess

		} else {

			mb.WriteString(TailFail)

			ev = m.log.Error()

			extString = m.fromOptions.extStrOnFail
			extAny = m.fromOptions.extAnyOnFail
		}

		m.applyExtractors(r, ev, extString, extAny)

		ev.Err(err)
		ev.Str("exec_time", execTime.String())
		ev.Str("client_ip", r.RemoteAddr)

		ev.Msg(mb.String())
	})
}

// applyExtractors calls extractors providing passed context and adding
// returned values from extractors to the 'fieldsOut'.
func (_ *middleware) applyExtractors(
	r *http.Request, ev *zerolog.Event,
	extractorsString []extractorStr, extractorsAny []extractorAny,
) {
	for i, n := 0, len(extractorsString); i < n; i++ {
		key := extractorsString[i].key
		valueString := extractorsString[i].ext(r)
		if valueString != "" {
			ev.Str(key, valueString)
		}
	}

	for i, n := 0, len(extractorsAny); i < n; i++ {
		key := extractorsAny[i].key
		valueAny := extractorsAny[i].ext(r)
		if valueAny != nil {
			ev.Any(key, valueAny)
		}
	}
}

var _ ekaweb.Middleware = (*middleware)(nil)

////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

// NewMiddleware returns a new logging middleware, that allows you to log any
// incoming HTTP requests. Backed by https://github.com/uber-go/zap .
// You may pass some Option(s) to adjust its behaviour.
func NewMiddleware(opts ...Option) ekaweb.Middleware {

	var m = middleware{log: log.Logger}

	for i, n := 0, len(opts); i < n; i++ {
		if opts[i] != nil {
			opts[i](&m)
		}
	}

	return &m
}

// newStringExtractor is just extractorStr constructor.
func newStringExtractor(key string, cb CallbackExtractorString) extractorStr {
	return extractorStr{key, cb}
}

// newAnyExtractor is just extractorAny constructor.
func newAnyExtractor(key string, cb CallbackExtractorAny) extractorAny {
	return extractorAny{key, cb}
}
