package domainError

import (
	"errors"

	"github.com/rizalarfiyan/be-plant-factory/internal/shared/code"
)

type ManualValidationField struct {
	Code   string
	Params map[string]any
}

type ManualValidationError struct {
	Fields map[string]ManualValidationField
}

func NewManualValidation(key, code string, params ...map[string]any) *ManualValidationError {
	var param map[string]any
	if len(params) > 0 {
		param = params[0]
	}

	return &ManualValidationError{
		Fields: map[string]ManualValidationField{
			key: {
				Code:   code,
				Params: param,
			},
		},
	}
}

func NewManualValidations(fields map[string]ManualValidationField) *ManualValidationError {
	return &ManualValidationError{
		Fields: fields,
	}
}

func (e *ManualValidationError) Error() string {
	return string(code.ValidationError.Code)
}

func IsManualValidationError(err error) (*ManualValidationError, bool) {
	var errs *ManualValidationError
	if ok := errors.As(err, &errs); ok {
		return errs, true
	}

	return nil, false
}
