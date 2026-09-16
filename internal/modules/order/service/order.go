package service

import (
	"context"
	orderModel "kialkuz/shop-with-loyalty/internal/modules/order/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

//go:generate go run go.uber.org/mock/mockgen -source=order.go -destination=mocks/order_mock.go -package=mocks -typed
type OrderRepository interface {
	GetByNumber(ctx context.Context, number string) (*orderModel.Order, error)
	Add(ctx context.Context, order orderModel.Order) error
	UpdateTx(ctx context.Context, tx pgx.Tx, order orderModel.Order) error
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]orderModel.Order, error)
	GetOrdersForGetAccrual(ctx context.Context) (map[string]orderModel.Order, error)
}

type OrderService struct {
	repository OrderRepository
}

func NewOrderService(
	repository OrderRepository,
) *OrderService {
	return &OrderService{
		repository: repository,
	}
}

func (s *OrderService) GetByNumber(ctx context.Context, number string) (*orderModel.Order, error) {
	return s.repository.GetByNumber(ctx, number)
}

func (s *OrderService) Add(ctx context.Context, order orderModel.Order) error {
	return s.repository.Add(ctx, order)
}

func (s *OrderService) UpdateTx(ctx context.Context, tx pgx.Tx, order orderModel.Order) error {
	return s.repository.UpdateTx(ctx, tx, order)
}

func (s *OrderService) GetByUserID(ctx context.Context, userID uuid.UUID) ([]orderModel.Order, error) {
	return s.repository.GetByUserID(ctx, userID)
}

func (s *OrderService) GetOrdersForGetAccrual(ctx context.Context) (map[string]orderModel.Order, error) {
	return s.repository.GetOrdersForGetAccrual(ctx)
}
