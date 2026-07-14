package usertenant

import (
	"context"

	"github.com/google/uuid"
	domainError "github.com/rizalarfiyan/be-plant-factory/internal/domain/error"
	"github.com/rizalarfiyan/be-plant-factory/internal/shared/code"
)

func (u userTenantUseCase) AssignRole(ctx context.Context, tenantID, userID, roleID uuid.UUID) error {
	updated, err := u.userTenantRepo.Upsert(ctx, tenantID, userID, roleID)
	if err != nil {
		return err
	}

	if !updated {
		return domainError.New(code.UserTenantExists)
	}

	return nil
}
