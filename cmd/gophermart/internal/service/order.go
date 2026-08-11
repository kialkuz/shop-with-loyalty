package service

import "kialkuz/shop-with-loyalty/internal/contracts/repository"

type OrderService struct {
	repository repository.OrderRepository
}

func NewOrderService(repository repository.OrderRepository) *OrderService {
	return &OrderService{repository: repository}
}
