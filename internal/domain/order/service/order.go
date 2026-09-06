package service

import (
	"context"
	balanceService "kialkuz/shop-with-loyalty/internal/domain/balance/service"
	"kialkuz/shop-with-loyalty/internal/domain/drawal/model"
	drawalService "kialkuz/shop-with-loyalty/internal/domain/drawal/service"
	orderModel "kialkuz/shop-with-loyalty/internal/domain/order/model"
	"kialkuz/shop-with-loyalty/internal/infrastructure"
	pkgErrors "kialkuz/shop-with-loyalty/pkg/errors"

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
	transactionManager infrastructure.Transaction
	repository         OrderRepository
	balanceService     *balanceService.BalanceService
	drawalService      *drawalService.DrawalService
}

func NewOrderService(
	transactionManager infrastructure.Transaction,
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
	return s.transactionManager.RunTransaction(ctx, func(tx pgx.Tx) error {
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

func (s *OrderService) UpdateAccrual(
	ctx context.Context,
	userID uuid.UUID,
	orderForUpdate orderModel.Order,
) error {
	return s.transactionManager.RunTransaction(ctx, func(tx pgx.Tx) error {
		err := s.UpdateTx(ctx, tx, orderForUpdate)
		if err != nil {
			return err
		}

		userBalance, err := s.balanceService.GetByUserID(ctx, userID)
		if err != nil {
			return err
		}

		userBalance.Current += *orderForUpdate.Accrual

		err = s.balanceService.UpdateTx(ctx, tx, *userBalance)
		if err != nil {
			return err
		}

		return nil
	})
}
