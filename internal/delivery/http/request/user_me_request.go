package request

import "github.com/rizalarfiyan/be-plant-factory/internal/domain/entity"

type UserMeGettingStartedRequest struct {
	Name   string `json:"name" validate:"required,alphaspace,min=3,max=64" example:"Rizal Arfiyan"`
	Avatar string `json:"avatar" validate:"required,startswith=avatar:" example:"avatar:coding"`
}

func (r UserMeGettingStartedRequest) ToEntity() entity.UserMeGettingStarted {
	return entity.UserMeGettingStarted{
		Name:   r.Name,
		Avatar: r.Avatar,
	}
}
