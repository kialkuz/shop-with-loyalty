package scheduler

import (
	"kialkuz/shop-with-loyalty/internal/config"
	balanceServ "kialkuz/shop-with-loyalty/internal/domain/balance/service"
	orderServ "kialkuz/shop-with-loyalty/internal/domain/order/service"
	"kialkuz/shop-with-loyalty/internal/infrastructure"
	"kialkuz/shop-with-loyalty/internal/scheduler/accrual"

	"go.uber.org/zap"
)

func New(
	config *config.Config,
	sugar *zap.SugaredLogger,
	transactionManager infrastructure.Transaction,
	ordersService *orderServ.OrderService,
	balanceService *balanceServ.BalanceService,
) *Runner {
	jobs := []Job{
		accrual.New(config, sugar, transactionManager, ordersService, balanceService),
	}

	return NewRunner(jobs)
}
