package handler

import (
	"github.com/gofiber/fiber/v3"
	"github.com/rizalarfiyan/be-plant-factory/internal/delivery/http/response"
	"github.com/rizalarfiyan/be-plant-factory/internal/shared/code"
	"github.com/rizalarfiyan/be-plant-factory/internal/usecase/health"
	"github.com/samber/do/v2"
)

type HealthHandler interface {
	Check(c fiber.Ctx) error
}

type healthHandler struct {
	health health.HealthUseCase `do:""`
}

func NewHealthHandler(i do.Injector) (HealthHandler, error) {
	return do.InvokeStruct[healthHandler](i)
}

// Check godoc
//
//	@Summary		Check Health
//	@Description	Check Health of the application
//	@ID				check-health
//	@Tags			Health
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	response.BaseSwaggerResponse{data=response.HealthResponse}	"Health check successful. Available code (SUCCESS)"
//	@Router			/health [get]
func (h healthHandler) Check(c fiber.Ctx) error {
	res := response.NewHealthResponse(h.health.Check())
	return response.New(c, code.OK, response.WithData(res))
}
