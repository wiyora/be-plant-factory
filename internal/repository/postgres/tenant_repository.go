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

type TenantRepository interface {
	List(ctx context.Context, req entity.TenantFilter) ([]ListTenant, uint64, error)
	GetById(ctx context.Context, id uuid.UUID) (Tenant, error)
	Create(ctx context.Context, tenant entity.Tenant) error
	Update(ctx context.Context, tenant entity.Tenant) (bool, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status entity.TenantStatus) (bool, error)
}

type tenantRepository struct {
	db     *pgxpool.Pool
	mapper helper.PgxErrorMapper
}

func NewTenantRepository(i do.Injector) (TenantRepository, error) {
	db := do.MustInvoke[*postgres.Database](i)
	mapper := helper.NewPgxErrorMapper()
	mapper.Register("tenants_name_key", domainError.NewManualValidation("name", "DB_EXISTS"))

	return tenantRepository{
		db:     db.Pool(),
		mapper: mapper,
	}, nil
}

func (r tenantRepository) List(ctx context.Context, req entity.TenantFilter) ([]ListTenant, uint64, error) {
	log := logger.WithLayerCtx(ctx, logger.LayerPostgresRepository)

	countQuery := pgq.Select("COUNT(*)").From("tenants")
	countQuery = r.listFilters(countQuery, req)

	countSQL, countArgs, err := countQuery.SQL()
	if err != nil {
		log.Error().Err(err).Msg("failed to build tenant count query")
		return nil, 0, err
	}

	var total uint64
	db := GetDB(ctx, r.db)
	if err := db.QueryRow(ctx, countSQL, countArgs...).Scan(&total); err != nil {
		log.Error().Err(err).Msg("failed to execute tenant count query")
		return nil, 0, err
	}

	query := pgq.Select("id", "name", "logo", "status", "created_at").From("tenants")
	query = r.listFilters(query, req)
	query = query.OrderBy(req.Order.String()).Limit(req.Pagination.PageSize).Offset(req.Pagination.Offset())

	sql, args, err := query.SQL()
	if err != nil {
		log.Error().Err(err).Msg("failed to build tenant list query")
		return nil, 0, err
	}

	rows, err := db.Query(ctx, sql, args...)
	if err != nil {
		log.Error().Err(err).Msg("failed to execute tenant list query")
		return nil, 0, err
	}
	defer rows.Close()

	items, err := pgx.CollectRows(rows, pgx.RowToStructByName[ListTenant])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, total, nil
		}

		log.Error().Err(err).Msg("failed to scan tenant list rows")
		return nil, 0, err
	}

	return items, total, nil
}

func (r tenantRepository) listFilters(query pgq.SelectBuilder, req entity.TenantFilter) pgq.SelectBuilder {
	if req.Status.Valid() {
		query = query.Where("status = ?", req.Status)
	}

	if req.Search.HasSearch() {
		query = query.Where("name % ?", req.Search)
	}

	return query
}

func (r tenantRepository) GetById(ctx context.Context, id uuid.UUID) (Tenant, error) {
	log := logger.WithLayerCtx(ctx, logger.LayerPostgresRepository)

	query := `
	SELECT id, name, logo, status, created_at, updated_at
	FROM tenants
	WHERE id = $1`

	db := GetDB(ctx, r.db)
	rows, err := db.Query(ctx, query, id)
	if err != nil {
		log.Error().Err(err).Msg("failed to execute get by id query")
		return Tenant{}, err
	}
	defer rows.Close()

	item, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Tenant])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Tenant{}, nil
		}

		log.Error().Err(err).Msg("failed to execute collect one row")
		return Tenant{}, err
	}

	return item, nil
}

func (r tenantRepository) Create(ctx context.Context, tenant entity.Tenant) error {
	log := logger.WithLayerCtx(ctx, logger.LayerPostgresRepository)

	query := `
	INSERT INTO tenants (name, logo)
	VALUES ($1, $2)`

	db := GetDB(ctx, r.db)
	_, err := db.Exec(ctx, query, tenant.Name, tenant.Logo)
	if err != nil {
		return r.mapper.MapError(err, func(unhandledErr error) {
			log.Error().Err(unhandledErr).Msg("failed to execute tenant create query")
		})
	}

	return nil
}

func (r tenantRepository) Update(ctx context.Context, tenant entity.Tenant) (bool, error) {
	log := logger.WithLayerCtx(ctx, logger.LayerPostgresRepository)

	query := `
	UPDATE tenants
	SET name = $1, logo = $2, updated_at = NOW()
	WHERE id = $3`

	db := GetDB(ctx, r.db)
	conn, err := db.Exec(ctx, query, tenant.Name, tenant.Logo, tenant.ID)
	if err != nil {
		return false, r.mapper.MapError(err, func(unhandledErr error) {
			log.Error().Err(unhandledErr).Msg("failed to update tenant")
		})
	}

	return conn.RowsAffected() > 0, nil
}

func (r tenantRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status entity.TenantStatus) (bool, error) {
	log := logger.WithLayerCtx(ctx, logger.LayerPostgresRepository)

	query := `
	UPDATE tenants
	SET status = $1, updated_at = NOW()
	WHERE id = $2 AND status != $1`

	db := GetDB(ctx, r.db)
	conn, err := db.Exec(ctx, query, status, id)
	if err != nil {
		return false, r.mapper.MapError(err, func(unhandledErr error) {
			log.Error().Err(unhandledErr).Msg("failed to update tenant status")
		})
	}

	return conn.RowsAffected() > 0, nil
}
