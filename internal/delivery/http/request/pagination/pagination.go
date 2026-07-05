package pagination

import (
	"github.com/gofiber/fiber/v3"
	"github.com/rizalarfiyan/be-plant-factory/internal/domain/entity"
	domainError "github.com/rizalarfiyan/be-plant-factory/internal/domain/error"
	"github.com/rizalarfiyan/be-plant-factory/internal/shared/code"
)

type Parser func(c fiber.Ctx) (entity.Pagination, error)

type Option func(*config)

type config struct {
	keyPage      string
	keyLimit     string
	defaultPage  uint64
	defaultLimit uint64
	maxLimit     uint64
}

func WithKeys(pageKey, limitKey string) Option {
	return func(c *config) {
		if pageKey != "" {
			c.keyPage = pageKey
		}
		if limitKey != "" {
			c.keyLimit = limitKey
		}
	}
}

func WithDefaults(page, limit uint64) Option {
	return func(c *config) {
		if page > 0 {
			c.defaultPage = page
		}
		if limit > 0 {
			c.defaultLimit = limit
		}
	}
}

func WithMaxLimit(max uint64) Option {
	return func(c *config) {
		if max > 0 {
			c.maxLimit = max
		}
	}
}

func Parse(opts ...Option) Parser {
	cfg := config{
		keyPage:      "page",
		keyLimit:     "limit",
		defaultPage:  1,
		defaultLimit: 10,
		maxLimit:     50,
	}

	for _, opt := range opts {
		opt(&cfg)
	}

	return func(c fiber.Ctx) (entity.Pagination, error) {
		page := fiber.Query(c, cfg.keyPage, cfg.defaultPage)
		limit := fiber.Query(c, cfg.keyLimit, cfg.defaultLimit)

		if page <= 0 {
			return entity.Pagination{}, domainError.New(code.InvalidPageQuery, domainError.WithParams(map[string]any{
				"value": page,
				"key":   cfg.keyPage,
			}))
		}

		if limit <= 0 || limit > cfg.maxLimit {
			return entity.Pagination{}, domainError.New(code.InvalidLimitQuery, domainError.WithParams(map[string]any{
				"value": limit,
				"key":   cfg.keyLimit,
			}))
		}

		return entity.Pagination{
			Page:  page,
			Limit: limit,
		}, nil
	}
}
