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
	usertenant "github.com/rizalarfiyan/be-plant-factory/internal/usecase/user-tenant"
	"github.com/samber/do/v2"
)

type UserTenantHandler interface {
	ListTenantUsers(c fiber.Ctx) error
	AssignUserRole(c fiber.Ctx) error
	RemoveTenantUser(c fiber.Ctx) error
	ListUserTenants(c fiber.Ctx) error
}

type userTenantHandler struct {
	useCase usertenant.UserTenantUseCase `do:""`
}

func NewUserTenantHandler(i do.Injector) (UserTenantHandler, error) {
	return do.InvokeStruct[userTenantHandler](i)
}

// ListTenantUsers godoc
//
//	@Summary		List Tenant Users
//	@Description	Get list of users in a tenant with pagination and search
//	@ID				tenant-user-list
//	@Tags			Tenant User
//	@Produce		json
//	@Security		CookieAccessToken
//	@Param			id			path		string																			true	"Tenant ID"
//	@Param			page		query		int																				false	"Page"		default(1)
//	@Param			page_size	query		int																				false	"Page Size"	default(10)
//	@Param			search		query		string																			false	"Search by user name"
//	@Param			order_by	query		string																			false	"Order by"	Enums(id, created_at, name, email)	default(id)
//	@Param			sort_by		query		string																			false	"Sort by"	Enums(asc, desc)					default(asc)
//	@Success		200			{object}	response.BasePaginationSwaggerResponse{data=[]response.ListTenantUserResponse}	"List fetched successfully. Available code (LIST_FETCHED)"
//	@Failure		400			{object}	response.BaseSwaggerValidationResponse{}										"Bad Request. Available code (VALIDATION_ERROR, BAD_REQUEST)"
//	@Failure		401			{object}	response.BaseSwaggerEmptyResponse{}												"Unauthorized. Available code (UNAUTHORIZED)"
//	@Failure		500			{object}	response.BaseSwaggerEmptyResponse{}												"Internal Server Error. Available code (INTERNAL_SERVER_ERROR)"
//	@Router			/tenant/{id}/users [get]
func (h userTenantHandler) ListTenantUsers(c fiber.Ctx) error {
	tenantID, err := request.GetParamID(c, "id")
	if err != nil {
		return response.New(c, code.InvalidParamID)
	}

	req := entity.TenantUserFilter{
		TenantID: tenantID,
	}

	err = helper.QueryBind(c,
		helper.QueryField(&req.Search, search.Parse()),
		helper.QueryField(&req.Pagination, pagination.Parse()),
		helper.QueryField(&req.Order, sorting.Parse(
			sorting.WithDefaults("id", entity.Asc),
			sorting.WithAllowedOrderBy("id", "created_at", "name", "email"),
			sorting.WithMappingOrderBy(map[string]string{
				"id":         "id",
				"created_at": "created_at",
				"name":       "name",
				"email":      "email",
			}),
		)),
	)
	if err != nil {
		return err
	}

	items, total, err := h.useCase.ListByTenant(c.Context(), req)
	if err != nil {
		return err
	}

	res := make([]response.ListTenantUserResponse, len(items))
	for i, item := range items {
		res[i] = response.NewListTenantUserResponse(item)
	}

	return response.New(c, code.ListFetched, response.WithData(res), response.WithPagination(req.Pagination.ToResult(total)))
}

// AssignUserRole godoc
//
//	@Summary		Assign or Update User Role in Tenant
//	@Description	Assign or update user role in tenant
//	@ID				tenant-user-assign-role
//	@Tags			Tenant User
//	@Accept			json
//	@Produce		json
//	@Security		CookieAccessToken
//	@Param			id		path		string										true	"Tenant ID"
//	@Param			user_id	path		string										true	"User ID"
//	@Param			payload	body		request.AssignUserRoleRequest				true	"Assign User Role Payload"
//	@Success		200		{object}	response.BaseSwaggerEmptyResponse{}			"User role assigned successfully. Available code (UPDATED)"
//	@Failure		400		{object}	response.BaseSwaggerValidationResponse{}	"Bad Request. Available code (VALIDATION_ERROR, BAD_REQUEST)"
//	@Failure		401		{object}	response.BaseSwaggerEmptyResponse{}			"Unauthorized. Available code (UNAUTHORIZED)"
//	@Failure		500		{object}	response.BaseSwaggerEmptyResponse{}			"Internal Server Error. Available code (INTERNAL_SERVER_ERROR)"
//	@Router			/tenant/{id}/users/{user_id} [put]
func (h userTenantHandler) AssignUserRole(c fiber.Ctx) error {
	tenantID, err := request.GetParamID(c, "id")
	if err != nil {
		return response.New(c, code.InvalidParamID)
	}

	userID, err := request.GetParamID(c, "user_id")
	if err != nil {
		return response.New(c, code.InvalidParamID)
	}

	payload, err := helper.BindJSON[request.AssignUserRoleRequest](c)
	if err != nil {
		return response.NewValidate(c, err)
	}

	if err := h.useCase.AssignRole(c.Context(), tenantID, userID, payload.RoleID); err != nil {
		return err
	}

	return response.New(c, code.Updated)
}

