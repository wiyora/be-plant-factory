package response

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/rizalarfiyan/be-plant-factory/internal/domain/entity"
	domainError "github.com/rizalarfiyan/be-plant-factory/internal/domain/error"
	"github.com/rizalarfiyan/be-plant-factory/internal/shared/code"
	"github.com/rizalarfiyan/be-plant-factory/internal/shared/helper"
)

type BaseResponse struct {
	Message    MessageResponse     `json:"message"`
	Data       any                 `json:"data"`
	Pagination *PaginationResponse `json:"pagination,omitempty"`
}

type ParamsResponse = map[string]any

type FieldsResponse = map[string]CodeParamsResponse

type CodeParamsResponse struct {
	Code   string         `json:"code"`
	Params ParamsResponse `json:"params,omitempty"`
}

type MessageResponse struct {
	CodeParamsResponse
	Fields FieldsResponse `json:"fields,omitempty"`
}

type PaginationResponse struct {
	Page       uint64 `json:"page"`
	PageSize   uint64 `json:"page_size"`
	Total      uint64 `json:"total"`
	TotalPages uint64 `json:"total_pages"`
}

type responseConfig struct {
	res      BaseResponse
	httpCode int
}

type Option func(*responseConfig)

func WithParams(params ParamsResponse) Option {
	return func(r *responseConfig) {
		r.res.Message.Params = params
	}
}

func WithFields(fields FieldsResponse) Option {
	return func(r *responseConfig) {
		r.res.Message.Fields = fields
	}
}

func WithData(data any) Option {
	return func(r *responseConfig) {
		r.res.Data = data
	}
}

func WithPagination(pagination entity.PaginationResult) Option {
	return func(r *responseConfig) {
		r.res.Pagination = &PaginationResponse{
			Page:       pagination.Page,
			PageSize:   pagination.PageSize,
			Total:      pagination.Total,
			TotalPages: pagination.TotalPages,
		}
	}
}

func New(c fiber.Ctx, code code.AppCode, opts ...Option) error {
	r := &responseConfig{
		res: BaseResponse{
			Message: MessageResponse{
				CodeParamsResponse: CodeParamsResponse{
					Code: code.Code,
				},
			},
		},
		httpCode: code.HttpCode,
	}

	for _, opt := range opts {
		opt(r)
	}

	return c.Status(r.httpCode).JSON(r.res)
}

func NewValidate(c fiber.Ctx, err error) error {
	if manualErr, ok := domainError.IsManualValidationError(err); ok {
		return manualErr
	}

	mapSkipParamCode := map[string]struct{}{
		"db_exists": {},
	}

	var errs validator.ValidationErrors
	if errors.As(err, &errs) {
		fields := make(FieldsResponse, len(errs))
		for _, val := range errs {
			code := helper.ToSnakeCase(val.Tag())

			var params = make(ParamsResponse)
			if _, skip := mapSkipParamCode[code]; !skip && val.Param() != "" {
				params["value"] = val.Value()
			}

			fields[val.Field()] = CodeParamsResponse{
				Code:   strings.ToUpper(code),
				Params: params,
			}
		}

		return New(c, code.ValidationError, WithFields(fields))
	}

	return New(c, code.InvalidBodyRequest)
}
