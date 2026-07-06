package handler

import (
	"github.com/gofiber/fiber/v3"
	"github.com/rizalarfiyan/be-plant-factory/internal/delivery/http/request"
	"github.com/rizalarfiyan/be-plant-factory/internal/delivery/http/response"
	"github.com/rizalarfiyan/be-plant-factory/internal/shared/code"
	"github.com/rizalarfiyan/be-plant-factory/internal/shared/helper"
	"github.com/rizalarfiyan/be-plant-factory/internal/usecase/storage"
	"github.com/samber/do/v2"
)

type StorageHandler interface {
	PresignedUpload(c fiber.Ctx) error
}

type storageHandler struct {
	useCase storage.StorageUseCase `do:""`
}

func NewStorageHandler(i do.Injector) (StorageHandler, error) {
	return do.InvokeStruct[storageHandler](i)
}

// PresignedUpload godoc
//
//	@Summary		Generate Presigned Upload
//	@Description	Generate presigned Upload for file upload. Validation: mime-type REQUIRED, INVALID; filesize REQUIRED, MAX_FILE_SIZE_EXCEEDED, MIN_FILE_SIZE_EXCEEDED, GT 0; extensions REQUIRED, INVALID; type REQUIRED, must be one of (avatar)
//	@ID				storage-presigned-upload
//	@Tags			Storage
//	@Accept			json
//	@Produce		json
//	@Security		CookieAccessToken
//	@Param			payload	body		request.StoragePresignedUploadRequest										true	"Presigned Upload Payload"
//	@Success		200		{object}	response.BaseSwaggerResponse{data=response.StoragePresignedUploadResponse}	"Presigned Upload generated successfully. Available code (OK)"
//	@Failure		400		{object}	response.BaseSwaggerValidationResponse{}									"Bad Request - invalid input data. Available code (VALIDATION_ERROR, BAD_REQUEST)"
//	@Failure		401		{object}	response.BaseSwaggerEmptyResponse{}											"Unauthorized. Available code (UNAUTHORIZED)"
//	@Failure		500		{object}	response.BaseSwaggerEmptyResponse{}											"Internal Server Error. Available code (INTERNAL_SERVER_ERROR)"
//	@Router			/storage/presigned-upload [post]
func (h storageHandler) PresignedUpload(c fiber.Ctx) error {
	payload, err := helper.BindJSON[request.StoragePresignedUploadRequest](c)
	if err != nil {
		return response.NewValidate(c, err)
	}

	req := payload.ToEntity()
	item, err := h.useCase.GeneratePresignedUpload(c.Context(), req)
	if err != nil {
		return err
	}

	res := response.NewStoragePresignedUploadResponse(item)
	return response.New(c, code.OK, response.WithData(res))
}
