package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rizalarfiyan/be-plant-factory/internal/domain/entity"
	"github.com/rizalarfiyan/be-plant-factory/internal/infrastructure/logger"
	"github.com/rizalarfiyan/be-plant-factory/internal/infrastructure/postgres"
	"github.com/rizalarfiyan/be-plant-factory/internal/shared/helper"
	"github.com/samber/do/v2"
)

type SessionRepository interface {
	Create(ctx context.Context, session entity.Session) (uuid.UUID, error)
	Delete(ctx context.Context, userID, sessionID uuid.UUID) error
	GetByRefreshToken(ctx context.Context, refreshToken string) (Session, error)
	Update(ctx context.Context, sessionID uuid.UUID, refreshTokenHash string, expiredAt time.Time) error
}

type sessionRepository struct {
	db     *pgxpool.Pool
	mapper helper.PgxErrorMapper
}

func NewSessionRepository(i do.Injector) (SessionRepository, error) {
	db := do.MustInvoke[*postgres.Database](i)
	mapper := helper.NewPgxErrorMapper()

	return &sessionRepository{
		db:     db.Pool(),
		mapper: mapper,
	}, nil
}

func (r sessionRepository) Create(ctx context.Context, session entity.Session) (uuid.UUID, error) {
	log := logger.WithLayerCtx(ctx, logger.LayerPostgresRepository)

	query := `
		INSERT INTO sessions (user_id, refresh_token_hash, device_name, ip_address, expired_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`

	db := GetDB(ctx, r.db)
	row := db.QueryRow(ctx, query, session.UserID, session.RefreshTokenHash, session.DeviceName, session.IPAddress, session.ExpiredAt)
	var id uuid.UUID
	if err := row.Scan(&id); err != nil {
		return uuid.Nil, r.mapper.MapError(err, func(unhandledErr error) {
			log.Error().Err(unhandledErr).Msg("failed to execute create query")
		})
	}

	return id, nil
}

func (r sessionRepository) Delete(ctx context.Context, userID, sessionID uuid.UUID) error {
	log := logger.WithLayerCtx(ctx, logger.LayerPostgresRepository)

	query := `
		DELETE FROM sessions
		WHERE user_id = $1 AND id = $2`

	db := GetDB(ctx, r.db)
	if _, err := db.Exec(ctx, query, userID, sessionID); err != nil {
		return r.mapper.MapError(err, func(unhandledErr error) {
			log.Error().Err(unhandledErr).Msg("failed to execute delete query")
		})
	}

	return nil
}

func (r sessionRepository) GetByRefreshToken(ctx context.Context, refreshToken string) (Session, error) {
	log := logger.WithLayerCtx(ctx, logger.LayerPostgresRepository)

	query := `
		SELECT id, user_id, refresh_token_hash, device_name, ip_address, expired_at, created_at, updated_at
		FROM sessions
		WHERE refresh_token_hash = $1`

	db := GetDB(ctx, r.db)
	rows, err := db.Query(ctx, query, refreshToken)
	if err != nil {
		log.Error().Err(err).Msg("failed to execute get by refresh token query")
		return Session{}, err
	}
	defer rows.Close()

	item, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Session])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Session{}, nil
		}

		log.Error().Err(err).Msg("failed to execute collect one row")
		return Session{}, err
	}

	return item, nil
}

func (r sessionRepository) Update(ctx context.Context, sessionID uuid.UUID, refreshTokenHash string, expiredAt time.Time) error {
	log := logger.WithLayerCtx(ctx, logger.LayerPostgresRepository)

	query := `
		UPDATE sessions
		SET refresh_token_hash = $1, expired_at = $2
		WHERE id = $3`

	db := GetDB(ctx, r.db)
	if _, err := db.Exec(ctx, query, refreshTokenHash, expiredAt, sessionID); err != nil {
		return r.mapper.MapError(err, func(unhandledErr error) {
			log.Error().Err(unhandledErr).Msg("failed to execute update query")
		})
	}

	return nil
}
