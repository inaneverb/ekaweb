package ekaweb_bind

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

type CustomRuleApplicator func(v *validator.Validate, tr ut.Translator) error
type CustomTypeApplicator func(v *validator.Validate) error

func RegisterCustomValidators(cb ...CustomRuleApplicator) error {
	v := Validator.Engine().(*validator.Validate)

	for i, n := 0, len(cb); i < n; i++ {
		if cb[i] == nil {
			continue
		}
		err := cb[i](v, defaultEnTranslator)
		if err != nil {
			return err
		}
	}

	return nil
}

func RegisterCustomTypes(cb ...CustomTypeApplicator) error {
	v := Validator.Engine().(*validator.Validate)

	for i, n := 0, len(cb); i < n; i++ {
		if cb[i] == nil {
			continue
		}
		err := cb[i](v)
		if err != nil {
			return err
		}
	}

	return nil
}
