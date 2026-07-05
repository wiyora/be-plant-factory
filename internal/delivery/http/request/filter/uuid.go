package filter

import (
	"bytes"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	domainError "github.com/rizalarfiyan/be-plant-factory/internal/domain/error"
	"github.com/rizalarfiyan/be-plant-factory/internal/shared/code"
	"github.com/rizalarfiyan/be-plant-factory/internal/shared/helper"
)

type UUIDParser func(c fiber.Ctx) (uuid.UUID, error)
type UUIDsParser func(c fiber.Ctx) ([]uuid.UUID, error)

func ParseUUID(key string, appCodes ...code.AppCode) UUIDParser {
	appCode := helper.GetDefault(appCodes)
	if helper.IsEmptyStruct(appCode) {
		appCode = code.InvalidFilterUUIDQuery
	}

	return func(c fiber.Ctx) (uuid.UUID, error) {
		valStr := strings.TrimSpace(c.Query(key))
		if valStr == "" {
			return uuid.Nil, nil
		}

		parsedUUID, err := uuid.Parse(valStr)
		if err != nil {
			return uuid.Nil, domainError.New(appCode, domainError.WithParams(map[string]any{
				"value": valStr,
				"key":   key,
			}))
		}

		return parsedUUID, nil
	}
}

func ParseUUIDs(key string, appCodes ...code.AppCode) UUIDsParser {
	appCode := helper.GetDefault(appCodes)
	if helper.IsEmptyStruct(appCode) {
		appCode = code.InvalidFilterUUIDQuery
	}

	return func(c fiber.Ctx) ([]uuid.UUID, error) {
		valBytes := c.Request().URI().QueryArgs().PeekMulti(key)
		if len(valBytes) == 0 {
			return nil, nil
		}

		result := make([]uuid.UUID, 0, len(valBytes))
		seen := make(map[uuid.UUID]struct{}, len(valBytes))

		for _, b := range valBytes {
			trimmed := bytes.TrimSpace(b)
			if len(trimmed) == 0 {
				continue
			}

			id, err := uuid.ParseBytes(trimmed)
			if err != nil {
				return nil, domainError.New(appCode, domainError.WithParams(map[string]any{
					"value": string(trimmed),
					"key":   key,
				}))
			}

			if _, exists := seen[id]; !exists {
				seen[id] = struct{}{}
				result = append(result, id)
			}
		}

		return result, nil
	}
}
