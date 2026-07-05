package code

import "github.com/gofiber/fiber/v3"

var (
	UserHasSameStatus = AppCode{"USER_HAS_SAME_STATUS", fiber.StatusUnprocessableEntity}
)
