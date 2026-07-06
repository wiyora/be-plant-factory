package domainError

import "errors"

var (
	ErrValidationStateNotSet = errors.New("VALIDATION_STATE_NOT_SET")
	ErrStorageDiffInvalid    = errors.New("STORAGE_DIFF_INVALID")
)
