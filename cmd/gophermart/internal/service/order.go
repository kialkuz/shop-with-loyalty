package service

import (
	"context"
	"kialkuz/shop-with-loyalty/internal/model/order"

	"github.com/google/uuid"
)

type OrderRepository interface {
	GetByNumber(ctx context.Context, number string) (*order.Order, error)
	Add(ctx context.Context, order order.Order) error
	GetByUserId(ctx context.Context, userId uuid.UUID) ([]order.Order, error)
}

type OrderService struct {
	repository OrderRepository
}

func NewOrderService(repository OrderRepository) *OrderService {
	return &OrderService{repository: repository}
}

func (s *OrderService) GetByNumber(ctx context.Context, number string) (*order.Order, error) {
	return s.repository.GetByNumber(ctx, number)
}

func (s *OrderService) Add(ctx context.Context, order order.Order) error {
	return s.repository.Add(ctx, order)
}

func (s *OrderService) GetByUserId(ctx context.Context, userId uuid.UUID) ([]order.Order, error) {
	return s.repository.GetByUserId(ctx, userId)
}
