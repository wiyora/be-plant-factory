package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rizalarfiyan/be-plant-factory/internal/domain/entity"
	"github.com/rizalarfiyan/be-plant-factory/internal/infrastructure/logger"
	"github.com/rizalarfiyan/be-plant-factory/internal/infrastructure/postgres"
	"github.com/rizalarfiyan/be-plant-factory/internal/shared/helper"
	"github.com/samber/do/v2"
)

type UserRepository interface {
	UpsertSocialUser(ctx context.Context, user entity.User) (UpsertSocialUser, error)
	GetById(ctx context.Context, id uuid.UUID) (User, error)
	UpdateGettingStarted(ctx context.Context, req entity.UserMeGettingStarted) error
}

type userRepository struct {
	db     *pgxpool.Pool
	mapper helper.PgxErrorMapper
}

func NewUserRepository(i do.Injector) (UserRepository, error) {
	db := do.MustInvoke[*postgres.Database](i)
	mapper := helper.NewPgxErrorMapper()

	return userRepository{
		db:     db.Pool(),
		mapper: mapper,
	}, nil
}

func (r userRepository) UpsertSocialUser(ctx context.Context, user entity.User) (UpsertSocialUser, error) {
	log := logger.WithLayerCtx(ctx, logger.LayerPostgresRepository)

	query := `
	INSERT INTO users (email, name, avatar, current_step)
	VALUES ($1, $2, $3, $4)
	ON CONFLICT (email) DO UPDATE
	SET last_logged_in_at = now()
	RETURNING id, current_step`

	db := GetDB(ctx, r.db)
	row := db.QueryRow(ctx, query, user.Email, user.Name, user.Avatar, user.CurrentStep)
	var result UpsertSocialUser
	if err := row.Scan(&result.ID, &result.CurrentStep); err != nil {
		log.Error().Err(err).Msg("failed to execute query")
		return UpsertSocialUser{}, err
	}

	return result, nil
}

func (r userRepository) GetById(ctx context.Context, id uuid.UUID) (User, error) {
	log := logger.WithLayerCtx(ctx, logger.LayerPostgresRepository)

	query := `
	SELECT id, email, name, avatar, is_super_admin, current_step, last_logged_in_at, created_at, updated_at
	FROM users
	WHERE id = $1`

	db := GetDB(ctx, r.db)
	rows, err := db.Query(ctx, query, id)
	if err != nil {
		log.Error().Err(err).Msg("failed to execute get by id query")
		return User{}, err
	}
	defer rows.Close()

	item, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[User])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return User{}, nil
		}

		log.Error().Err(err).Msg("failed to execute collect one row")
		return User{}, err
	}

	return item, nil
}

func (r userRepository) UpdateGettingStarted(ctx context.Context, req entity.UserMeGettingStarted) error {
	log := logger.WithLayerCtx(ctx, logger.LayerPostgresRepository)

	query := `
	UPDATE users
	SET name = $1, avatar = $2, current_step = $3, updated_at = now()
	WHERE id = $4`

	db := GetDB(ctx, r.db)
	_, err := db.Exec(ctx, query, req.Name, req.Avatar, entity.CurrentStepCompleted, req.UserID)
	if err != nil {
		return r.mapper.MapError(err, func(unhandledErr error) {
			log.Error().Err(unhandledErr).Msg("failed to update getting started")
		})
	}

	return nil
}
