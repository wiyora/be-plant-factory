package domainError

import "errors"

var (
	ErrValidationStateNotSet = errors.New("VALIDATION_STATE_NOT_SET")
)
