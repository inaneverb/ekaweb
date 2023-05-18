package ekaweb_zap

import (
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/inaneverb/ekaweb"
	"github.com/inaneverb/ekaweb/private"
)

type (
	// middleware implements ekaweb.Middleware, giving an ability to create
	// a new logging middleware, backed by https://github.com/uber-go/zap .
	middleware struct {
		log *zap.Logger

		fromOptions struct {
			extStrOnSuccess []extractorStr
			extAnyOnSuccess []extractorAny

			extStrOnFail []extractorStr
			extAnyOnFail []extractorAny
		}

		afterBuilding struct {
			fieldsNumberOnSuccess int
			fieldsNumberOnFail    int
		}
	}

	// extractorStr is just a pair of http.Request's string value extraction
	// callback, and a key with which the extracted value will be logged.
	extractorStr struct {
		key string
		ext CallbackExtractorString
	}

	// extractorAny is the same as extractorStr but allows to extract
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

// CheckErrorBefore implements ekaweb.MiddlewareExtended. Always returns false.
// Thus, it requests to skip http.Request's error checking during middlewares &
// handler call stack unwinding.
//
// Speaking simply: If log middleware is after A callback, then during theirs
// combine process (zip them to only 1 handler), the error check middleware
// won't be added between A & logging one.
func (m *middleware) CheckErrorBefore() bool {
	return false
}

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

		var log func(logger *zap.Logger, message string, fields ...zap.Field)
		var fields []zap.Field

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

			log = (*zap.Logger).Debug
			fields = make([]zap.Field, 0, m.afterBuilding.fieldsNumberOnSuccess)

			extString = m.fromOptions.extStrOnSuccess
			extAny = m.fromOptions.extAnyOnSuccess

		} else {

			mb.WriteString(TailFail)

			log = (*zap.Logger).Error
			fields = make([]zap.Field, 0, m.afterBuilding.fieldsNumberOnFail)

			extString = m.fromOptions.extStrOnFail
			extAny = m.fromOptions.extAnyOnFail
		}

		m.applyExtractors(r, &fields, extString, extAny)

		fields = append(fields, zap.Error(err))
		fields = append(fields, zap.Duration("exec_time", execTime))
		fields = append(fields, zap.String("client_ip", r.RemoteAddr))

		log(m.log, mb.String(), fields...)
	})
}

// applyExtractors calls extractors providing passed context and adding
// returned values from extractors to the 'fieldsOut'.
func (_ *middleware) applyExtractors(
	r *http.Request, fieldsOut *[]zap.Field,
	extractorsString []extractorStr, extractorsAny []extractorAny,
) {
	for i, n := 0, len(extractorsString); i < n; i++ {
		key := extractorsString[i].key
		valueString := extractorsString[i].ext(r)
		if valueString != "" {
			*fieldsOut = append(*fieldsOut, zap.String(key, valueString))
		}
	}

	for i, n := 0, len(extractorsAny); i < n; i++ {
		key := extractorsAny[i].key
		valueAny := extractorsAny[i].ext(r)
		if valueAny != nil {
			*fieldsOut = append(*fieldsOut, zap.Any(key, valueAny))
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

	const AttachedFieldsMinimum = 10

	var m = middleware{}
	m.log = zap.L()

	for i, n := 0, len(opts); i < n; i++ {
		if opts[i] != nil {
			opts[i](&m)
		}
	}

	var n01 = len(m.fromOptions.extStrOnSuccess)
	var n02 = len(m.fromOptions.extAnyOnSuccess)

	var n11 = len(m.fromOptions.extStrOnFail)
	var n12 = len(m.fromOptions.extAnyOnFail)

	m.afterBuilding.fieldsNumberOnSuccess = AttachedFieldsMinimum + n01 + n02
	m.afterBuilding.fieldsNumberOnFail = AttachedFieldsMinimum + n11 + n12

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
