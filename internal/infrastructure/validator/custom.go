package validator

import (
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/rizalarfiyan/be-plant-factory/internal/domain/entity"
)

type EnumValid interface {
	Valid() bool
}

func (v *Validate) CustomValidation() error {
	customValidations := []func() error{
		v.registerEnumValidation,
		v.registerStorageTypeValidation,
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

func (v *Validate) registerStorageTypeValidation() error {
	return v.val.RegisterValidation("storage", func(fl validator.FieldLevel) bool {
		field := fl.Field()
		if field.Kind() != reflect.String {
			return false
		}

		param := fl.Param()
		if param == "" {
			return false
		}

		storageType := entity.StorageType(param)
		return storageType.IsValidFile(field.String())
	})
}
