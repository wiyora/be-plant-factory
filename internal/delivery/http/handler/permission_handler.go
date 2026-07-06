package handler

import (
	"github.com/gofiber/fiber/v3"
	"github.com/rizalarfiyan/be-plant-factory/internal/delivery/http/response"
	"github.com/rizalarfiyan/be-plant-factory/internal/shared/code"
	"github.com/rizalarfiyan/be-plant-factory/internal/usecase/permission"
	"github.com/samber/do/v2"
)

type PermissionHandler interface {
	All(c fiber.Ctx) error
}

type permissionHandler struct {
	useCase permission.PermissionUseCase `do:""`
}

func NewPermissionHandler(i do.Injector) (PermissionHandler, error) {
	return do.InvokeStruct[permissionHandler](i)
}

// GetAllPermissions godoc
//
//	@Summary		Get All Permissions
//	@Description	Get all available permission modules with their permissions
//	@ID				permission-list
//	@Tags			Permission
//	@Produce		json
//	@Security		CookieAccessToken
//	@Success		200	{object}	response.BaseSwaggerResponse{data=[]response.PermissionModule}	"Permissions fetched successfully. Available code (SUCCESS)"
//	@Failure		401	{object}	response.BaseSwaggerEmptyResponse{}								"Unauthorized. Available code (UNAUTHORIZED)"
//	@Failure		500	{object}	response.BaseSwaggerEmptyResponse{}								"Internal Server Error. Available code (INTERNAL_SERVER_ERROR)"
//	@Router			/permission [get]
func (h permissionHandler) All(c fiber.Ctx) error {
	items := h.useCase.All(c.Context())

	res := make([]response.PermissionModule, len(items))
	for i, item := range items {
		res[i] = response.NewPermissionModule(item)
	}

	return response.New(c, code.OK, response.WithData(res))
}
