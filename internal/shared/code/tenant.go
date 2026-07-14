package code

import "github.com/gofiber/fiber/v3"

var (
	TenantHasSameStatus = AppCode{"TENANT_HAS_SAME_STATUS", fiber.StatusUnprocessableEntity}
	TenantNotFound      = AppCode{"TENANT_NOT_FOUND", fiber.StatusUnprocessableEntity}
)
