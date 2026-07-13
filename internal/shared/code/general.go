package code

import "github.com/gofiber/fiber/v3"

var (
	InvalidBodyRequest = AppCode{"INVALID_BODY_REQUEST", fiber.StatusBadRequest}
	ListFetched        = AppCode{"LIST_FETCHED", fiber.StatusOK}
	DetailFetched      = AppCode{"DETAIL_FETCHED", fiber.StatusOK}
	Updated            = AppCode{"UPDATED", fiber.StatusOK}
	Deleted            = AppCode{"DELETED", fiber.StatusOK}
	Restored           = AppCode{"RESTORED", fiber.StatusOK}
	ValidationError    = AppCode{"VALIDATION_ERROR", fiber.StatusBadRequest}
	RateLimitExceeded  = AppCode{"RATE_LIMIT_EXCEEDED", fiber.StatusTooManyRequests}

	InvalidPath               = AppCode{"INVALID_PATH", fiber.StatusBadRequest}
	InvalidPageQuery          = AppCode{"INVALID_PAGE_QUERY", fiber.StatusBadRequest}
	InvalidPageSizeQuery      = AppCode{"INVALID_PAGE_SIZE_QUERY", fiber.StatusBadRequest}
	InvalidOrderByQuery       = AppCode{"INVALID_ORDER_BY_QUERY", fiber.StatusBadRequest}
	InvalidSortByQuery        = AppCode{"INVALID_SORT_BY_QUERY", fiber.StatusBadRequest}
	InvalidMinimalSearchQuery = AppCode{"INVALID_MINIMAL_SEARCH_QUERY", fiber.StatusBadRequest}
	InvalidFilterEnumQuery    = AppCode{"INVALID_FILTER_ENUM_QUERY", fiber.StatusBadRequest}
	InvalidFilterUUIDQuery    = AppCode{"INVALID_FILTER_UUID_QUERY", fiber.StatusBadRequest}
	InvalidFilterQuery        = AppCode{"INVALID_FILTER_QUERY", fiber.StatusBadRequest}

	InvalidParamID        = AppCode{"INVALID_PARAM_ID", fiber.StatusBadRequest}
	InvalidStartDateQuery = AppCode{"INVALID_START_DATE_QUERY", fiber.StatusBadRequest}
	InvalidEndDateQuery   = AppCode{"INVALID_END_DATE_QUERY", fiber.StatusBadRequest}
	InvalidDateRangeQuery = AppCode{"INVALID_DATE_RANGE_QUERY", fiber.StatusBadRequest}
)
