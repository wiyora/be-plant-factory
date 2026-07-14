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

type RoleRepository interface {
	List(ctx context.Context, req entity.RoleFilter) ([]ListRole, uint64, error)
	GetById(ctx context.Context, id uuid.UUID) (Role, error)
	Create(ctx context.Context, role entity.Role) error
	Update(ctx context.Context, role entity.Role) (bool, error)
	Delete(ctx context.Context, id uuid.UUID) (bool, error)
	Dropdown(ctx context.Context, req entity.DropdownFilter) ([]DropdownRole, uint64, error)
}

type roleRepository struct {
	db     *pgxpool.Pool
	mapper helper.PgxErrorMapper
}

func NewRoleRepository(i do.Injector) (RoleRepository, error) {
	db := do.MustInvoke[*postgres.Database](i)
	mapper := helper.NewPgxErrorMapper()
	mapper.Register("roles_name_key", domainError.NewManualValidation("name", "DB_EXISTS"))

	return roleRepository{
		db:     db.Pool(),
		mapper: mapper,
	}, nil
}

func (r roleRepository) List(ctx context.Context, req entity.RoleFilter) ([]ListRole, uint64, error) {
	log := logger.WithLayerCtx(ctx, logger.LayerPostgresRepository)

	countQuery := pgq.Select("COUNT(*)").From("roles")
	countQuery = r.listFilters(countQuery, req)

	countSQL, countArgs, err := countQuery.SQL()
	if err != nil {
		log.Error().Err(err).Msg("failed to build role count query")
		return nil, 0, err
	}

	var total uint64
	db := GetDB(ctx, r.db)
	if err := db.QueryRow(ctx, countSQL, countArgs...).Scan(&total); err != nil {
		log.Error().Err(err).Msg("failed to execute role count query")
		return nil, 0, err
	}

	subQuery := pgq.Select("id", "name", "total_permission").From("roles")
	subQuery = r.listFilters(subQuery, req)
	subQuery = subQuery.OrderBy(req.Order.String()).Limit(req.Pagination.PageSize).Offset(req.Pagination.Offset())

	query := pgq.Select("pr.id", "pr.name", "pr.total_permission", "sub.total_user").
		FromSelect(subQuery, "pr").
		LeftJoin("LATERAL (SELECT COUNT(*) AS total_user FROM user_tenants ut WHERE ut.role_id = pr.id) sub ON TRUE")

	sql, args, err := query.SQL()
	if err != nil {
		log.Error().Err(err).Msg("failed to build role list query")
		return nil, 0, err
	}

	rows, err := db.Query(ctx, sql, args...)
	if err != nil {
		log.Error().Err(err).Msg("failed to execute role list query")
		return nil, 0, err
	}
	defer rows.Close()

	items, err := pgx.CollectRows(rows, pgx.RowToStructByName[ListRole])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, total, nil
		}

		log.Error().Err(err).Msg("failed to scan role list rows")
		return nil, 0, err
	}

	return items, total, nil
}

func (r roleRepository) listFilters(query pgq.SelectBuilder, req entity.RoleFilter) pgq.SelectBuilder {
	if req.Search.HasSearch() {
		query = query.Where("name % ?", req.Search)
	}

	return query
}

func (r roleRepository) GetById(ctx context.Context, id uuid.UUID) (Role, error) {
	log := logger.WithLayerCtx(ctx, logger.LayerPostgresRepository)

	query := `
	SELECT id, name, total_permission, permissions, created_at, updated_at
	FROM roles
	WHERE id = $1`

	db := GetDB(ctx, r.db)
	rows, err := db.Query(ctx, query, id)
	if err != nil {
		log.Error().Err(err).Msg("failed to execute get role by id query")
		return Role{}, err
	}
	defer rows.Close()

	item, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Role])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Role{}, nil
		}

		log.Error().Err(err).Msg("failed to execute collect one row")
		return Role{}, err
	}

	return item, nil
}

func (r roleRepository) Create(ctx context.Context, role entity.Role) error {
	log := logger.WithLayerCtx(ctx, logger.LayerPostgresRepository)

	query := `
	INSERT INTO roles (name, total_permission, permissions)
	VALUES ($1, $2, $3)`

	db := GetDB(ctx, r.db)
	_, err := db.Exec(ctx, query, role.Name, role.TotalPermission, role.Permissions)
	if err != nil {
		return r.mapper.MapError(err, func(unhandledErr error) {
			log.Error().Err(unhandledErr).Msg("failed to execute role create query")
		})
	}

	return nil
}

func (r roleRepository) Update(ctx context.Context, role entity.Role) (bool, error) {
	log := logger.WithLayerCtx(ctx, logger.LayerPostgresRepository)

	query := `
	UPDATE roles
	SET name = $1, total_permission = $2, permissions = $3, updated_at = NOW()
	WHERE id = $4`

	db := GetDB(ctx, r.db)
	conn, err := db.Exec(ctx, query, role.Name, role.TotalPermission, role.Permissions, role.ID)
	if err != nil {
		return false, r.mapper.MapError(err, func(unhandledErr error) {
			log.Error().Err(unhandledErr).Msg("failed to update role")
		})
	}

	return conn.RowsAffected() > 0, nil
}

func (r roleRepository) Delete(ctx context.Context, id uuid.UUID) (bool, error) {
	log := logger.WithLayerCtx(ctx, logger.LayerPostgresRepository)

	query := `DELETE FROM roles WHERE id = $1`

	db := GetDB(ctx, r.db)
	conn, err := db.Exec(ctx, query, id)
	if err != nil {
		return false, r.mapper.MapError(err, func(unhandledErr error) {
			log.Error().Err(unhandledErr).Msg("failed to delete role")
		})
	}

	return conn.RowsAffected() > 0, nil
}

func (r roleRepository) dropdownFilters(query pgq.SelectBuilder, req entity.DropdownFilter) pgq.SelectBuilder {
	hasActiveIDs := len(req.ActiveIDs) > 0
	hasSearch := req.Search.HasSearch()

	if hasActiveIDs && hasSearch {
		query = query.Where("id = ANY(?) OR name % ?", req.ActiveIDs, req.Search)
	} else if hasSearch {
		query = query.Where("name % ?", req.Search)
	}

	return query
}

func (r roleRepository) dropdownOrderBy(query pgq.SelectBuilder, req entity.DropdownFilter) pgq.SelectBuilder {
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

func (r roleRepository) Dropdown(ctx context.Context, req entity.DropdownFilter) ([]DropdownRole, uint64, error) {
	log := logger.WithLayerCtx(ctx, logger.LayerPostgresRepository)

	countQuery := pgq.Select("COUNT(*)").From("roles")
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

	query := pgq.Select("id", "name").From("roles")
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

	items, err := pgx.CollectRows(rows, pgx.RowToStructByName[DropdownRole])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, total, nil
		}
		log.Error().Err(err).Msg("failed to scan role dropdown rows")
		return nil, 0, err
	}

	return items, total, nil
}
