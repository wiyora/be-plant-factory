package handler

import (
	"github.com/gofiber/fiber/v3"
	"github.com/rizalarfiyan/be-plant-factory/internal/delivery/http/request"
	"github.com/rizalarfiyan/be-plant-factory/internal/delivery/http/response"
	"github.com/rizalarfiyan/be-plant-factory/internal/shared/code"
	"github.com/rizalarfiyan/be-plant-factory/internal/shared/helper"
	userMe "github.com/rizalarfiyan/be-plant-factory/internal/usecase/user-me"
	"github.com/samber/do/v2"
)

type UserMeHandler interface {
	GettingStarted(c fiber.Ctx) error
}

type userMeHandler struct {
	useCase userMe.UserMeUseCase `do:""`
}

func NewUserMeHandler(i do.Injector) (UserMeHandler, error) {
	return do.InvokeStruct[userMeHandler](i)
}

// GettingStarted godoc
//
//	@Summary		Getting Started
//	@Description	Complete the getting started process for the user. Validation:
//	@Description	name = REQUIRED, ALPHASPACE, MIN 3, MAX 64
//	@Description	avatar = REQUIRED, STARTSWITH "avatar:"
//	@ID				user-me-getting-started
//	@Tags			UserMe
//	@Accept			json
//	@Produce		json
//	@Security		CookieAccessToken
//	@Param			payload	body		request.UserMeGettingStartedRequest			true	"Getting Started Payload"
//	@Success		200		{object}	response.BaseSwaggerEmptyResponse{}			"Getting started process completed successfully. Available code (SUCCESS)"
//	@Failure		400		{object}	response.BaseSwaggerValidationResponse{}	"Bad Request - invalid input data. Available code (VALIDATION_ERROR, BAD_REQUEST)"
//	@Failure		422		{object}	response.BaseSwaggerEmptyResponse{}			"Unprocessable Entity - validation failed. Available code (INVALID_CURRENT_STEP)"
//	@Failure		401		{object}	response.BaseSwaggerEmptyResponse{}			"Unauthorized - user not authenticated. Available code (UNAUTHORIZED)"
//	@Failure		500		{object}	response.BaseSwaggerEmptyResponse{}			"Internal Server Error. Available code (INTERNAL_SERVER_ERROR)"
//	@Router			/user/me/getting-started [post]
func (h userMeHandler) GettingStarted(c fiber.Ctx) error {
	user, ok := helper.GetLocalUser(c)
	if !ok {
		return response.New(c, code.Unauthorized)
	}

	payload, err := helper.BindJSON[request.UserMeGettingStartedRequest](c)
	if err != nil {
		return response.NewValidate(c, err)
	}

	req := payload.ToEntity()
	req.UserID = user.ID
	req.CurrentStep = user.CurrentStep

	if err := h.useCase.GettingStarted(c.Context(), req); err != nil {
		return err
	}

	return response.New(c, code.OK)
}
