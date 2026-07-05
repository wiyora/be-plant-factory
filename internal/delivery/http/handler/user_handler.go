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
	domainError "github.com/rizalarfiyan/be-plant-factory/internal/domain/error"
	"github.com/rizalarfiyan/be-plant-factory/internal/shared/code"
	"github.com/rizalarfiyan/be-plant-factory/internal/shared/helper"
	"github.com/rizalarfiyan/be-plant-factory/internal/usecase/user"
	"github.com/samber/do/v2"
)

type UserHandler interface {
	List(c fiber.Ctx) error
	Detail(c fiber.Ctx) error
	Create(c fiber.Ctx) error
	Update(c fiber.Ctx) error
	UpdateStatus(c fiber.Ctx) error
}

type userHandler struct {
	useCase user.UserUseCase `do:""`
}

func NewUserHandler(i do.Injector) (UserHandler, error) {
	return do.InvokeStruct[userHandler](i)
}

// List godoc
//
//	@Summary		List Users
//	@Description	Get list of users with pagination and search
//	@ID				user-list
//	@Tags			User
//	@Produce		json
//	@Security		CookieAccessToken
//	@Param			search		query		string																		false	"Search by name (min 3 chars)"
//	@Param			page		query		int																			false	"Page number"		default(1)
//	@Param			limit		query		int																			false	"Items per page"	default(10)
//	@Param			order_by	query		string																		false	"Order by field"	default(id)		enums(id, created_at, name, email)
//	@Param			sort_by		query		string																		false	"Sort direction"	default(desc)	enums(asc, desc)
//	@Param			status		query		string																		false	"Filter by status"	enums(active, inactive, banned)
//	@Success		200			{object}	response.BasePaginationSwaggerResponse{data=[]response.ListUserResponse}	"Users fetched successfully"
//	@Failure		400			{object}	response.BaseSwaggerValidationResponse{}									"Bad Request"
//	@Failure		401			{object}	response.BaseSwaggerEmptyResponse{}											"Unauthorized"
//	@Failure		500			{object}	response.BaseSwaggerEmptyResponse{}											"Internal Server Error"
//	@Router			/user [get]
func (h userHandler) List(c fiber.Ctx) error {
	var req entity.UserFilter

	err := helper.QueryBind(c,
		helper.QueryField(&req.Search, search.Parse()),
		helper.QueryField(&req.Pagination, pagination.Parse()),
		helper.QueryField(&req.Order,
			sorting.Parse(
				sorting.WithDefaults("id", entity.Desc),
				sorting.WithAllowedOrderBy("id", "created_at", "name", "email"),
				sorting.WithMappingOrderBy(map[string]string{
					"id":         "id",
					"created_at": "id",
					"name":       "name",
					"email":      "email",
				}),
			)),
		helper.QueryField(&req.Status, filter.ParseEnum[entity.UserStatus]("status")),
	)
	if err != nil {
		return err
	}

	items, total, err := h.useCase.List(c.Context(), req)
	if err != nil {
		return err
	}

	res := make([]response.ListUserResponse, len(items))
	for i, item := range items {
		res[i] = response.NewListUserResponse(item)
	}

	return response.New(c, code.ListFetched, response.WithData(res), response.WithPagination(req.Pagination.ToResult(total)))
}

// Detail godoc
//
//	@Summary		Get User Detail
//	@Description	Get user detail by ID
//	@ID				user-detail
//	@Tags			User
//	@Produce		json
//	@Security		CookieAccessToken
//	@Param			id	path		string															true	"User ID"
//	@Success		200	{object}	response.BaseSwaggerResponse{data=response.DetailUserResponse}	"User detail fetched successfully"
//	@Failure		400	{object}	response.BaseSwaggerEmptyResponse{}								"Bad Request - invalid ID"
//	@Failure		401	{object}	response.BaseSwaggerEmptyResponse{}								"Unauthorized"
//	@Failure		422	{object}	response.BaseSwaggerEmptyResponse{}								"User not found"
//	@Failure		500	{object}	response.BaseSwaggerEmptyResponse{}								"Internal Server Error"
//	@Router			/user/{id} [get]
func (h userHandler) Detail(c fiber.Ctx) error {
	id, err := request.GetParamID(c, "id")
	if err != nil {
		return response.New(c, code.InvalidParamID)
	}

	user, err := h.useCase.Detail(c.Context(), id)
	if err != nil {
		return err
	}

	return response.New(c, code.DetailFetched, response.WithData(response.NewDetailUserResponse(user)))
}

// Create godoc
//
//	@Summary		Create User
//	@Description	Create a new user
//	@ID				user-create
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Security		CookieAccessToken
//	@Param			payload	body		request.CreateUserRequest					true	"Create User Payload"
//	@Success		201		{object}	response.BaseSwaggerEmptyResponse{}			"User created successfully"
//	@Failure		400		{object}	response.BaseSwaggerValidationResponse{}	"Bad Request - invalid input data"
//	@Failure		401		{object}	response.BaseSwaggerEmptyResponse{}			"Unauthorized"
//	@Failure		500		{object}	response.BaseSwaggerEmptyResponse{}			"Internal Server Error"
//	@Router			/user [post]
func (h userHandler) Create(c fiber.Ctx) error {
	payload, err := helper.BindJSON[request.CreateUserRequest](c)
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
//	@Summary		Update User
//	@Description	Update user name and avatar
//	@ID				user-update
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Security		CookieAccessToken
//	@Param			id		path		string										true	"User ID"
//	@Param			payload	body		request.UpdateUserRequest					true	"Update User Payload"
//	@Success		200		{object}	response.BaseSwaggerEmptyResponse{}			"User updated successfully"
//	@Failure		400		{object}	response.BaseSwaggerValidationResponse{}	"Bad Request - invalid input data"
//	@Failure		401		{object}	response.BaseSwaggerEmptyResponse{}			"Unauthorized"
//	@Failure		500		{object}	response.BaseSwaggerEmptyResponse{}			"Internal Server Error"
//	@Router			/user/{id} [put]
func (h userHandler) Update(c fiber.Ctx) error {
	id, err := request.GetParamID(c, "id")
	if err != nil {
		return response.New(c, code.InvalidParamID)
	}

	payload, err := helper.BindJSON[request.UpdateUserRequest](c)
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
//	@Summary		Update User Status
//	@Description	Update user status (active, inactive, banned)
//	@ID				user-update-status
//	@Tags			User
//	@Produce		json
//	@Security		CookieAccessToken
//	@Param			id		path		string								true	"User ID"
//	@Param			status	path		string								true	"Status"	Enums(active, inactive, banned)
//	@Success		200		{object}	response.BaseSwaggerEmptyResponse{}	"User status updated successfully"
//	@Failure		400		{object}	response.BaseSwaggerEmptyResponse{}	"Bad Request - invalid status"
//	@Failure		401		{object}	response.BaseSwaggerEmptyResponse{}	"Unauthorized"
//	@Failure		500		{object}	response.BaseSwaggerEmptyResponse{}	"Internal Server Error"
//	@Router			/user/{id}/status/{status} [post]
func (h userHandler) UpdateStatus(c fiber.Ctx) error {
	id, err := request.GetParamID(c, "id")
	if err != nil {
		return response.New(c, code.InvalidParamID)
	}

	status := entity.UserStatus(c.Params("status"))
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
