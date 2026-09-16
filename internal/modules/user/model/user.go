package model

import (
	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID
	Login    string `validate:"required,gte=4,lte=50"`
	Password string `validate:"required,gte=5,lte=20"`
}
