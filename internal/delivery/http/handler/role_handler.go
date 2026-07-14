package handler

import (
	"github.com/gofiber/fiber/v3"
	"github.com/rizalarfiyan/be-plant-factory/internal/delivery/http/request"
	"github.com/rizalarfiyan/be-plant-factory/internal/delivery/http/request/filter"
	"github.com/rizalarfiyan/be-plant-factory/internal/delivery/http/request/pagination"
	"github.com/rizalarfiyan/be-plant-factory/internal/delivery/http/request/search"
	"github.com/rizalarfiyan/be-plant-factory/internal/delivery/http/request/sorting"
	"github.com/rizalarfiyan/be-plant-factory/internal/delivery/http/response"
	"github.com/rizalarfiyan/be-plant-factory/internal/domain/entity"
	"github.com/rizalarfiyan/be-plant-factory/internal/shared/code"
	"github.com/rizalarfiyan/be-plant-factory/internal/shared/helper"
	"github.com/rizalarfiyan/be-plant-factory/internal/usecase/role"
	"github.com/samber/do/v2"
)

type RoleHandler interface {
	List(c fiber.Ctx) error
	Detail(c fiber.Ctx) error
	Create(c fiber.Ctx) error
	Update(c fiber.Ctx) error
	Delete(c fiber.Ctx) error
	Dropdown(c fiber.Ctx) error
}

type roleHandler struct {
	useCase role.RoleUseCase `do:""`
}

func NewRoleHandler(i do.Injector) (RoleHandler, error) {
	return do.InvokeStruct[roleHandler](i)
}

// List godoc
//
//	@Summary		List Roles
//	@Description	Get list of roles with pagination and search
//	@ID				role-list
//	@Tags			Role
//	@Produce		json
//	@Security		CookieAccessToken
//	@Param			search		query		string																		false	"Search by name (min 3 chars)"
//	@Param			page		query		int																			false	"Page number"		default(1)
//	@Param			page_size	query		int																			false	"Items per page"	default(10)
//	@Param			order_by	query		string																		false	"Order by field"	default(id)		enums(id, name, created_at, total_permission)
//	@Param			sort_by		query		string																		false	"Sort direction"	default(desc)	enums(asc, desc)
//	@Success		200			{object}	response.BasePaginationSwaggerResponse{data=[]response.ListRoleResponse}	"Roles fetched successfully. Available code (LIST_FETCHED)"
//	@Failure		400			{object}	response.BaseSwaggerValidationResponse{}									"Bad Request. Available code (VALIDATION_ERROR, BAD_REQUEST)"
//	@Failure		401			{object}	response.BaseSwaggerEmptyResponse{}											"Unauthorized. Available code (UNAUTHORIZED)"
//	@Failure		500			{object}	response.BaseSwaggerEmptyResponse{}											"Internal Server Error. Available code (INTERNAL_SERVER_ERROR)"
//	@Router			/role [get]
func (h roleHandler) List(c fiber.Ctx) error {
	var req entity.RoleFilter

	err := helper.QueryBind(c,
		helper.QueryField(&req.Search, search.Parse()),
		helper.QueryField(&req.Pagination, pagination.Parse()),
		helper.QueryField(&req.Order,
			sorting.Parse(
				sorting.WithDefaults("id", entity.Desc),
				sorting.WithAllowedOrderBy("id", "name", "created_at", "total_permission"),
				sorting.WithMappingOrderBy(map[string]string{
					"id":               "id",
					"name":             "name",
					"created_at":       "id",
					"total_permission": "total_permission",
				}),
			)),
	)
	if err != nil {
		return err
	}

	items, total, err := h.useCase.List(c.Context(), req)
	if err != nil {
		return err
	}

	res := make([]response.ListRoleResponse, len(items))
	for i, item := range items {
		res[i] = response.NewListRoleResponse(item)
	}

	return response.New(c, code.ListFetched, response.WithData(res), response.WithPagination(req.Pagination.ToResult(total)))
}

