package filter

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/rizalarfiyan/be-plant-factory/internal/domain/entity"
	domainError "github.com/rizalarfiyan/be-plant-factory/internal/domain/error"
	"github.com/rizalarfiyan/be-plant-factory/internal/shared/code"
	"github.com/rizalarfiyan/be-plant-factory/internal/shared/helper"
)

type ValidEnum interface {
	~string
	Valid() bool
}

type EnumParser[T ValidEnum] func(c fiber.Ctx) (T, error)

func ParseEnum[T ValidEnum](key string, defaultEnums ...T) EnumParser[T] {
	defaultEnum := helper.GetDefault(defaultEnums)

	return func(c fiber.Ctx) (T, error) {
		valStr := strings.TrimSpace(c.Query(key))
		if valStr == "" {
			return defaultEnum, nil
		}

		enumVal := T(strings.ToLower(valStr))
		if !enumVal.Valid() {
			return defaultEnum, domainError.New(code.InvalidFilterEnumQuery, domainError.WithParams(map[string]any{
				"value": valStr,
				"key":   key,
			}))
		}

		return enumVal, nil
	}
}

func ParseSoftDelete() EnumParser[entity.SoftDeleteFilter] {
	return ParseEnum("deleted", entity.WithoutDeleted)
}
