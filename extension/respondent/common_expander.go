package ekaweb_respondent

import (
	"errors"

	"github.com/inaneverb/ekacore/ekaunsafe/v4"
	"github.com/inaneverb/ekaweb"
)

// CommonExpander is an error "expander". Implements Expander interface.
// The core of the logic is Expand() method. It allows you to transform error
// to the Manifest object that can be used as an HTTP response.
//
// The rules of HOW error will be transformed into Manifest
// should be provided by you. Technically CommonExpander is just a database
// of such rules and it just "routes" provided error
// to the correspondent manifest (or its generator).
//
// It's allow-to-use object after instantiating this type but feel free to use
// its constructor NewCommonExpander().
type CommonExpander struct {

	// The logic is:
	// 1. Direct. An errors that are OK to be compared directly are stored here.
	// 2. Deep. Here stored an errors that must be matched with go1.13 errors.Is() API.
	// 3. Fallback. Call this extractor if other ways generates nothing.

	direct   map[error]ManifestExtractor
	deep     []commonExpanderDeepPair
	check    []commonExpanderCheckPair
	fallback ManifestExtractor
}

// ManifestExtractor is an alias for function that must take an error,
// "expand" it and return an Manifest the fields of which are populated
// by the error's content. It acts almost the same as Expander.Expand() method do.
type ManifestExtractor = func(err error) *Manifest

// commonExpanderDeepPair is a unit that contains:
//  1. An error that will be used at the error pattern matching
//  2. A manifest extractor callback that will be called
//     if pattern matching is succeeded.
type commonExpanderDeepPair struct {
	Pattern   error
	Extractor ManifestExtractor
}

// commonExpanderCheckPair is a unit that contains:
//  1. A callback that will tell whether provided error is that one,
//     its extractor wants to handle.
//  2. A manifest extractor callback that will be called
//     if check callback reports true.
type commonExpanderCheckPair struct {
	Checker   func(err error) bool
	Extractor ManifestExtractor
}

// Expand returns an Manifest based on the provided error.
// The registered rules and extractors are used to figure out what Manifest
// is correspondent to the passed error. Implements Expander interface.
func (ce *CommonExpander) Expand(err error) *Manifest {

	// "if err == nil" statement reports false when err is (*someUserErrorType)(nil)
	// for example. So, we're using more explicit check here.

	if ce == nil || ekaunsafe.UnpackInterface(err).Word == nil {
		return nil
	}

	var manifest *Manifest
	var iErr = ekaunsafe.UnpackInterface(err)
	var errKind = ekaunsafe.ReflectTypeOfRType(iErr.Type).Kind()

	if isGoHashableObject(errKind) {
		if extractor := ce.direct[err]; extractor != nil {
			manifest = extractor(err)
		}
	}

	for i, n := 0, len(ce.deep); i < n && manifest == nil; i++ {
		if errors.Is(err, ce.deep[i].Pattern) {
			manifest = ce.deep[i].Extractor(err)
			break
		}
	}

	for i, n := 0, len(ce.check); i < n && manifest == nil; i++ {
		if ce.check[i].Checker(err) {
			manifest = ce.check[i].Extractor(err)
			break
		}
	}

	if ce.fallback != nil && manifest == nil {
		manifest = ce.fallback(err)
	}

	return manifest
}

// ExtractorFor allows you to specify custom extractor that will populate Manifest
// from the error that is matched with the provided one by the Expand() method.
//
// A deep flag leads to match errors in Expand() method using go1.13 errors.Is() API.
// It's slower than direct comparison, so use it only when you're sure.
//
// This method can be chained.
func (ce *CommonExpander) ExtractorFor(err error, deep bool, cb ManifestExtractor) *CommonExpander {

	if ce == nil || err == nil || cb == nil {
		return ce
	}

	if deep {
		ce.deep = append(ce.deep, commonExpanderDeepPair{err, cb})
	} else {
		if ce.direct == nil {
			ce.direct = make(map[error]ManifestExtractor)
		}
		ce.direct[err] = cb
	}

	return ce
}

// ExtractorByChecker allows you to specify custom extractor that will populate Manifest
// from the error if your checker reports that it wants to handle that error.
//
// This method can be chained.
func (ce *CommonExpander) ExtractorByChecker(checker func(err error) bool, cb ManifestExtractor) *CommonExpander {
	if ce != nil && checker != nil && cb != nil {
		ce.check = append(ce.check, commonExpanderCheckPair{checker, cb})
	}
	return ce
}

// ManifestFor allows you to specify an Manifest that must be returned
// by Expand() method for the error that is matched with the provided one
//
// A deep flag leads to match errors in Expand() method using go1.13 errors.Is() API.
// It's slower than direct comparison, so use it only when you're sure.
//
// This method can be chained.
func (ce *CommonExpander) ManifestFor(err error, deep bool, manifest *Manifest) *CommonExpander {
	return ce.ExtractorFor(err, deep, ce.createSimpleExtractor(manifest))
}

