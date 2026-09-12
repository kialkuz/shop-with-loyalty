package usecase

import (
	"context"
	"kialkuz/shop-with-loyalty/internal/infrastructure"
	balanceServ "kialkuz/shop-with-loyalty/internal/modules/balance/service"
	orderModel "kialkuz/shop-with-loyalty/internal/modules/order/model"
	orderServ "kialkuz/shop-with-loyalty/internal/modules/order/service"

	"github.com/jackc/pgx/v5"
)

type AddAccrual struct {
	transactionManager infrastructure.Transaction
	orderService       *orderServ.OrderService
	balanceService     *balanceServ.BalanceService
}

func NewAddAccrualUseCase(
	transactionManager infrastructure.Transaction,
	orderService *orderServ.OrderService,
	balanceService *balanceServ.BalanceService,
) *AddAccrual {
	return &AddAccrual{
		transactionManager: transactionManager,
		orderService:       orderService,
		balanceService:     balanceService,
	}
}

func (s *AddAccrual) Execute(
	ctx context.Context,
	orderForUpdate orderModel.Order,
) error {
	return s.transactionManager.RunTransaction(ctx, func(tx pgx.Tx) error {
		if err := s.orderService.UpdateTx(ctx, tx, orderForUpdate); err != nil {
			return err
		}

		userBalance, err := s.balanceService.GetByUserID(ctx, orderForUpdate.UserID)
		if err != nil {
			return err
		}

		userBalance.Current += *orderForUpdate.Accrual

		if err := s.balanceService.UpdateTx(ctx, tx, *userBalance); err != nil {
			return err
		}

		return nil
	})
}