// RemoveTenantUser godoc
//
//	@Summary		Remove User from Tenant
//	@Description	Remove user from tenant
//	@ID				tenant-user-remove
//	@Tags			Tenant User
//	@Produce		json
//	@Security		CookieAccessToken
//	@Param			id		path		string								true	"Tenant ID"
//	@Param			user_id	path		string								true	"User ID"
//	@Success		200		{object}	response.BaseSwaggerEmptyResponse{}	"User removed successfully. Available code (DELETED)"
//	@Failure		400		{object}	response.BaseSwaggerEmptyResponse{}	"Bad Request. Available code (INVALID_PARAM_ID)"
//	@Failure		401		{object}	response.BaseSwaggerEmptyResponse{}	"Unauthorized. Available code (UNAUTHORIZED)"
//	@Failure		500		{object}	response.BaseSwaggerEmptyResponse{}	"Internal Server Error. Available code (INTERNAL_SERVER_ERROR)"
//	@Router			/tenant/{id}/users/{user_id} [delete]
func (h userTenantHandler) RemoveTenantUser(c fiber.Ctx) error {
	tenantID, err := request.GetParamID(c, "id")
	if err != nil {
		return response.New(c, code.InvalidParamID)
	}

	userID, err := request.GetParamID(c, "user_id")
	if err != nil {
		return response.New(c, code.InvalidParamID)
	}

	if err := h.useCase.RemoveUser(c.Context(), tenantID, userID); err != nil {
		return err
	}

	return response.New(c, code.Deleted)
}

// ListUserTenants godoc
//
//	@Summary		List User Tenants
//	@Description	Get list of tenants for a user with pagination and filters
//	@ID				user-tenant-list
//	@Tags			User Tenant
//	@Produce		json
//	@Security		CookieAccessToken
//	@Param			id			path		string																			true	"User ID"
//	@Param			page		query		int																				false	"Page"									default(1)
//	@Param			page_size	query		int																				false	"Page Size"								default(10)
//	@Param			order_by	query		string																			false	"Order by"								Enums(id, created_at, name, email)	default(id)
//	@Param			sort_by		query		string																			false	"Sort by"								Enums(asc, desc)					default(asc)
//	@Param			tenant_ids	query		[]string																		false	"Tenant IDs (comma-separated UUIDs)"	collectionFormat(csv)
//	@Param			role_ids	query		[]string																		false	"Role IDs (comma-separated UUIDs)"		collectionFormat(csv)
//	@Success		200			{object}	response.BasePaginationSwaggerResponse{data=[]response.ListUserTenantResponse}	"List fetched successfully. Available code (LIST_FETCHED)"
//	@Failure		400			{object}	response.BaseSwaggerValidationResponse{}										"Bad Request. Available code (VALIDATION_ERROR, BAD_REQUEST)"
//	@Failure		401			{object}	response.BaseSwaggerEmptyResponse{}												"Unauthorized. Available code (UNAUTHORIZED)"
//	@Failure		500			{object}	response.BaseSwaggerEmptyResponse{}												"Internal Server Error. Available code (INTERNAL_SERVER_ERROR)"
//	@Router			/user/{id}/tenants [get]
func (h userTenantHandler) ListUserTenants(c fiber.Ctx) error {
	userID, err := request.GetParamID(c, "id")
	if err != nil {
		return response.New(c, code.InvalidParamID)
	}

	req := entity.UserTenantFilter{
		UserID: userID,
	}

	err = helper.QueryBind(c,
		helper.QueryField(&req.Pagination, pagination.Parse()),
		helper.QueryField(&req.Order, sorting.Parse(
			sorting.WithDefaults("id", entity.Asc),
			sorting.WithAllowedOrderBy("id", "created_at", "name", "email"),
			sorting.WithMappingOrderBy(map[string]string{
				"id":         "id",
				"created_at": "created_at",
				"name":       "name",
				"email":      "email",
			}),
		)),
		helper.QueryField(&req.TenantIDs, filter.ParseUUIDs("tenant_ids")),
		helper.QueryField(&req.RoleIDs, filter.ParseUUIDs("role_ids")),
	)
	if err != nil {
		return err
	}

	items, total, err := h.useCase.ListByUser(c.Context(), req)
	if err != nil {
		return err
	}

	res := make([]response.ListUserTenantResponse, len(items))
	for i, item := range items {
		res[i] = response.NewListUserTenantResponse(item)
	}

	return response.New(c, code.ListFetched, response.WithData(res), response.WithPagination(req.Pagination.ToResult(total)))
}
