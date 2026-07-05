package filter

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	domainError "github.com/rizalarfiyan/be-plant-factory/internal/domain/error"
	"github.com/rizalarfiyan/be-plant-factory/internal/shared/code"
	"github.com/rizalarfiyan/be-plant-factory/internal/shared/constant"
	"github.com/rizalarfiyan/be-plant-factory/internal/shared/helper"
)

type ValidationParser[T fiber.GenericType] func(c fiber.Ctx) (T, error)

func ParseValidate[T fiber.GenericType](key string, rule string, appCodes ...code.AppCode) ValidationParser[T] {
	appCode := helper.GetDefault(appCodes)
	if helper.IsEmptyStruct(appCode) {
		appCode = code.InvalidFilterQuery
	}

	return func(ctx fiber.Ctx) (T, error) {
		var zero T
		val, ok := fiber.GetState[*validator.Validate](ctx.App().State(), constant.AppStateValidatorKey)
		if !ok {
			return zero, domainError.New(appCode)
		}

		value := fiber.Query[T](ctx, key)
		err := val.Var(value, rule)
		if err != nil {
			return zero, domainError.New(appCode, domainError.WithParams(map[string]any{
				"value": value,
				"key":   key,
			}))
		}

		return value, nil
	}
}

func ParseISO4217(key string, appCodes ...code.AppCode) ValidationParser[string] {
	return ParseValidate[string](key, "omitempty,iso4217", appCodes...)
}
