package dto

import "github.com/google/uuid"

type Token struct {
	Jti    uuid.UUID
	UserID uuid.UUID
}
