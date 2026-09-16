package model

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	Number     Number
	Status     *Status
	Accrual    *int
	UploadedAt time.Time
}

func (m *Order) IsLoadedByUser(userID uuid.UUID) bool {
	return m.UserID == userID
}
