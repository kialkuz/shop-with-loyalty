package model

import "github.com/google/uuid"

type Balance struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Current   int
	WithDrawn int
}
