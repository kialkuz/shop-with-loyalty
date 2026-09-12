package scheduler

import (
	"kialkuz/shop-with-loyalty/internal/config"
	"kialkuz/shop-with-loyalty/internal/infrastructure"
	accrualServ "kialkuz/shop-with-loyalty/internal/modules/accrual/service"
	accrualUsecase "kialkuz/shop-with-loyalty/internal/modules/accrual/usecase"
	balanceServ "kialkuz/shop-with-loyalty/internal/modules/balance/service"
	orderServ "kialkuz/shop-with-loyalty/internal/modules/order/service"
	"kialkuz/shop-with-loyalty/internal/scheduler/accrual"

	"go.uber.org/zap"
)

func New(
	config *config.Config,
	sugar *zap.SugaredLogger,
	transactionManager infrastructure.Transaction,
	orderService *orderServ.OrderService,
	balanceService *balanceServ.BalanceService,
	accrualService *accrualServ.AccrualService,
) *Runner {
	usecaseAddAccrual := accrualUsecase.NewAddAccrualUseCase(transactionManager, orderService, balanceService)

	accrualJobs := []Job{
		accrual.New(config, sugar, orderService, accrualService, usecaseAddAccrual),
	}

	return NewRunner(accrualJobs)
}
