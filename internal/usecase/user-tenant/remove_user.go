package usertenant

import (
	"context"

	"github.com/google/uuid"
	domainError "github.com/rizalarfiyan/be-plant-factory/internal/domain/error"
	"github.com/rizalarfiyan/be-plant-factory/internal/shared/code"
)

func (u userTenantUseCase) RemoveUser(ctx context.Context, tenantID, userID uuid.UUID) error {
	deleted, err := u.userTenantRepo.Delete(ctx, tenantID, userID)
	if err != nil {
		return err
	}

	if !deleted {
		return domainError.New(code.InvalidParamID)
	}

	return nil
}
