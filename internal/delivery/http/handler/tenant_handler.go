package handler

import (
	"github.com/gofiber/fiber/v3"
	"github.com/rizalarfiyan/be-plant-factory/internal/delivery/http/request"
	"github.com/rizalarfiyan/be-plant-factory/internal/delivery/http/request/pagination"
	"github.com/rizalarfiyan/be-plant-factory/internal/delivery/http/request/search"
	"github.com/rizalarfiyan/be-plant-factory/internal/delivery/http/request/sorting"
	"github.com/rizalarfiyan/be-plant-factory/internal/delivery/http/response"
	"github.com/rizalarfiyan/be-plant-factory/internal/domain/entity"
	domainError "github.com/rizalarfiyan/be-plant-factory/internal/domain/error"
	"github.com/rizalarfiyan/be-plant-factory/internal/shared/code"
	"github.com/rizalarfiyan/be-plant-factory/internal/shared/helper"
	"github.com/rizalarfiyan/be-plant-factory/internal/usecase/tenant"
	"github.com/samber/do/v2"
)

type TenantHandler interface {
	List(c fiber.Ctx) error
	Detail(c fiber.Ctx) error
	Create(c fiber.Ctx) error
	Update(c fiber.Ctx) error
	UpdateStatus(c fiber.Ctx) error
}

type tenantHandler struct {
	useCase tenant.TenantUseCase `do:""`
}

func NewTenantHandler(i do.Injector) (TenantHandler, error) {
	return do.InvokeStruct[tenantHandler](i)
}

