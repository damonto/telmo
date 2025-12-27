package validator

import (
	"errors"
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
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errs := make([]string, 0, len(validationErrors))
			for _, e := range validationErrors {
				errs = append(errs, e.Translate(v.trans))
			}
			if len(errs) == 0 {
				return err
			}
			return errors.New(strings.Join(errs, ", "))
		}
		return fmt.Errorf("validation failed: %w", err)
	}
	return nil
}
