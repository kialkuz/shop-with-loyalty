package db

import (
	"kialkuz/shop-with-loyalty/internal/infrastructure/storage"
)

type OrderRepository struct {
	db *storage.DB
}

func NewOrderRepository(db *storage.DB) *OrderRepository {
	return &OrderRepository{db: db}
}