// List godoc
//
//	@Summary		List Tenants
//	@Description	Get list of tenants with pagination and search
//	@ID				tenant-list
//	@Tags			Tenant
//	@Produce		json
//	@Security		CookieAccessToken
//	@Param			search		query		string																		false	"Search by name (min 3 chars)"
//	@Param			page		query		int																			false	"Page number"		default(1)
//	@Param			limit		query		int																			false	"Items per page"	default(10)
//	@Param			order_by	query		string																		false	"Order by field"	default(id)		enums(id, created_at, name)
//	@Param			sort_by		query		string																		false	"Sort direction"	default(desc)	enums(asc, desc)
//	@Success		200			{object}	response.BasePaginationSwaggerResponse{data=[]response.ListTenantResponse}	"Tenants fetched successfully. Available code (LIST_FETCHED)"
//	@Failure		400			{object}	response.BaseSwaggerValidationResponse{}									"Bad Request. Available code (VALIDATION_ERROR, BAD_REQUEST)"
//	@Failure		401			{object}	response.BaseSwaggerEmptyResponse{}											"Unauthorized. Available code (UNAUTHORIZED)"
//	@Failure		500			{object}	response.BaseSwaggerEmptyResponse{}											"Internal Server Error. Available code (INTERNAL_SERVER_ERROR)"
//	@Router			/tenant [get]
func (h tenantHandler) List(c fiber.Ctx) error {
	var req entity.TenantFilter

	err := helper.QueryBind(c,
		helper.QueryField(&req.Search, search.Parse()),
		helper.QueryField(&req.Pagination, pagination.Parse()),
		helper.QueryField(&req.Order,
			sorting.Parse(
				sorting.WithDefaults("id", entity.Desc),
				sorting.WithAllowedOrderBy("id", "created_at", "name"),
				sorting.WithMappingOrderBy(map[string]string{
					"id":         "id",
					"created_at": "id",
					"name":       "name",
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

	res := make([]response.ListTenantResponse, len(items))
	for i, item := range items {
		res[i] = response.NewListTenantResponse(item)
	}

	return response.New(c, code.ListFetched, response.WithData(res), response.WithPagination(req.Pagination.ToResult(total)))
}

// Detail godoc
//
//	@Summary		Get Tenant Detail
//	@Description	Get tenant detail by ID
//	@ID				tenant-detail
//	@Tags			Tenant
//	@Produce		json
//	@Security		CookieAccessToken
//	@Param			id	path		string																true	"Tenant ID"
//	@Success		200	{object}	response.BaseSwaggerResponse{data=response.DetailTenantResponse}	"Tenant detail fetched successfully. Available code (DETAIL_FETCHED)"
//	@Failure		400	{object}	response.BaseSwaggerEmptyResponse{}									"Bad Request - invalid ID. Available code (INVALID_PARAM_ID)"
//	@Failure		401	{object}	response.BaseSwaggerEmptyResponse{}									"Unauthorized. Available code (UNAUTHORIZED)"
//	@Failure		422	{object}	response.BaseSwaggerEmptyResponse{}									"Tenant not found. Available code (TENANT_NOT_FOUND)"
//	@Failure		500	{object}	response.BaseSwaggerEmptyResponse{}									"Internal Server Error. Available code (INTERNAL_SERVER_ERROR)"
//	@Router			/tenant/{id} [get]
func (h tenantHandler) Detail(c fiber.Ctx) error {
	id, err := request.GetParamID(c, "id")
	if err != nil {
		return response.New(c, code.InvalidParamID)
	}

	item, err := h.useCase.Detail(c.Context(), id)
	if err != nil {
		return err
	}

	return response.New(c, code.DetailFetched, response.WithData(response.NewDetailTenantResponse(item)))
}

// Create godoc
//
//	@Summary		Create Tenant
//	@Description	Create a new tenant. Validation: name REQUIRED, ALPHASPACE, MIN 3, MAX 32; logo REQUIRED, STORAGE tenant-logo
//	@ID				tenant-create
//	@Tags			Tenant
//	@Accept			json
//	@Produce		json
//	@Security		CookieAccessToken
//	@Param			payload	body		request.CreateTenantRequest					true	"Create Tenant Payload"
//	@Success		201		{object}	response.BaseSwaggerEmptyResponse{}			"Tenant created successfully. Available code (CREATED)"
//	@Failure		400		{object}	response.BaseSwaggerValidationResponse{}	"Bad Request - invalid input data. Available code (VALIDATION_ERROR, BAD_REQUEST)"
//	@Failure		401		{object}	response.BaseSwaggerEmptyResponse{}			"Unauthorized. Available code (UNAUTHORIZED)"
//	@Failure		500		{object}	response.BaseSwaggerEmptyResponse{}			"Internal Server Error. Available code (INTERNAL_SERVER_ERROR)"
//	@Router			/tenant [post]
func (h tenantHandler) Create(c fiber.Ctx) error {
	payload, err := helper.BindJSON[request.CreateTenantRequest](c)
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
//	@Summary		Update Tenant
//	@Description	Update tenant name and logo. Validation: name REQUIRED, ALPHASPACE, MIN 3, MAX 32; logo STORAGE tenant-logo (optional)
//	@ID				tenant-update
//	@Tags			Tenant
//	@Accept			json
//	@Produce		json
//	@Security		CookieAccessToken
//	@Param			id		path		string										true	"Tenant ID"
//	@Param			payload	body		request.UpdateTenantRequest					true	"Update Tenant Payload"
//	@Success		200		{object}	response.BaseSwaggerEmptyResponse{}			"Tenant updated successfully. Available code (UPDATED)"
//	@Failure		400		{object}	response.BaseSwaggerValidationResponse{}	"Bad Request - invalid input data. Available code (VALIDATION_ERROR, BAD_REQUEST)"
//	@Failure		401		{object}	response.BaseSwaggerEmptyResponse{}			"Unauthorized. Available code (UNAUTHORIZED)"
//	@Failure		422		{object}	response.BaseSwaggerEmptyResponse{}			"Tenant not found. Available code (TENANT_NOT_FOUND)"
//	@Failure		500		{object}	response.BaseSwaggerEmptyResponse{}			"Internal Server Error. Available code (INTERNAL_SERVER_ERROR)"
//	@Router			/tenant/{id} [put]
func (h tenantHandler) Update(c fiber.Ctx) error {
	id, err := request.GetParamID(c, "id")
	if err != nil {
		return response.New(c, code.InvalidParamID)
	}

	payload, err := helper.BindJSON[request.UpdateTenantRequest](c)
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

// UpdateStatus godoc
//
//	@Summary		Update Tenant Status
//	@Description	Update tenant status (active, inactive, suspended). Validation: status must be one of active, inactive, suspended
//	@ID				tenant-update-status
//	@Tags			Tenant
//	@Produce		json
//	@Security		CookieAccessToken
//	@Param			id		path		string								true	"Tenant ID"
//	@Param			status	path		string								true	"Status"	Enums(active, inactive, suspended)
//	@Success		200		{object}	response.BaseSwaggerEmptyResponse{}	"Tenant status updated successfully. Available code (UPDATED)"
//	@Failure		400		{object}	response.BaseSwaggerEmptyResponse{}	"Bad Request - invalid status. Available code (INVALID_PATH, VALIDATION_ERROR)"
//	@Failure		401		{object}	response.BaseSwaggerEmptyResponse{}	"Unauthorized. Available code (UNAUTHORIZED)"
//	@Failure		422		{object}	response.BaseSwaggerEmptyResponse{}	"Tenant has same status. Available code (TENANT_HAS_SAME_STATUS)"
//	@Failure		500		{object}	response.BaseSwaggerEmptyResponse{}	"Internal Server Error. Available code (INTERNAL_SERVER_ERROR)"
//	@Router			/tenant/{id}/status/{status} [post]
func (h tenantHandler) UpdateStatus(c fiber.Ctx) error {
	id, err := request.GetParamID(c, "id")
	if err != nil {
		return response.New(c, code.InvalidParamID)
	}

	status := entity.TenantStatus(c.Params("status"))
	if !status.Valid() {
		return domainError.New(code.InvalidPath, domainError.WithParams(map[string]any{
			"value": status.String(),
		}))
	}

	if err := h.useCase.UpdateStatus(c.Context(), id, status); err != nil {
		return err
	}

	return response.New(c, code.Updated)
}
