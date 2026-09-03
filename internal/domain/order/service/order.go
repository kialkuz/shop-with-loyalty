package service

import (
	"context"
	balanceService "kialkuz/shop-with-loyalty/internal/domain/balance/service"
	"kialkuz/shop-with-loyalty/internal/domain/drawal/model"
	drawalService "kialkuz/shop-with-loyalty/internal/domain/drawal/service"
	orderModel "kialkuz/shop-with-loyalty/internal/domain/order/model"
	infrastructureInterfaces "kialkuz/shop-with-loyalty/internal/infrastructure/interfaces"
	pkgErrors "kialkuz/shop-with-loyalty/pkg/errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type OrderRepository interface {
	GetByNumber(ctx context.Context, number string) (*orderModel.Order, error)
	Add(ctx context.Context, order orderModel.Order) error
	UpdateTx(ctx context.Context, tx pgx.Tx, order orderModel.Order) error
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]orderModel.Order, error)
	GetOrdersForGetAccrual(ctx context.Context) (map[string]orderModel.Order, error)
}

type OrderService struct {
	transactionManager infrastructureInterfaces.TransactionManager
	repository         OrderRepository
	balanceService     *balanceService.BalanceService
	drawalService      *drawalService.DrawalService
}

func NewOrderService(
	transactionManager infrastructureInterfaces.TransactionManager,
	repository OrderRepository,
	balanceService *balanceService.BalanceService,
	drawalService *drawalService.DrawalService,
) *OrderService {
	return &OrderService{
		transactionManager: transactionManager,
		repository:         repository,
		balanceService:     balanceService,
		drawalService:      drawalService,
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

func (s *OrderService) WithDraw(ctx context.Context, userID uuid.UUID, drawal model.Drawal) error {
	return s.transactionManager.WithinTransaction(ctx, func(tx pgx.Tx) error {
		balance, err := s.balanceService.GetByUserIDWithBlockForUpdate(ctx, tx, userID)
		if err != nil {
			return err
		}

		if balance.Current < drawal.Sum {
			return pkgErrors.ErrLessDrawals
		}

		balance.Current -= drawal.Sum
		balance.WithDrawn += drawal.Sum
		err = s.balanceService.UpdateTx(ctx, tx, *balance)
		if err != nil {
			return err
		}

		err = s.drawalService.AddTx(ctx, tx, drawal)
		if err != nil {
			return err
		}

		return nil
	})
}
