package jwt

import (
	"time"

	"github.com/google/uuid"
	"github.com/rizalarfiyan/be-plant-factory/internal/config"
	"github.com/samber/do/v2"
)

type Service interface {
	GenerateAccessToken(userID, sessionID uuid.UUID, issuedAt time.Time) (string, error)
	ParseAccessToken(tokenString string) (AccessClaims, error)
}

type service struct {
	conf *config.Config `do:""`
}

func New(i do.Injector) (Service, error) {
	return do.InvokeStruct[service](i)
}
