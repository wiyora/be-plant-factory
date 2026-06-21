package code

import "github.com/gofiber/fiber/v3"

var (
	HealthCheckSuccess = AppCode{"HEALTH_CHECK_SUCCESS", fiber.StatusOK}
)
