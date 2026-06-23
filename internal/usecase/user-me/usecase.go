package userMe

import (
	"context"

	"github.com/rizalarfiyan/be-plant-factory/internal/domain/entity"
	"github.com/rizalarfiyan/be-plant-factory/internal/repository/cache"
	"github.com/rizalarfiyan/be-plant-factory/internal/repository/postgres"
	"github.com/samber/do/v2"
)

type UserMeUseCase interface {
	GettingStarted(ctx context.Context, req entity.UserMeGettingStarted) error
}

type userMeUseCase struct {
	userRepo  postgres.UserRepository `do:""`
	userCache cache.UserRepository    `do:""`
}

func NewUserMeUseCase(i do.Injector) (UserMeUseCase, error) {
	return do.InvokeStruct[userMeUseCase](i)
}
