package validator

import (
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

// to cache the validator (recommended by the docs)
var uni *ut.UniversalTranslator
var validate = validator.New(validator.WithRequiredStructEnabled())
var trans ut.Translator

type Validator interface {
	Validate() error
}

func Struct(i interface{}) error {
	return validate.Struct(i)
}

func Translate(err error) string {
	if translable, ok := err.(validator.ValidationErrors); ok {
		for _, err := range translable {
			return err.Translate(trans)
		}
	}

	return err.Error()
}

func init() {
	en := en.New()
	uni = ut.New(en, en)
	trans, _ = uni.GetTranslator("en")
	en_translations.RegisterDefaultTranslations(validate, trans)

	// TODO: register custom validators
}
