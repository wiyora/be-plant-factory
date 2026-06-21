package repository

import (
	"github.com/rizalarfiyan/be-plant-factory/internal/repository/postgres"
	"github.com/samber/do/v2"
)

var Package = do.Package(
	// Generic Manager
	do.Lazy(postgres.NewTransactionManager),

	// Postgres repositories

	// Redis repositories

	// Cache repositories
)
