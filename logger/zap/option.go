package ekaweb_zap

import (
	"go.uber.org/zap"
)

// Option is a callback that allows to modify middleware under its construction.
type Option = func(m *middleware)

// WithLogger returns a new Option, that allows you to specify, what zap.Logger
// object shall be used to write log entries.
func WithLogger(log *zap.Logger) Option {
	return func(m *middleware) {
		if log != nil {
			m.log = log
		}
	}
}

// WithStringExtractorOnSuccess returns a new Option, that allows you to
// register a new string extractor for success finished HTTP requests.
//
// You may use other functions to register diff extractors for diff behavior.
// E.g: Specify an extractor for only failed HTTP requests (those, which are
// ended with an attached error).
// They're up to you: WithStringExtractorOnFail(),
// WithAnyExtractorOnSuccess(), WithAnyExtractorOnFail().
//
// Read more: CallbackExtractorString docs.
func WithStringExtractorOnSuccess(key string, callback CallbackExtractorString) Option {
	return withStringExtractor(key, callback, true)
}

// WithStringExtractorOnFail returns an Option that registers string extractor
// for failed HTTP requests. Read more: WithStringExtractorOnSuccess().
func WithStringExtractorOnFail(key string, callback CallbackExtractorString) Option {
	return withStringExtractor(key, callback, true)
}

// WithAnyExtractorOnSuccess returns an Option that registers "any" extractor
// for succeeded HTTP requests. Read more: WithStringExtractorOnSuccess().
func WithAnyExtractorOnSuccess(key string, callback CallbackExtractorAny) Option {
	return withAnyExtractor(key, callback, true)
}

// WithAnyExtractorOnFail returns an Option that registers "any" extractor
// for failed HTTP requests. Read more: WithStringExtractorOnSuccess().
func WithAnyExtractorOnFail(key string, callback CallbackExtractorAny) Option {
	return withAnyExtractor(key, callback, false)
}

// WithStringExtractor returns an Option that registers string extractor
// for any HTTP requests. Read more: WithStringExtractorOnSuccess().
func WithStringExtractor(key string, callback CallbackExtractorString) Option {
	return withMany(
		withStringExtractor(key, callback, true),
		withStringExtractor(key, callback, false),
	)
}

// WithAnyExtractor returns an Option that registers any extractor
// for any HTTP requests. Read more: WithStringExtractorOnSuccess().
func WithAnyExtractor(key string, callback CallbackExtractorAny) Option {
	return withMany(
		withAnyExtractor(key, callback, true),
		withAnyExtractor(key, callback, false),
	)
}

////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func withStringExtractor(
	key string, callback CallbackExtractorString, onSuccess bool) Option {

	return func(m *middleware) {
		var extractors = &m.fromOptions.extStrOnSuccess
		if !onSuccess {
			extractors = &m.fromOptions.extStrOnFail
		}
		if key != "" && callback != nil {
			*extractors = append(*extractors, newStringExtractor(key, callback))
		}
	}
}

func withAnyExtractor(
	key string, callback CallbackExtractorAny, onSuccess bool) Option {

	return func(m *middleware) {
		var extractors = &m.fromOptions.extAnyOnSuccess
		if !onSuccess {
			extractors = &m.fromOptions.extAnyOnFail
		}
		if key != "" && callback != nil {
			*extractors = append(*extractors, newAnyExtractor(key, callback))
		}
	}
}

func withMany(opt ...Option) Option {
	return func(m *middleware) {
		for i, n := 0, len(opt); i < n; i++ {
			opt[i](m)
		}
	}
}
