package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/henvic/pgq"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rizalarfiyan/be-plant-factory/internal/domain/entity"
	domainError "github.com/rizalarfiyan/be-plant-factory/internal/domain/error"
	"github.com/rizalarfiyan/be-plant-factory/internal/infrastructure/logger"
	"github.com/rizalarfiyan/be-plant-factory/internal/infrastructure/postgres"
	"github.com/rizalarfiyan/be-plant-factory/internal/shared/helper"
	"github.com/samber/do/v2"
)

type UserRepository interface {
	UpsertSocialUser(ctx context.Context, user entity.User) (UpsertSocialUser, error)
	GetById(ctx context.Context, id uuid.UUID) (User, error)
	UpdateGettingStarted(ctx context.Context, req entity.UserMeGettingStarted) error
	List(ctx context.Context, req entity.UserFilter) ([]ListUser, uint64, error)
	Create(ctx context.Context, user entity.User) error
	Update(ctx context.Context, user entity.User) (bool, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status entity.UserStatus) (bool, error)
	Dropdown(ctx context.Context, req entity.DropdownFilter) ([]DropdownUser, uint64, error)
}

type userRepository struct {
	db     *pgxpool.Pool
	mapper helper.PgxErrorMapper
}

func NewUserRepository(i do.Injector) (UserRepository, error) {
	db := do.MustInvoke[*postgres.Database](i)
	mapper := helper.NewPgxErrorMapper()
	mapper.Register("users_email_key", domainError.NewManualValidation("email", "DB_EXISTS"))

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
	SET last_logged_in_at = NOW()
	RETURNING id, current_step, status`

	db := GetDB(ctx, r.db)
	rows, err := db.Query(ctx, query, user.Email, user.Name, user.Avatar, user.CurrentStep)
	if err != nil {
		log.Error().Err(err).Msg("failed to execute get by id query")
		return UpsertSocialUser{}, err
	}
	defer rows.Close()

	item, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[UpsertSocialUser])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return UpsertSocialUser{}, nil
		}

		log.Error().Err(err).Msg("failed to execute collect one row")
		return UpsertSocialUser{}, err
	}

	return item, nil
}

func (r userRepository) GetById(ctx context.Context, id uuid.UUID) (User, error) {
	log := logger.WithLayerCtx(ctx, logger.LayerPostgresRepository)

	query := `
	SELECT id, email, name, avatar, is_super_admin, current_step, status, last_logged_in_at, created_at, updated_at
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
	SET name = $1, avatar = $2, current_step = $3, updated_at = NOW()
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

func (r userRepository) List(ctx context.Context, req entity.UserFilter) ([]ListUser, uint64, error) {
	log := logger.WithLayerCtx(ctx, logger.LayerPostgresRepository)

	countQuery := pgq.Select("COUNT(*)").From("users")
	countQuery = r.listFilters(countQuery, req)

	countSQL, countArgs, err := countQuery.SQL()
	if err != nil {
		log.Error().Err(err).Msg("failed to build user count query")
		return nil, 0, err
	}

	var total uint64
	db := GetDB(ctx, r.db)
	if err := db.QueryRow(ctx, countSQL, countArgs...).Scan(&total); err != nil {
		log.Error().Err(err).Msg("failed to execute user count query")
		return nil, 0, err
	}

	query := pgq.Select("id", "email", "name", "avatar", "status", "created_at").From("users")
	query = r.listFilters(query, req)
	query = query.OrderBy(req.Order.String()).Limit(req.Pagination.PageSize).Offset(req.Pagination.Offset())

	sql, args, err := query.SQL()
	if err != nil {
		log.Error().Err(err).Msg("failed to build user list query")
		return nil, 0, err
	}

	rows, err := db.Query(ctx, sql, args...)
	if err != nil {
		log.Error().Err(err).Msg("failed to execute user list query")
		return nil, 0, err
	}
	defer rows.Close()

	items, err := pgx.CollectRows(rows, pgx.RowToStructByName[ListUser])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, total, nil
		}

		log.Error().Err(err).Msg("failed to scan user list rows")
		return nil, 0, err
	}

	return items, total, nil
}

func (r userRepository) listFilters(query pgq.SelectBuilder, req entity.UserFilter) pgq.SelectBuilder {
	if req.Status.Valid() {
		query = query.Where("status = ?", req.Status)
	}

	if req.Search.HasSearch() {
		query = query.Where("name % ?", req.Search)
	}

	return query
}

