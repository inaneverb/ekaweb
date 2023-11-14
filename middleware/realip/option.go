package ekaweb_realip

import (
	"slices"
)

// Option is a callback that allows to modify Middleware under its construction.
type Option func(m *middleware)

func WithAdditionalHeaders(additionalHeaders ...string) Option {
	return func(m *middleware) {
		for _, additionalHeader := range additionalHeaders {
			if slices.Contains(gGenericHeaders, additionalHeader) ||
				slices.Contains(m.additionalHeaders, additionalHeader) {

				continue
			}
			m.additionalHeaders = append(m.additionalHeaders, additionalHeader)
		}
	}
}

func WithIPSaver(cb IPSaver) Option {
	return func(m *middleware) {
		if cb != nil {
			m.realIpSaver = cb
		}
	}
}
