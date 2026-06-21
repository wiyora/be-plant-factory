package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/rizalarfiyan/be-plant-factory/internal/infrastructure/logger"
	"github.com/rs/zerolog"
	"github.com/samber/do/v2"
)

type Validate struct {
	val *validator.Validate
}

func New(i do.Injector) (*Validate, error) {
	rawLog := do.MustInvoke[*zerolog.Logger](i)

	log := logger.WithLayer(rawLog, logger.LayerValidator)

	log.Info().Msg("Initializing validator")

	return &Validate{
		val: validator.New(),
	}, nil
}

type StructManualValidate interface {
	Validate() error
	Validator() *validator.Validate
}

func (v *Validate) Validate(out any) error {
	if err := v.val.Struct(out); err != nil {
		return err
	}

	if manual, ok := out.(StructManualValidate); ok {
		if err := manual.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func (v *Validate) Validator() *validator.Validate {
	return v.val
}