// FallbackExtractor allows you to specify an extractor that will populate Manifest
// from the error when other extractors returns a nil or if there's no other extractors.
//
// This method can be chained.
func (ce *CommonExpander) FallbackExtractor(cb ManifestExtractor) *CommonExpander {
	if ce != nil && cb != nil {
		ce.fallback = cb
	}
	return ce
}

func (_ *CommonExpander) createSimpleExtractor(manifest *Manifest) ManifestExtractor {
	if manifest != nil {
		return func(_ error) *Manifest { return manifest }
	} else {
		return nil
	}
}

// NewCommonExpander is a CommonExpander's constructor. It initializes all
// internal fields and returns a ready-to-use object with some default values.
func NewCommonExpander() *CommonExpander {
	return &CommonExpander{
		direct: make(map[error]ManifestExtractor),
		deep:   make([]commonExpanderDeepPair, 0, 32),
		check:  make([]commonExpanderCheckPair, 0, 32),
		fallback: func(err error) *Manifest {
			return &Manifest{
				Status:      ekaweb.StatusInternalServerError,
				Error:       "Unrecoverable internal server error",
				ErrorCode:   ekaweb.StatusInternalServerError * 100,
				ErrorDetail: err.Error(),
			}
		},
	}
}

/////////////////////////////////////////////////////////////////////////////////
//////////////////////////////// A SHORT ALIASES ////////////////////////////////
/////////////////////////////////////////////////////////////////////////////////

// WithoutDetail is an alias for ManifestFor() method. It populates Status,
// ErrorCode, Error fields of the Manifest by the provided values.
func (ce *CommonExpander) WithoutDetail(err error, status, errorCode int, msg string) *CommonExpander {
	return ce.ManifestFor(err, false, &Manifest{
		Status: status, ErrorCode: errorCode, Error: msg,
	})
}

// WithDetail is an alias for ManifestFor() method. It populates Status,
// ErrorCode, Error, ErrorDetail fields of the Manifest by the provided values.
func (ce *CommonExpander) WithDetail(err error, status, errorCode int, msg, reason string) *CommonExpander {
	return ce.ManifestFor(err, false, &Manifest{
		Status: status, ErrorCode: errorCode, Error: msg, ErrorDetail: reason,
	})
}

// WithDetails is an alias for ManifestFor() method. It populates Status,
// ErrorCode, Error, ErrorDetails fields of the Manifest by the provided values.
func (ce *CommonExpander) WithDetails(err error, status, errorCode int, msg string, reasons []string) *CommonExpander {
	return ce.ManifestFor(err, false, &Manifest{
		Status: status, ErrorCode: errorCode, Error: msg, ErrorDetails: reasons,
	})
}

// WithCustomFillers is an alias for ManifestFor() method. It populates Status,
// ErrorCode, Error, ErrorDetailExtractor fields of the Manifest by the provided values.
func (ce *CommonExpander) WithCustomFillers(err error, status, errorCode int, msg string, ext ...ManifestCustomFiller) *CommonExpander {
	return ce.ManifestFor(err, false, &Manifest{
		Status: status, ErrorCode: errorCode, Error: msg, customFillers: ext,
	})
}

// DeepWithoutDetail is an alias for ManifestFor() method. It populates Status,
// ErrorCode, Error fields of the Manifest by the provided values.
// It also marks err to be deeply matched in the Expand() method.
func (ce *CommonExpander) DeepWithoutDetail(err error, status, errorCode int, msg string) *CommonExpander {
	return ce.ManifestFor(err, true, &Manifest{
		Status: status, ErrorCode: errorCode, Error: msg,
	})
}

// DeepWithDetail is an alias for ManifestFor() method. It populates Status,
// ErrorCode, Error, ErrorDetail fields of the Manifest by the provided values.
// It also marks err to be deeply matched in the Expand() method.
func (ce *CommonExpander) DeepWithDetail(err error, status, errorCode int, msg, reason string) *CommonExpander {
	return ce.ManifestFor(err, true, &Manifest{
		Status: status, ErrorCode: errorCode, Error: msg, ErrorDetail: reason,
	})
}

// DeepWithDetails is an alias for ManifestFor() method. It populates Status,
// ErrorCode, Error, ErrorDetails fields of the Manifest by the provided values.
// It also marks err to be deeply matched in the Expand() method.
func (ce *CommonExpander) DeepWithDetails(err error, status, errorCode int, msg string, reasons []string) *CommonExpander {
	return ce.ManifestFor(err, true, &Manifest{
		Status: status, ErrorCode: errorCode, Error: msg, ErrorDetails: reasons,
	})
}

// DeepWithCustomFillers is an alias for ManifestFor() method. It populates Status,
// ErrorCode, Error, ErrorDetailsExtractor fields of the Manifest by the provided values.
// It also marks err to be deeply matched in the Expand() method.
func (ce *CommonExpander) DeepWithCustomFillers(err error, status, errorCode int, msg string, ext ...ManifestCustomFiller) *CommonExpander {
	return ce.ManifestFor(err, true, &Manifest{
		Status: status, ErrorCode: errorCode, Error: msg, customFillers: ext,
	})
}
