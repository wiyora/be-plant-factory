package request

import "github.com/rizalarfiyan/be-plant-factory/internal/domain/entity"

type CreateUserRequest struct {
	Email  string `json:"email" validate:"required,email,max=255" example:"rizalarfiyan@plant-factory.com"`
	Name   string `json:"name" validate:"required,alphaspace,min=3,max=64" example:"Rizal Arfiyan"`
	Avatar string `json:"avatar" validate:"required,startswith=avatar:" example:"avatar: coding"`
}

func (r CreateUserRequest) ToEntity() entity.User {
	return entity.User{
		Email:  r.Email,
		Name:   r.Name,
		Avatar: r.Avatar,
	}
}

type UpdateUserRequest struct {
	Name   string `json:"name" validate:"required,alphaspace,min=3,max=64" example:"Rizal Arfiyan"`
	Avatar string `json:"avatar" validate:"required,startswith=avatar:" example:"avatar: coding"`
}

func (r UpdateUserRequest) ToEntity() entity.User {
	return entity.User{
		Name:   r.Name,
		Avatar: r.Avatar,
	}
}
