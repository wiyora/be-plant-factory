package domainError

import (
	"errors"

	"github.com/rizalarfiyan/be-plant-factory/internal/shared/code"
)

type AppError struct {
	Code   code.AppCode
	Params map[string]any
	Data   any
}

type Option func(*AppError)

func WithParams(params map[string]any) Option {
	return func(e *AppError) {
		e.Params = params
	}
}

func WithData(data any) Option {
	return func(e *AppError) {
		e.Data = data
	}
}

func New(appCode code.AppCode, opts ...Option) *AppError {
	e := &AppError{
		Code: appCode,
	}

	for _, opt := range opts {
		opt(e)
	}

	return e
}

func (e *AppError) Error() string {
	return string(e.Code.Code)
}

func IsAppError(err error) (*AppError, bool) {
	var errs *AppError
	if ok := errors.As(err, &errs); ok {
		return errs, true
	}

	return nil, false
}
