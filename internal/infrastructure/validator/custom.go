package validator

import (
	"github.com/go-playground/validator/v10"
)

type EnumValid interface {
	Valid() bool
}

func (v *Validate) CustomValidation() error {
	customValidations := []func() error{
		v.registerEnumValidation,
	}

	for _, validation := range customValidations {
		if err := validation(); err != nil {
			return err
		}
	}

	return nil
}

func (v *Validate) registerEnumValidation() error {
	return v.val.RegisterValidation("enum", func(fl validator.FieldLevel) bool {
		if enum, ok := fl.Field().Interface().(EnumValid); ok {
			return enum.Valid()
		}

		return false
	})
}
