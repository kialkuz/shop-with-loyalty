package db

import (
	"kialkuz/gophermart/internal/infrastructure/storage"
)

type ActiveTokensRepository struct {
	db *storage.DB
}

func NewActiveTokensRepository(db *storage.DB) *ActiveTokensRepository {
	return &ActiveTokensRepository{db: db}
}
