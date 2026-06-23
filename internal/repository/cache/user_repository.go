package cache

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/rizalarfiyan/be-plant-factory/internal/config"
	domainError "github.com/rizalarfiyan/be-plant-factory/internal/domain/error"
	"github.com/rizalarfiyan/be-plant-factory/internal/infrastructure/logger"
	"github.com/rizalarfiyan/be-plant-factory/internal/repository/postgres"
	"github.com/rizalarfiyan/be-plant-factory/internal/repository/redis"
	"github.com/rizalarfiyan/be-plant-factory/internal/shared/code"
	"github.com/rizalarfiyan/be-plant-factory/internal/shared/helper"
	"github.com/samber/do/v2"
)

type UserRepository interface {
	GetById(ctx context.Context, userId uuid.UUID) (AuthContext, error)
	Clear(ctx context.Context, userID uuid.UUID) error
}

type userRepository struct {
	conf     *config.Config
	cache    redis.CacheRepository
	userRepo postgres.UserRepository
}

func NewUserRepository(i do.Injector) (UserRepository, error) {
	conf := do.MustInvoke[*config.Config](i)
	cacheRepo := do.MustInvoke[redis.CacheRepository](i)
	userRepo := do.MustInvoke[postgres.UserRepository](i)

	return userRepository{
		conf:     conf,
		cache:    cacheRepo,
		userRepo: userRepo,
	}, nil
}

func (r userRepository) key(userID uuid.UUID) string {
	return fmt.Sprintf("auth:user:%s", userID.String())
}

func (r userRepository) GetById(ctx context.Context, userID uuid.UUID) (AuthContext, error) {
	log := logger.WithLayerCtx(ctx, logger.LayerCacheRepository)

	var authCtx AuthContext
	cacheKey := r.key(userID)

	if err := r.cache.Get(ctx, cacheKey, &authCtx); err == nil {
		return authCtx, nil
	}

	user, err := r.userRepo.GetById(ctx, userID)
	if err != nil {
		log.Error().Err(err).Msg("failed to fetch user from database")
		return AuthContext{}, err
	}

	if helper.IsEmptyStruct(user) {
		log.Warn().Msgf("user with id %s not found in database", userID.String())
		return AuthContext{}, domainError.New(code.UserNotFound)
	}

	authCtx = AuthContext{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		Avatar:      user.Avatar,
		CurrentStep: user.CurrentStep,
	}

	if err := r.cache.Set(ctx, cacheKey, authCtx, r.conf.Auth.CacheUserDuration); err != nil {
		log.Error().Err(err).Msg("failed to set user in cache")
	}

	return authCtx, nil
}

func (r userRepository) Clear(ctx context.Context, userID uuid.UUID) error {
	log := logger.WithLayerCtx(ctx, logger.LayerCacheRepository)

	if err := r.cache.Delete(ctx, r.key(userID)); err != nil {
		log.Error().Err(err).Msg("failed to delete user from cache")
		return err
	}

	return nil
}
