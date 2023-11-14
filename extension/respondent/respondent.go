package ekaweb_respondent

import (
	"errors"
	"net/http"

	"github.com/inaneverb/ekacore/ekaunsafe/v4"
	"github.com/inaneverb/ekaweb/v2"
)

// respondent is a combination of Replacer, Expander and Applicator
// that are work together to represent error as some net IO response.
type respondent struct {
	expander Expander

	fromOptions struct {
		replacer   Replacer
		applicator Applicator
	}
}

// HttpContext is a special type that contains http.ResponseWriter and http.Request.
// It's useful when you need to pass these to some function but has a limitation
// to pass only 1 argument.
type HttpContext struct {
	W http.ResponseWriter
	R *http.Request
}

//goland:noinspection GoErrorStringFormat
var (
	// ErrEmptyReplacement is returned when Replacer returns a nil instead of error.
	// CommonReplacer could do that if it got a half-nil error as an in argument
	// (half-nil means nil value, non-nil type, like (*errType)(nil)).
	ErrEmptyReplacement = errors.New("Middleware.Respondent got unexpectedly empty replacement")

	// ErrEmptyManifest is returned when Expander returns a nil Manifest.
	// CommonExpander could do that if there's no registered manifest generator
	// for such error and no default fallback way.
	ErrEmptyManifest = errors.New("Middleware.Respondent got unexpectedly empty manifest")

	// ----------

	errBadExpander = errors.New("Middleware.Respondent is not initialized properly: Nil or incorrect Expander")
)

var (
	_ ekaweb.ErrorHandler     = (*respondent)(nil).Callback
	_ ekaweb.ErrorHandlerHTTP = (*respondent)(nil).CallbackForHTTP
)

/////////////////////////////////////////////////////////////////////////////////
///////////////////////////////// MIDDLEWARE ////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////////

// Callback is the root of Respondent.
// It's a function, that is, presents the flow of transforming occurred error
// into the some net IO response. It could return an error, but it's abnormal case.
func (rp *respondent) Callback(ctx any, err error) {

	err = rp.fromOptions.replacer.Replace(err)
	if err == nil {
		return
	}

	manifest := rp.expander.Expand(err)
	if manifest == nil {
		return
	}

	rp.fromOptions.applicator.Apply(ctx, manifest)
}

// CallbackForHTTP is a Callback() inside but with a different signature.
// The signature is intended to be compatible with types.ErrorHandlerHTTP.
func (rp *respondent) CallbackForHTTP(w http.ResponseWriter, r *http.Request, err error) {
	rp.Callback(HttpContext{w, r}, err)
}

/////////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////////

// newRespondent creates a new Respondent object based on provided Expander
// and custom options. Expander is only one required thing for Respondent.
//
// After Respondent object is created it's ready to be used as an error handler
// (kind of middleware) providing transformation occurred error into HTTP response.
func newRespondent(expander Expander, opts []Option) *respondent {
	const s = "Failed to initialize Respondent middleware."

	m := respondent{}
	m.fromOptions.replacer = newNoOpReplacer()
	m.fromOptions.applicator = NewCommonApplicator()

	for i, n := 0, len(opts); i < n; i++ {
		if opts[i] != nil {
			opts[i](&m)
		}
	}

	if ekaunsafe.UnpackInterface(expander).Word == nil {
		panic(errBadExpander)
	}

	m.expander = expander
	return &m
}

func NewGeneric(expander Expander, opts ...Option) ekaweb.ErrorHandler {
	return newRespondent(expander, opts).Callback
}

func NewForHTTP(expander Expander, opts ...Option) ekaweb.ErrorHandlerHTTP {
	return newRespondent(expander, opts).CallbackForHTTP
}
