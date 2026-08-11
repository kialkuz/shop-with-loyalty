package model

import (
	"time"

	"github.com/google/uuid"
)

type Token struct {
	Jti       string
	UserId    uuid.UUID
	IssuedAt  time.Time
	ExpiresAt time.Time
}
