package helper

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

type PgxErrorMapper interface {
	Register(constraint string, mappedErr error)
	MapError(err error, onUnhandled func(err error)) error
}

type pgxErrorMapper struct {
	constraints map[string]error
}

func NewPgxErrorMapper() PgxErrorMapper {
	return &pgxErrorMapper{
		constraints: make(map[string]error),
	}
}

func (m *pgxErrorMapper) Register(constraint string, mappedErr error) {
	m.constraints[constraint] = mappedErr
}

func (m *pgxErrorMapper) MapError(err error, onUnhandled func(err error)) error {
	if err == nil {
		return nil
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if mappedErr, ok := m.constraints[pgErr.ConstraintName]; ok {
			return mappedErr
		}
	}

	if onUnhandled != nil {
		onUnhandled(err)
	}

	return err
}
