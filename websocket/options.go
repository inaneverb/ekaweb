package ekaweb_socket

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/inaneverb/ekaweb/v2/private"
)

// Option is a callback that allows to modify some WebSocket provider behaviour.
type Option func(o *Options)

type CallbackCheckOrigin = func(r *http.Request) bool
type CallbackIDGenerator = func(ctx context.Context) string

type Options struct {
	IDGenerator     CallbackIDGenerator
	CheckOrigin     CallbackCheckOrigin
	ErrorHandler    ekaweb_private.ErrorHandler
	ResponseHeaders http.Header
}

var defaultOptions = Options{
	IDGenerator:     defaultIDGenerator,
	CheckOrigin:     nil,
	ErrorHandler:    nil,
	ResponseHeaders: nil,
}

func defaultIDGenerator(_ context.Context) string {
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}

// PrepareOptions returns an Options object based on default and provided ones.
// The caller MUST NOT change it. Only read operations are allowed.
//
// The returned object will be a composite object of default options
// with overwriting some values by the provided ones.
// If there's no provided options, returned object will be just default options
// (w/o unnecessary allocations).
func PrepareOptions(options []Option) *Options {

	if len(options) == 0 {
		return &defaultOptions
	}

	optionsSet := defaultOptions

	for i, n := 0, len(options); i < n; i++ {
		if options[i] != nil {
			options[i](&optionsSet)
		}
	}

	return &optionsSet
}

func WithCheckOrigin(cb CallbackCheckOrigin) Option {
	return func(o *Options) {
		if cb != nil {
			o.CheckOrigin = cb
		}
	}
}

func WithIDGenerator(cb CallbackIDGenerator) Option {
	return func(o *Options) {
		if cb != nil {
			o.IDGenerator = cb
		}
	}
}

func WithErrorHandler(cb ekaweb_private.ErrorHandler) Option {
	return func(o *Options) {
		if cb != nil {
			o.ErrorHandler = cb
		}
	}
}
