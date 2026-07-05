package filter

import (
	"github.com/gofiber/fiber/v3"
	"github.com/rizalarfiyan/be-plant-factory/internal/shared/helper"
)

func ParseBool(key string, defaultBools ...bool) func(c fiber.Ctx) (bool, error) {
	defaultBool := helper.GetDefault(defaultBools)

	return func(c fiber.Ctx) (bool, error) {
		return fiber.Query(c, key, defaultBool), nil
	}
}
