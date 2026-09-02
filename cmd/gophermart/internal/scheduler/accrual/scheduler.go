package accrual

import (
	"context"
	"errors"
	"kialkuz/shop-with-loyalty/internal/config"
	balanceServ "kialkuz/shop-with-loyalty/internal/domain/balance/service"
	orderServ "kialkuz/shop-with-loyalty/internal/domain/order/service"
	infrastructureInterfaces "kialkuz/shop-with-loyalty/internal/infrastructure/interfaces"
	"sync"

	"go.uber.org/zap"
)

var errorPrefix = errors.New("Error sheduler")

type AccrualScheduler struct {
	config             *config.Config
	sugar              *zap.SugaredLogger
	transactionManager infrastructureInterfaces.TransactionManager
	orderService       *orderServ.OrderService
	balanceService     *balanceServ.BalanceService
	m                  sync.Mutex
}

func New(
	config *config.Config,
	sugar *zap.SugaredLogger,
	transactionManager infrastructureInterfaces.TransactionManager,
	orderService *orderServ.OrderService,
	balanceService *balanceServ.BalanceService,
) *AccrualScheduler {
	return &AccrualScheduler{
		config:             config,
		sugar:              sugar,
		transactionManager: transactionManager,
		orderService:       orderService,
		balanceService:     balanceService,
	}
}

func (sh *AccrualScheduler) Run(ctx context.Context) {
	go sh.GetList(ctx)
}
