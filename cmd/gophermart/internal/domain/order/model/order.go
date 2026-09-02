package model

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID         uuid.UUID
	UserId     uuid.UUID
	Number     Number
	Status     *Status
	Accrual    *int
	UploadedAt time.Time
}

func (m *Order) IsLoadedByUser(userId uuid.UUID) bool {
	return m.UserId == userId
}
