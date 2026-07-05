package search

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/rizalarfiyan/be-plant-factory/internal/domain/entity"
	domainError "github.com/rizalarfiyan/be-plant-factory/internal/domain/error"
	"github.com/rizalarfiyan/be-plant-factory/internal/shared/code"
)

type Parser func(c fiber.Ctx) (entity.Search, error)

type config struct {
	keySearch string
}

type Option func(*config)

func WithKey(searchKey string) Option {
	return func(c *config) {
		if searchKey != "" {
			c.keySearch = searchKey
		}
	}
}

func Parse(opts ...Option) Parser {
	cfg := config{
		keySearch: "search",
	}

	for _, opt := range opts {
		opt(&cfg)
	}

	return func(c fiber.Ctx) (entity.Search, error) {
		val := strings.TrimSpace(c.Query(cfg.keySearch))
		if val == "" {
			return entity.Search(""), nil
		}

		if len([]rune(val)) < 3 {
			return entity.Search(""), domainError.New(code.InvalidMinimalSearchQuery, domainError.WithParams(map[string]any{
				"value": val,
				"key":   cfg.keySearch,
			}))
		}

		return entity.Search(val), nil
	}
}
