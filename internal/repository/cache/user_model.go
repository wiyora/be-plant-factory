package cache

import (
	"github.com/google/uuid"
	"github.com/rizalarfiyan/be-plant-factory/internal/domain/entity"
)

type AuthContext struct {
	ID          uuid.UUID          `json:"id"`
	Email       string             `json:"email"`
	Name        string             `json:"name"`
	Avatar      string             `json:"avatar"`
	CurrentStep entity.CurrentStep `json:"current_step"`
}

func (a AuthContext) ToEntity() entity.AuthContext {
	return entity.AuthContext{
		ID:          a.ID,
		Email:       a.Email,
		Name:        a.Name,
		Avatar:      a.Avatar,
		CurrentStep: a.CurrentStep,
	}
}
