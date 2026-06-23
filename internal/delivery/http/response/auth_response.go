package response

import (
	"github.com/google/uuid"
	"github.com/rizalarfiyan/be-plant-factory/internal/domain/entity"
)

type AuthGoogleLoginResponse struct {
	RedirectUrl string `json:"redirect_url" example:"https://accounts.google.com/o/oauth2/auth/..."`
}

func NewAuthGoogleLoginResponse(redirectUrl string) AuthGoogleLoginResponse {
	return AuthGoogleLoginResponse{
		RedirectUrl: redirectUrl,
	}
}

type AuthMeResponse struct {
	ID          uuid.UUID          `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Email       string             `json:"email" example:"rizalarfiyan@ominotes.com"`
	Name        string             `json:"name" example:"Rizal Arfiyan"`
	Avatar      string             `json:"avatar" example:"https://ominotes.com/avatar.jpg"`
	CurrentStep entity.CurrentStep `json:"current_step" example:"completed"`
}

func NewAuthMeResponse(user entity.AuthContext) AuthMeResponse {
	return AuthMeResponse{
		ID:          user.ID,
		Email:       user.Email,
		Name:        user.Name,
		Avatar:      user.Avatar,
		CurrentStep: user.CurrentStep,
	}
}