// Detail godoc
//
//	@Summary		Get Role Detail
//	@Description	Get role detail by ID with permissions
//	@ID				role-detail
//	@Tags			Role
//	@Produce		json
//	@Security		CookieAccessToken
//	@Param			id	path		string															true	"Role ID"
//	@Success		200	{object}	response.BaseSwaggerResponse{data=response.DetailRoleResponse}	"Role detail fetched successfully. Available code (DETAIL_FETCHED)"
//	@Failure		400	{object}	response.BaseSwaggerEmptyResponse{}								"Bad Request - invalid ID. Available code (INVALID_PARAM_ID)"
//	@Failure		401	{object}	response.BaseSwaggerEmptyResponse{}								"Unauthorized. Available code (UNAUTHORIZED)"
//	@Failure		404	{object}	response.BaseSwaggerEmptyResponse{}								"Role not found. Available code (NOT_FOUND)"
//	@Failure		422	{object}	response.BaseSwaggerEmptyResponse{}								"Role not found. Available code (ROLE_NOT_FOUND)"
//	@Failure		500	{object}	response.BaseSwaggerEmptyResponse{}								"Internal Server Error. Available code (INTERNAL_SERVER_ERROR)"
//	@Router			/role/{id} [get]
func (h roleHandler) Detail(c fiber.Ctx) error {
	id, err := request.GetParamID(c, "id")
	if err != nil {
		return response.New(c, code.InvalidParamID)
	}

	item, err := h.useCase.Detail(c.Context(), id)
	if err != nil {
		return err
	}

	res := response.NewDetailRoleResponse(item)
	return response.New(c, code.DetailFetched, response.WithData(res))
}

// Create godoc
//
//	@Summary		Create Role
//	@Description	Create a new role. Validation: name REQUIRED, MIN 3, MAX 32; permissions REQUIRED, MIN 1
//	@ID				role-create
//	@Tags			Role
//	@Accept			json
//	@Produce		json
//	@Security		CookieAccessToken
//	@Param			payload	body		request.CreateUpdateRoleRequest				true	"Create Role Payload"
//	@Success		201		{object}	response.BaseSwaggerEmptyResponse{}			"Role created successfully. Available code (CREATED)"
//	@Failure		400		{object}	response.BaseSwaggerValidationResponse{}	"Bad Request - invalid input data. Available code (VALIDATION_ERROR, BAD_REQUEST)"
//	@Failure		401		{object}	response.BaseSwaggerEmptyResponse{}			"Unauthorized. Available code (UNAUTHORIZED)"
//	@Failure		422		{object}	response.BaseSwaggerEmptyResponse{}			"Role name already exists. Available code (ROLE_NAME_EXISTS)"
//	@Failure		500		{object}	response.BaseSwaggerEmptyResponse{}			"Internal Server Error. Available code (INTERNAL_SERVER_ERROR)"
//	@Router			/role [post]
func (h roleHandler) Create(c fiber.Ctx) error {
	payload, err := helper.BindJSON[request.CreateUpdateRoleRequest](c)
	if err != nil {
		return response.NewValidate(c, err)
	}

	req := payload.ToEntity()
	if err := h.useCase.Create(c.Context(), req); err != nil {
		return err
	}

	return response.New(c, code.Created)
}

// Update godoc
//
//	@Summary		Update Role
//	@Description	Update role name and permissions. Validation: name REQUIRED, MIN 3, MAX 32; permissions REQUIRED, MIN 1
//	@ID				role-update
//	@Tags			Role
//	@Accept			json
//	@Produce		json
//	@Security		CookieAccessToken
//	@Param			id		path		string										true	"Role ID"
//	@Param			payload	body		request.CreateUpdateRoleRequest				true	"Update Role Payload"
//	@Success		200		{object}	response.BaseSwaggerEmptyResponse{}			"Role updated successfully. Available code (UPDATED)"
//	@Failure		400		{object}	response.BaseSwaggerValidationResponse{}	"Bad Request - invalid input data. Available code (VALIDATION_ERROR, BAD_REQUEST)"
//	@Failure		401		{object}	response.BaseSwaggerEmptyResponse{}			"Unauthorized. Available code (UNAUTHORIZED)"
//	@Failure		422		{object}	response.BaseSwaggerEmptyResponse{}			"Role not found or name already exists. Available code (ROLE_NOT_FOUND, ROLE_NAME_EXISTS)"
//	@Failure		500		{object}	response.BaseSwaggerEmptyResponse{}			"Internal Server Error. Available code (INTERNAL_SERVER_ERROR)"
//	@Router			/role/{id} [put]
func (h roleHandler) Update(c fiber.Ctx) error {
	id, err := request.GetParamID(c, "id")
	if err != nil {
		return response.New(c, code.InvalidParamID)
	}

	payload, err := helper.BindJSON[request.CreateUpdateRoleRequest](c)
	if err != nil {
		return response.NewValidate(c, err)
	}

	req := payload.ToEntity()
	req.ID = id
	if err := h.useCase.Update(c.Context(), req); err != nil {
		return err
	}

	return response.New(c, code.Updated)
}

