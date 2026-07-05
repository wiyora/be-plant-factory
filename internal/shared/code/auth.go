package code

import "github.com/gofiber/fiber/v3"

var (
	InvalidState            = AppCode{"INVALID_STATE", fiber.StatusBadRequest}
	ProviderError           = AppCode{"PROVIDER_ERROR", fiber.StatusInternalServerError}
	TokenExchangeFailed     = AppCode{"TOKEN_EXCHANGE_FAILED", fiber.StatusBadRequest}
	InvalidCallback         = AppCode{"INVALID_CALLBACK", fiber.StatusBadRequest}
	FailedCallback          = AppCode{"FAILED_CALLBACK", fiber.StatusUnprocessableEntity}
	InvalidRefreshToken     = AppCode{"INVALID_REFRESH_TOKEN", fiber.StatusUnauthorized}
	RefreshTokenNotEligible = AppCode{"REFRESH_TOKEN_NOT_ELIGIBLE", fiber.StatusUnauthorized}
	InvalidRedirectURL      = AppCode{"INVALID_REDIRECT_URL", fiber.StatusBadRequest}
	UserNotFound            = AppCode{"USER_NOT_FOUND", fiber.StatusUnprocessableEntity}
	UserStatusInactive      = AppCode{"USER_STATUS_INACTIVE", fiber.StatusUnauthorized}
	UserStatusBanned        = AppCode{"USER_STATUS_BANNED", fiber.StatusUnauthorized}
	InvalidCurrentStep      = AppCode{"INVALID_CURRENT_STEP", fiber.StatusUnprocessableEntity}
)
