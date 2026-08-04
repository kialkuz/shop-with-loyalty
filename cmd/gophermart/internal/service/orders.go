package service

import "kialkuz/gophermart/internal/infrastructure/repository/db"

type OrdersService struct {
	repository *db.OrdersRepository
}

func NewOrdersService(repository *db.OrdersRepository) *OrdersService {
	return &OrdersService{repository: repository}
}
