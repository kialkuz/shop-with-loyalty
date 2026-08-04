package db

import (
	"kialkuz/gophermart/internal/infrastructure/storage"
)

type OrdersRepository struct {
	db *storage.DB
}

func NewOrdersRepository(db *storage.DB) *OrdersRepository {
	return &OrdersRepository{db: db}
}
