package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/henvic/pgq"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rizalarfiyan/be-plant-factory/internal/domain/entity"
	"github.com/rizalarfiyan/be-plant-factory/internal/infrastructure/logger"
	"github.com/rizalarfiyan/be-plant-factory/internal/infrastructure/postgres"
	"github.com/rizalarfiyan/be-plant-factory/internal/shared/helper"
	"github.com/samber/do/v2"
)

type UserTenantRepository interface {
	ListByTenant(ctx context.Context, req entity.TenantUserFilter) ([]ListTenantUser, uint64, error)
	ListByUser(ctx context.Context, req entity.UserTenantFilter) ([]ListUserTenant, uint64, error)
	Upsert(ctx context.Context, tenantID, userID, roleID uuid.UUID) (bool, error)
	Delete(ctx context.Context, tenantID, userID uuid.UUID) (bool, error)
}

type userTenantRepository struct {
	db     *pgxpool.Pool
	mapper helper.PgxErrorMapper
}

func NewUserTenantRepository(i do.Injector) (UserTenantRepository, error) {
	db := do.MustInvoke[*postgres.Database](i)
	mapper := helper.NewPgxErrorMapper()

	return userTenantRepository{
		db:     db.Pool(),
		mapper: mapper,
	}, nil
}

func (r userTenantRepository) listByTenantFilters(query pgq.SelectBuilder, req entity.TenantUserFilter) pgq.SelectBuilder {
	query = query.Where("ut.tenant_id = ?", req.TenantID)
	if req.Search.HasSearch() {
		query = query.Where("u.name % ?", req.Search)
	}
	return query
}

func (r userTenantRepository) ListByTenant(ctx context.Context, req entity.TenantUserFilter) ([]ListTenantUser, uint64, error) {
	log := logger.WithLayerCtx(ctx, logger.LayerPostgresRepository)

	countQuery := pgq.Select("COUNT(*)").From("user_tenants ut")
	if req.Search.HasSearch() {
		countQuery = countQuery.Join("users u ON ut.user_id = u.id")
	}
	countQuery = r.listByTenantFilters(countQuery, req)

	countSQL, countArgs, err := countQuery.SQL()
	if err != nil {
		log.Error().Err(err).Msg("failed to build count query")
		return nil, 0, err
	}

	var total uint64
	db := GetDB(ctx, r.db)
	if err := db.QueryRow(ctx, countSQL, countArgs...).Scan(&total); err != nil {
		log.Error().Err(err).Msg("failed execute count query")
		return nil, 0, err
	}

	filteredQuery := pgq.Select("ut.user_id", "ut.role_id", "ut.created_at").
		From("user_tenants ut")
	if req.Search.HasSearch() {
		filteredQuery = filteredQuery.Join("users u ON ut.user_id = u.id")
	}
	filteredQuery = r.listByTenantFilters(filteredQuery, req)
	filteredQuery = filteredQuery.OrderBy(req.Order.String()).
		Limit(req.Pagination.PageSize).
		Offset(req.Pagination.Offset())

	mainQuery := pgq.Select(
		"u.id as user_id",
		"u.name as user_name",
		"u.email as user_email",
		"u.avatar as user_avatar",
		"u.status as user_status",
		"r.id as role_id",
		"r.name as role_name",
		"fut.created_at as assigned_date",
	).With("filtered_user_tenants", filteredQuery).
		From("filtered_user_tenants fut").
		Join("users u ON fut.user_id = u.id").
		Join("roles r ON fut.role_id = r.id").
		OrderBy(req.Order.String())

	sql, args, err := mainQuery.SQL()
	if err != nil {
		log.Error().Err(err).Msg("failed to build tenant users query")
		return nil, 0, err
	}

	rows, err := db.Query(ctx, sql, args...)
	if err != nil {
		log.Error().Err(err).Msg("failed execute tenant users query")
		return nil, 0, err
	}
	defer rows.Close()

	items, err := pgx.CollectRows(rows, pgx.RowToStructByName[ListTenantUser])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, total, nil
		}
		log.Error().Err(err).Msg("failed scan tenant users rows")
		return nil, 0, err
	}

	return items, total, nil
}

