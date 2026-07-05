package sorting

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/rizalarfiyan/be-plant-factory/internal/domain/entity"
	domainError "github.com/rizalarfiyan/be-plant-factory/internal/domain/error"
	"github.com/rizalarfiyan/be-plant-factory/internal/shared/code"
)

type Parser func(c fiber.Ctx) (entity.Order, error)

type Option func(*config)

type config struct {
	keyOrderBy     string
	keySortBy      string
	defaultOrderBy string
	defaultSortBy  entity.SortDirection
	allowedOrderBy map[string]bool
	mappingOrderBy map[string]string
}

func WithKeys(orderByKey, sortByKey string) Option {
	return func(c *config) {
		if orderByKey != "" {
			c.keyOrderBy = orderByKey
		}
		if sortByKey != "" {
			c.keySortBy = sortByKey
		}
	}
}

func WithDefaults(orderBy string, sortBy entity.SortDirection) Option {
	return func(c *config) {
		if orderBy != "" {
			c.defaultOrderBy = orderBy
		}
		if sortBy.Valid() {
			c.defaultSortBy = sortBy
		}
	}
}

func WithAllowedOrderBy(columns ...string) Option {
	return func(c *config) {
		c.allowedOrderBy = make(map[string]bool)
		for _, col := range columns {
			c.allowedOrderBy[col] = true
		}
	}
}

func WithMappingOrderBy(mapping map[string]string) Option {
	return func(c *config) {
		c.mappingOrderBy = mapping
	}
}

func Parse(opts ...Option) Parser {
	cfg := config{
		keyOrderBy:     "order_by",
		keySortBy:      "sort_by",
		defaultOrderBy: "id",
		defaultSortBy:  entity.Desc,
	}

	for _, opt := range opts {
		opt(&cfg)
	}

	return func(c fiber.Ctx) (entity.Order, error) {
		orderByStr := c.Query(cfg.keyOrderBy)
		if orderByStr == "" {
			orderByStr = cfg.defaultOrderBy
		}

		if len(cfg.allowedOrderBy) > 0 && !cfg.allowedOrderBy[orderByStr] {
			return entity.Order{}, domainError.New(code.InvalidOrderByQuery, domainError.WithParams(map[string]any{
				"value": orderByStr,
				"key":   cfg.keyOrderBy,
			}))
		}

		sortByStr := c.Query(cfg.keySortBy)
		if sortByStr == "" {
			sortByStr = cfg.defaultSortBy.String()
		}

		sortBy := entity.SortDirection(strings.ToLower(sortByStr))
		if !sortBy.Valid() {
			return entity.Order{}, domainError.New(code.InvalidSortByQuery, domainError.WithParams(map[string]any{
				"value": sortByStr,
				"key":   cfg.keySortBy,
			}))
		}

		if len(cfg.mappingOrderBy) > 0 {
			if mappedValue, exists := cfg.mappingOrderBy[orderByStr]; exists {
				orderByStr = mappedValue
			}
		}

		return entity.Order{
			OrderBy: orderByStr,
			SortBy:  sortBy,
		}, nil
	}
}
