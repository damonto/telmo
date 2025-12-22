package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

type CustomValidator struct {
	validator *validator.Validate
	trans     ut.Translator
}

func New() *CustomValidator {
	validate := validator.New()
	en := en.New()
	uni := ut.New(en, en)
	trans, _ := uni.GetTranslator("en")
	en_translations.RegisterDefaultTranslations(validate, trans)
	v := &CustomValidator{validator: validate, trans: trans}
	return v
}

func (v *CustomValidator) Validate(i any) error {
	if err := v.validator.Struct(i); err != nil {
		var errs []string
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			for _, e := range validationErrors {
				errs = append(errs, e.Translate(v.trans))
			}
		}
		return fmt.Errorf("%v", strings.Join(errs, ", "))
	}
	return nil
}