// Delete godoc
//
//	@Summary		Delete Role
//	@Description	Delete role by ID
//	@ID				role-delete
//	@Tags			Role
//	@Produce		json
//	@Security		CookieAccessToken
//	@Param			id	path		string								true	"Role ID"
//	@Success		200	{object}	response.BaseSwaggerEmptyResponse{}	"Role deleted successfully. Available code (DELETED)"
//	@Failure		400	{object}	response.BaseSwaggerEmptyResponse{}	"Bad Request - invalid ID. Available code (INVALID_PARAM_ID)"
//	@Failure		401	{object}	response.BaseSwaggerEmptyResponse{}	"Unauthorized. Available code (UNAUTHORIZED)"
//	@Failure		422	{object}	response.BaseSwaggerEmptyResponse{}	"Role not found. Available code (ROLE_NOT_FOUND)"
//	@Failure		500	{object}	response.BaseSwaggerEmptyResponse{}	"Internal Server Error. Available code (INTERNAL_SERVER_ERROR)"
//	@Router			/role/{id} [delete]
func (h roleHandler) Delete(c fiber.Ctx) error {
	id, err := request.GetParamID(c, "id")
	if err != nil {
		return response.New(c, code.InvalidParamID)
	}

	if err := h.useCase.Delete(c.Context(), id); err != nil {
		return err
	}

	return response.New(c, code.Deleted)
}

// Dropdown godoc
//
//	@Summary		Role Dropdown
//	@Description	Get role dropdown with search, pagination, and active_ids
//	@ID				role-dropdown
//	@Tags			Role
//	@Produce		json
//	@Security		CookieAccessToken
//	@Param			search		query		string																		false	"Search by name (min 3 chars)"
//	@Param			page		query		int																			false	"Page number"		default(1)
//	@Param			page_size	query		int																			false	"Items per page"	default(10)
//	@Param			active_ids	query		[]string																	false	"Active IDs"
//	@Success		200			{object}	response.BasePaginationSwaggerResponse{data=[]response.DropdownResponse}	"Roles fetched successfully. Available code (LIST_FETCHED)"
//	@Failure		400			{object}	response.BaseSwaggerValidationResponse{}									"Bad Request. Available code (VALIDATION_ERROR, BAD_REQUEST)"
//	@Failure		401			{object}	response.BaseSwaggerEmptyResponse{}											"Unauthorized. Available code (UNAUTHORIZED)"
//	@Failure		500			{object}	response.BaseSwaggerEmptyResponse{}											"Internal Server Error. Available code (INTERNAL_SERVER_ERROR)"
//	@Router			/role/dropdown [get]
func (h roleHandler) Dropdown(c fiber.Ctx) error {
	var req entity.DropdownFilter

	err := helper.QueryBind(c,
		helper.QueryField(&req.Search, search.Parse()),
		helper.QueryField(&req.Pagination, pagination.Parse()),
		helper.QueryField(&req.ActiveIDs, filter.ParseUUIDs("active_ids")),
	)
	if err != nil {
		return err
	}

	items, total, err := h.useCase.Dropdown(c.Context(), req)
	if err != nil {
		return err
	}

	res := make([]response.DropdownResponse, len(items))
	for i, item := range items {
		res[i] = response.NewDropdownResponse(item)
	}

	return response.New(c, code.ListFetched, response.WithData(res), response.WithPagination(req.Pagination.ToResult(total)))
}
