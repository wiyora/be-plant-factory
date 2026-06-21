// internal/repository/postgres/tx_manager.go
package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rizalarfiyan/be-plant-factory/internal/infrastructure/postgres"
	"github.com/samber/do/v2"
)

type transactionManager struct {
	db *pgxpool.Pool
}

type TransactionManager interface {
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

func NewTransactionManager(i do.Injector) (TransactionManager, error) {
	db := do.MustInvoke[*postgres.Database](i)
	return transactionManager{db: db.Pool()}, nil
}

func (tm transactionManager) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := tm.db.Begin(ctx)
	if err != nil {
		return err
	}

	txCtx := context.WithValue(ctx, txKey, tx)
	err = fn(txCtx)
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}
