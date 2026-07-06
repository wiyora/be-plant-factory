package code

import "github.com/gofiber/fiber/v3"

var (
	RoleNameExists = AppCode{"ROLE_NAME_EXISTS", fiber.StatusUnprocessableEntity}
	RoleNotFound   = AppCode{"ROLE_NOT_FOUND", fiber.StatusUnprocessableEntity}
	RoleInUse      = AppCode{"ROLE_IN_USE", fiber.StatusUnprocessableEntity}
)
