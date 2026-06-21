package code

import "github.com/gofiber/fiber/v3"

var (
	OK                  = AppCode{"SUCCESS", fiber.StatusOK}
	Created             = AppCode{"CREATED", fiber.StatusCreated}
	Accepted            = AppCode{"ACCEPTED", fiber.StatusAccepted}
	NoContent           = AppCode{"NO_CONTENT", fiber.StatusNoContent}
	BadRequest          = AppCode{"BAD_REQUEST", fiber.StatusBadRequest}
	Unauthorized        = AppCode{"UNAUTHORIZED", fiber.StatusUnauthorized}
	Forbidden           = AppCode{"FORBIDDEN", fiber.StatusForbidden}
	NotFound            = AppCode{"NOT_FOUND", fiber.StatusNotFound}
	Conflict            = AppCode{"CONFLICT", fiber.StatusConflict}
	InternalServerError = AppCode{"INTERNAL_SERVER_ERROR", fiber.StatusInternalServerError}
)

var HttpCodeAppCode = map[int]AppCode{
	fiber.StatusOK:                  OK,
	fiber.StatusCreated:             Created,
	fiber.StatusAccepted:            Accepted,
	fiber.StatusNoContent:           NoContent,
	fiber.StatusBadRequest:          BadRequest,
	fiber.StatusUnauthorized:        Unauthorized,
	fiber.StatusForbidden:           Forbidden,
	fiber.StatusNotFound:            NotFound,
	fiber.StatusConflict:            Conflict,
	fiber.StatusInternalServerError: InternalServerError,
}
