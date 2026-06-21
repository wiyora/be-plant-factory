package request

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func GetParamID(c fiber.Ctx, key string) (uuid.UUID, error) {
	id, err := uuid.Parse(c.Params(key))
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}
