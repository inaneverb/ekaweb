package ekaweb_bind

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/inaneverb/ekaweb"
)

func ScanAndValidateJSON(r *http.Request, to any) error {
	return scanAndValidate(r, to, bJSON)
}

func ScanAndValidateXML(r *http.Request, to any) error {
	return scanAndValidate(r, to, bXML)
}

func ScanAndValidateQuery(r *http.Request, to any) error {
	return scanAndValidate(r, to, bQuery)
}

func ScanAndValidateHeader(r *http.Request, to any) error {
	return scanAndValidate(r, to, bHeader)
}

func ScanAndValidateForm(r *http.Request, to any) error {
	return scanAndValidate(r, to, bForm)
}

func ScanAndValidateFormPost(r *http.Request, to any) error {
	return scanAndValidate(r, to, bFormPost)
}

func ScanAndValidateFormMultipart(r *http.Request, to any) error {
	return scanAndValidate(r, to, bFormMultipart)
}

func OnlyValidate(to any) error {
	return scanAndValidate(nil, to, bOnlyValidate)
}

func scanAndValidate(r *http.Request, to any, b Binding) error {

	err := b.Bind(r, to)
	if err == nil {
		return nil
	}

	if validationErr, ok := err.(validator.ValidationErrors); ok {
		return &errValidationFailed{validationErr}
	}

	ekaweb.ErrorDetailApply(r, "Malformed HTTP request "+b.Name()+" source")
	return ErrMalformedSource
}
