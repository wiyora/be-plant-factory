package social

import (
	"context"

	"github.com/rizalarfiyan/be-plant-factory/internal/domain/entity"
)

type SocialAuthProvider interface {
	GetLoginURL(state string) string
	GetUserInfo(ctx context.Context, code string) (entity.SocialUser, error)
}
