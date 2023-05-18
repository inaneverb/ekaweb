package ekaweb_bind

import (
	"net/http"
)

type onlyValidateBinding struct{}

func (onlyValidateBinding) Name() string {
	return "ONLY-VALIDATE"
}

func (onlyValidateBinding) Bind(_ *http.Request, obj any) error {
	return validate(obj)
}

var _ Binding = onlyValidateBinding{}
