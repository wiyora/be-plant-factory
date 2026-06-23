package entity

import "github.com/google/uuid"

type UserMeGettingStarted struct {
	UserID      uuid.UUID
	Name        string
	Avatar      string
	CurrentStep CurrentStep
}
