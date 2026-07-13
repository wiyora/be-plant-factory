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
	keyPage         string
	keyPageSize     string
	defaultPage     uint64
	defaultPageSize uint64
	maxPageSize     uint64
}

func WithKeys(pageKey, pageSizeKey string) Option {
	return func(c *config) {
		if pageKey != "" {
			c.keyPage = pageKey
		}
		if pageSizeKey != "" {
			c.keyPageSize = pageSizeKey
		}
	}
}

func WithDefaults(page, pageSize uint64) Option {
	return func(c *config) {
		if page > 0 {
			c.defaultPage = page
		}
		if pageSize > 0 {
			c.defaultPageSize = pageSize
		}
	}
}

func WithMaxPageSize(max uint64) Option {
	return func(c *config) {
		if max > 0 {
			c.maxPageSize = max
		}
	}
}

func Parse(opts ...Option) Parser {
	cfg := config{
		keyPage:         "page",
		keyPageSize:     "page_size",
		defaultPage:     1,
		defaultPageSize: 10,
		maxPageSize:     50,
	}

	for _, opt := range opts {
		opt(&cfg)
	}

	return func(c fiber.Ctx) (entity.Pagination, error) {
		page := fiber.Query(c, cfg.keyPage, cfg.defaultPage)
		pageSize := fiber.Query(c, cfg.keyPageSize, cfg.defaultPageSize)

		if page <= 0 {
			return entity.Pagination{}, domainError.New(code.InvalidPageQuery, domainError.WithParams(map[string]any{
				"value": page,
				"key":   cfg.keyPage,
			}))
		}

		if pageSize <= 0 || pageSize > cfg.maxPageSize {
			return entity.Pagination{}, domainError.New(code.InvalidPageSizeQuery, domainError.WithParams(map[string]any{
				"value": pageSize,
				"key":   cfg.keyPageSize,
			}))
		}

		return entity.Pagination{
			Page:     page,
			PageSize: pageSize,
		}, nil
	}
}
