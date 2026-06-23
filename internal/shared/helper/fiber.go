package helper

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/rizalarfiyan/be-plant-factory/internal/domain/entity"
	"github.com/rizalarfiyan/be-plant-factory/internal/shared/constant"
)

func GetLocal[T any](c fiber.Ctx, key string) (T, bool) {
	val := c.Locals(key)
	if val == nil {
		var zero T
		return zero, false
	}

	res, ok := val.(T)
	return res, ok
}

func GetLocalUser(c fiber.Ctx) (entity.AuthContext, bool) {
	return GetLocal[entity.AuthContext](c, constant.KeyAuthUser)
}

func GetLocalUserID(c fiber.Ctx) (uuid.UUID, bool) {
	return GetLocal[uuid.UUID](c, constant.KeyAuthUserID)
}

func GetLocalSessionID(c fiber.Ctx) (uuid.UUID, bool) {
	return GetLocal[uuid.UUID](c, constant.KeyAuthSessionID)
}

func GetLocalAccessToken(c fiber.Ctx) (entity.AccessTokenContext, bool) {
	return GetLocal[entity.AccessTokenContext](c, constant.KeyAuthAccessToken)
}