func (r userTenantRepository) listByUserFilters(query pgq.SelectBuilder, req entity.UserTenantFilter) pgq.SelectBuilder {
	query = query.Where("ut.user_id = ?", req.UserID)
	if len(req.TenantIDs) > 0 {
		query = query.Where("ut.tenant_id = ANY(?)", req.TenantIDs)
	}
	if len(req.RoleIDs) > 0 {
		query = query.Where("ut.role_id = ANY(?)", req.RoleIDs)
	}
	return query
}

func (r userTenantRepository) ListByUser(ctx context.Context, req entity.UserTenantFilter) ([]ListUserTenant, uint64, error) {
	log := logger.WithLayerCtx(ctx, logger.LayerPostgresRepository)

	countQuery := pgq.Select("COUNT(*)").From("user_tenants ut")
	countQuery = r.listByUserFilters(countQuery, req)

	countSQL, countArgs, err := countQuery.SQL()
	if err != nil {
		log.Error().Err(err).Msg("failed to build count query")
		return nil, 0, err
	}

	var total uint64
	db := GetDB(ctx, r.db)
	if err := db.QueryRow(ctx, countSQL, countArgs...).Scan(&total); err != nil {
		log.Error().Err(err).Msg("failed execute count query")
		return nil, 0, err
	}

	filteredQuery := pgq.Select("ut.user_id", "ut.tenant_id", "ut.role_id", "ut.created_at").
		From("user_tenants ut")
	filteredQuery = r.listByUserFilters(filteredQuery, req)
	filteredQuery = filteredQuery.OrderBy(req.Order.String()).
		Limit(req.Pagination.PageSize).
		Offset(req.Pagination.Offset())

	mainQuery := pgq.Select(
		"fut.user_id as user_tenant_id",
		"t.id as tenant_id",
		"t.name as tenant_name",
		"t.logo as tenant_logo",
		"r.id as role_id",
		"r.name as role_name",
		"fut.created_at as assigned_date",
	).With("filtered_user_tenants", filteredQuery).
		From("filtered_user_tenants fut").
		Join("tenants t ON fut.tenant_id = t.id").
		Join("roles r ON fut.role_id = r.id").
		OrderBy(req.Order.String())

	sql, args, err := mainQuery.SQL()
	if err != nil {
		log.Error().Err(err).Msg("failed to build user tenants query")
		return nil, 0, err
	}

	rows, err := db.Query(ctx, sql, args...)
	if err != nil {
		log.Error().Err(err).Msg("failed execute user tenants query")
		return nil, 0, err
	}
	defer rows.Close()

	items, err := pgx.CollectRows(rows, pgx.RowToStructByName[ListUserTenant])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, total, nil
		}
		log.Error().Err(err).Msg("failed scan user tenants rows")
		return nil, 0, err
	}

	return items, total, nil
}

func (r userTenantRepository) Upsert(ctx context.Context, tenantID, userID, roleID uuid.UUID) (bool, error) {
	log := logger.WithLayerCtx(ctx, logger.LayerPostgresRepository)

	query := `
		INSERT INTO user_tenants (user_id, tenant_id, role_id, created_at)
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT (user_id, tenant_id)
		DO UPDATE SET role_id = $3, updated_at = NOW()`

	db := GetDB(ctx, r.db)
	conn, err := db.Exec(ctx, query, userID, tenantID, roleID)
	if err != nil {
		return false, r.mapper.MapError(err, func(unhandledErr error) {
			log.Error().Err(unhandledErr).Msg("failed to upsert user tenant")
		})
	}

	return conn.RowsAffected() > 0, nil
}

func (r userTenantRepository) Delete(ctx context.Context, tenantID, userID uuid.UUID) (bool, error) {
	log := logger.WithLayerCtx(ctx, logger.LayerPostgresRepository)

	query := `DELETE FROM user_tenants WHERE tenant_id = $1 AND user_id = $2`

	db := GetDB(ctx, r.db)
	conn, err := db.Exec(ctx, query, tenantID, userID)
	if err != nil {
		log.Error().Err(err).Msg("failed to delete user tenant")
		return false, err
	}

	return conn.RowsAffected() > 0, nil
}
