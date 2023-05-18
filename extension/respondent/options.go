package ekaweb_respondent

import (
	"github.com/inaneverb/ekacore/ekaunsafe/v4"
)

// Option is a callback that allows to modify Responder under its construction.
type Option func(m *respondent)

// WithReplacer returns an Option that allows you to use custom Replacer
// instead of default no-op one.
func WithReplacer(replacer Replacer) Option {
	return func(r *respondent) {
		if ekaunsafe.UnpackInterface(replacer).Word != nil {
			r.fromOptions.replacer = replacer
		}
	}
}

// WithApplicator returns an Option that allows you to use custom Applicator
// instead of default one.
func WithApplicator(applicator Applicator) Option {
	return func(r *respondent) {
		if ekaunsafe.UnpackInterface(applicator).Word != nil {
			r.fromOptions.applicator = applicator
		}
	}
}
