package ekaweb_bind

import (
	"errors"
	"fmt"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	entr "github.com/go-playground/validator/v10/translations/en"

	"github.com/inaneverb/ekaweb/extension/respondent"
)

var (
	// defaultEnTranslator is the default EN translator that could be used
	// to provide fast access (or fast fallback) to the default translating.
	defaultEnTranslator ut.Translator
)

var (
	ErrMalformedSource  = errors.New("Extension.Binding: Malformed HTTP source")
	ErrValidationFailed = (*errValidationFailed)(nil)
)

// RespondentExpanderChecker is a checker for respondent.CommonExpander that will
// return true only if provided error was created by the any ScanAndValidate... function.
func RespondentExpanderChecker(err error) bool {
	_, ok := err.(validator.ValidationErrors)
	return ok
}

func NewRespondentManifestExtractor(status, errorCode int, message string) ekaweb_respondent.ManifestExtractor {
	return func(err error) *ekaweb_respondent.Manifest {

		manifest := ekaweb_respondent.Manifest{
			Status:    status,
			ErrorCode: errorCode,
			Error:     message,
		}

		var errList validator.ValidationErrors
		if errList1, ok := err.(validator.ValidationErrors); ok {
			errList = errList1
		} else if typedErr, ok := err.(*errValidationFailed); ok {
			errList = typedErr.originalErr
		} else {
			return nil
		}

		manifest.ErrorDetails = make([]string, len(errList))
		for i, n := 0, len(errList); i < n; i++ {
			// if defaultTranslator is nil, errList[i] underlying err.Error()
			// will be called.
			if defaultEnTranslator != nil {
				manifest.ErrorDetails[i] = errList[i].Translate(defaultEnTranslator)
			} else {
				manifest.ErrorDetails[i] = errList[i].Error()
			}
		}

		return &manifest
	}
}

func init() {
	var v = Validator.Engine().(*validator.Validate)
	var tr = en.New()
	var uni = ut.New(tr, tr)
	var foundDefaultEnTranslator bool

	defaultEnTranslator, foundDefaultEnTranslator = uni.GetTranslator("en")
	if !foundDefaultEnTranslator {
		panic("Extension.Binding: EN translator not found")
	}

	if err := entr.RegisterDefaultTranslations(v, defaultEnTranslator); err != nil {
		const D = "Extension.Binding: Failed to register defaults for EN translator: %s"
		panic(fmt.Sprintf(D, err.Error()))
	}
}

type errValidationFailed struct {
	originalErr validator.ValidationErrors
}

func (e *errValidationFailed) Error() string {
	return "Extension.Binding: Validation failed"
}

func (e *errValidationFailed) Is(other error) bool {
	_, ok := other.(*errValidationFailed)
	return ok
}
