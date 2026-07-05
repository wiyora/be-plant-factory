package helper

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	domainError "github.com/rizalarfiyan/be-plant-factory/internal/domain/error"
	"github.com/rizalarfiyan/be-plant-factory/internal/shared/constant"
)

func BindJSON[T any](ctx fiber.Ctx) (T, error) {
	var payload T
	if err := ctx.Bind().SkipValidation(true).JSON(&payload); err != nil {
		return payload, err
	}

	validate, isExist := fiber.GetState[*validator.Validate](ctx.App().State(), constant.AppStateValidatorKey)
	if !isExist {
		return payload, domainError.ErrValidationStateNotSet
	}

	if err := validate.StructCtx(ctx.Context(), payload); err != nil {
		return payload, err
	}

	if v, ok := any(payload).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return payload, err
		}
	}

	return payload, nil
}

type QueryBinder func(c fiber.Ctx) error

func QueryBind(c fiber.Ctx, binders ...QueryBinder) error {
	for _, bind := range binders {
		if err := bind(c); err != nil {
			return err
		}
	}
	return nil
}

func QueryField[T any](target *T, parseFn func(fiber.Ctx) (T, error)) QueryBinder {
	return func(c fiber.Ctx) error {
		val, err := parseFn(c)
		if err != nil {
			return err
		}

		*target = val
		return nil
	}
}
