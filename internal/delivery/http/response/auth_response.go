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
	ID          uuid.UUID              `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Email       string                 `json:"email" example:"rizalarfiyan@plant-factory.com"`
	Name        string                 `json:"name" example:"Rizal Arfiyan"`
	Avatar      string                 `json:"avatar" example:"https://plant-factory.com/avatar.jpg"`
	CurrentStep entity.CurrentStep     `json:"current_step" example:"completed"`
	Tenants     []AuthMeTenantResponse `json:"tenants"`
}

type AuthMeTenantResponse struct {
	ID          uuid.UUID `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Logo        string    `json:"logo" example:"https://plant-factory.com/logo.jpg"`
	Name        string    `json:"name" example:"Plant Factory"`
	Role        string    `json:"role" example:"admin"`
	Permissions []string  `json:"permissions" example:"[\"read\", \"write\"]"`
}

func NewAuthMeResponse(user entity.AuthContext) AuthMeResponse {
	return AuthMeResponse{
		ID:          user.ID,
		Email:       user.Email,
		Name:        user.Name,
		Avatar:      user.Avatar,
		CurrentStep: user.CurrentStep,
		Tenants:     []AuthMeTenantResponse{},
	}
}