func (r userRepository) Create(ctx context.Context, user entity.User) error {
	log := logger.WithLayerCtx(ctx, logger.LayerPostgresRepository)

	query := `
	INSERT INTO users (email, name, avatar)
	VALUES ($1, $2, $3)`

	db := GetDB(ctx, r.db)
	_, err := db.Exec(ctx, query, user.Email, user.Name, user.Avatar)
	if err != nil {
		return r.mapper.MapError(err, func(unhandledErr error) {
			log.Error().Err(unhandledErr).Msg("failed to execute user create query")
		})
	}

	return nil
}

func (r userRepository) Update(ctx context.Context, user entity.User) (bool, error) {
	log := logger.WithLayerCtx(ctx, logger.LayerPostgresRepository)

	query := `
	UPDATE users
	SET name = $1, avatar = $2, updated_at = NOW()
	WHERE id = $3`

	db := GetDB(ctx, r.db)
	conn, err := db.Exec(ctx, query, user.Name, user.Avatar, user.ID)
	if err != nil {
		return false, r.mapper.MapError(err, func(unhandledErr error) {
			log.Error().Err(unhandledErr).Msg("failed to update user")
		})
	}

	return conn.RowsAffected() > 0, nil
}

func (r userRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status entity.UserStatus) (bool, error) {
	log := logger.WithLayerCtx(ctx, logger.LayerPostgresRepository)

	query := `
	UPDATE users
	SET status = $1, updated_at = NOW()
	WHERE id = $2 AND status != $1`

	db := GetDB(ctx, r.db)
	conn, err := db.Exec(ctx, query, status, id)
	if err != nil {
		return false, r.mapper.MapError(err, func(unhandledErr error) {
			log.Error().Err(unhandledErr).Msg("failed to update user status")
		})
	}

	return conn.RowsAffected() > 0, nil
}

func (r userRepository) dropdownFilters(query pgq.SelectBuilder, req entity.DropdownFilter) pgq.SelectBuilder {
	hasActiveIDs := len(req.ActiveIDs) > 0
	hasSearch := req.Search.HasSearch()

	if hasActiveIDs && hasSearch {
		query = query.Where("id = ANY(?) OR name % ?", req.ActiveIDs, req.Search)
	} else if hasSearch {
		query = query.Where("name % ?", req.Search)
	}

	return query
}

func (r userRepository) dropdownOrderBy(query pgq.SelectBuilder, req entity.DropdownFilter) pgq.SelectBuilder {
	hasActiveIDs := len(req.ActiveIDs) > 0
	hasSearch := req.Search.HasSearch()

	if hasActiveIDs && hasSearch {
		return query.OrderByClause("CASE WHEN id = ANY(?) THEN 0 ELSE 1 END, similarity(name, ?) DESC, name ASC", req.ActiveIDs, req.Search)
	}

	if hasActiveIDs {
		return query.OrderByClause("CASE WHEN id = ANY(?) THEN 0 ELSE 1 END, name ASC", req.ActiveIDs)
	}

	if hasSearch {
		return query.OrderByClause("similarity(name, ?) DESC, name ASC", req.Search)
	}

	return query.OrderBy("name ASC")
}

func (r userRepository) Dropdown(ctx context.Context, req entity.DropdownFilter) ([]DropdownUser, uint64, error) {
	log := logger.WithLayerCtx(ctx, logger.LayerPostgresRepository)

	countQuery := pgq.Select("COUNT(*)").From("users")
	countQuery = r.dropdownFilters(countQuery, req)

	countSQL, countArgs, err := countQuery.SQL()
	if err != nil {
		log.Error().Err(err).Msg("failed to build role dropdown count query")
		return nil, 0, err
	}

	var total uint64
	db := GetDB(ctx, r.db)
	if err := db.QueryRow(ctx, countSQL, countArgs...).Scan(&total); err != nil {
		log.Error().Err(err).Msg("failed to execute role dropdown count query")
		return nil, 0, err
	}

	query := pgq.Select("id", "name").From("users")
	query = r.dropdownFilters(query, req)
	query = r.dropdownOrderBy(query, req)
	query = query.Limit(req.Pagination.PageSize).Offset(req.Pagination.Offset())

	sql, args, err := query.SQL()
	if err != nil {
		log.Error().Err(err).Msg("failed to build role dropdown query")
		return nil, 0, err
	}

	rows, err := db.Query(ctx, sql, args...)
	if err != nil {
		log.Error().Err(err).Msg("failed to execute role dropdown query")
		return nil, 0, err
	}
	defer rows.Close()

	items, err := pgx.CollectRows(rows, pgx.RowToStructByName[DropdownUser])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, total, nil
		}
		log.Error().Err(err).Msg("failed to scan role dropdown rows")
		return nil, 0, err
	}

	return items, total, nil
}
