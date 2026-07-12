package response

type BaseSwaggerResponse struct {
	Message MessageSwaggerResponse `json:"message"`
	Data    any                    `json:"data"`
}

type BaseSwaggerEmptyResponse struct {
	Message MessageSwaggerResponse `json:"message"`
	Data    any                    `json:"data" example:"null" extensions:"x-nullable"`
}

type MessageSwaggerResponse struct {
	Code   string         `json:"code" example:"APP_CODE"`
	Params map[string]any `json:"params,omitempty" swaggertype:"object,string" example:"key:value"`
}

type BaseSwaggerValidationResponse struct {
	Message MessageSwaggerValidationResponse `json:"message"`
	Data    any                              `json:"data" example:"null" extensions:"x-nullable"`
}

type MessageSwaggerValidationResponse struct {
	Code   string                                         `json:"code" example:"VALIDATION_ERROR"`
	Fields map[string]ValidationErrorFieldSwaggerResponse `json:"fields,omitempty"`
}

type ValidationErrorFieldSwaggerResponse struct {
	Code   string         `json:"code" example:"VALIDATION_CODE"`
	Params map[string]any `json:"params,omitempty" swaggertype:"object,string" example:"key:value"`
}

type BasePaginationSwaggerResponse struct {
	Message    MessageSwaggerResponse    `json:"message"`
	Data       any                       `json:"data"`
	Pagination PaginationSwaggerResponse `json:"pagination"`
}

type PaginationSwaggerResponse struct {
	Page       uint64 `json:"page" example:"1"`
	PageSize   uint64 `json:"page_size" example:"10"`
	Total      uint64 `json:"total" example:"100"`
	TotalPages uint64 `json:"total_pages" example:"10"`
}
