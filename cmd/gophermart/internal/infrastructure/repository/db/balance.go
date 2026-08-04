package db

import (
	"kialkuz/gophermart/internal/infrastructure/storage"
)

type BalanceRepository struct {
	db *storage.DB
}

func NewBalanceRepository(db *storage.DB) *BalanceRepository {
	return &BalanceRepository{db: db}
}
