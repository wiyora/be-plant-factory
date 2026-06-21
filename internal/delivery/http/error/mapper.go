package error

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/rizalarfiyan/be-plant-factory/internal/delivery/http/response"
	domainError "github.com/rizalarfiyan/be-plant-factory/internal/domain/error"
	"github.com/rizalarfiyan/be-plant-factory/internal/shared/code"
)

func Handle(c fiber.Ctx, err error) error {
	if err == nil {
		return response.New(c, code.InternalServerError)
	}

	if manualErr, ok := domainError.IsManualValidationError(err); ok {
		fields := make(response.FieldsResponse, len(manualErr.Fields))
		for field, code := range manualErr.Fields {
			fields[field] = response.CodeParamsResponse{
				Code:   strings.ToUpper(code.Code),
				Params: code.Params,
			}
		}

		return response.New(c, code.ValidationError, response.WithFields(fields))
	}

	if appErr, ok := domainError.IsAppError(err); ok {
		return response.New(c, appErr.Code, response.WithParams(appErr.Params), response.WithData(appErr.Data))
	}

	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		return response.New(c, code.Parse(fiberErr.Code))
	}

	return response.New(c, code.InternalServerError)
}
