package usecase

import (
	"context"
	"kialkuz/shop-with-loyalty/internal/infrastructure"
	balanceServ "kialkuz/shop-with-loyalty/internal/modules/balance/service"
	"kialkuz/shop-with-loyalty/internal/modules/drawal/model"
	drawalServ "kialkuz/shop-with-loyalty/internal/modules/drawal/service"
	pkgErrors "kialkuz/shop-with-loyalty/pkg/errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type WithDraw struct {
	transactionManager infrastructure.Transaction
	drawalService      *drawalServ.DrawalService
	balanceService     *balanceServ.BalanceService
}

func NewWithDrawUseCase(
	transactionManager infrastructure.Transaction,
	drawalService *drawalServ.DrawalService,
	balanceService *balanceServ.BalanceService,
) *WithDraw {
	return &WithDraw{
		transactionManager: transactionManager,
		drawalService:      drawalService,
		balanceService:     balanceService,
	}
}

func (s *WithDraw) Execute(ctx context.Context, userID uuid.UUID, drawal model.Drawal) error {
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
		if err := s.balanceService.UpdateTx(ctx, tx, *balance); err != nil {
			return err
		}

		if err := s.drawalService.AddTx(ctx, tx, drawal); err != nil {
			return err
		}

		return nil
	})
}
