package accrual

import (
	"context"
	"errors"
	"fmt"
	"kialkuz/shop-with-loyalty/internal/config"
	accrualServ "kialkuz/shop-with-loyalty/internal/modules/accrual/service"
	accrualUsecase "kialkuz/shop-with-loyalty/internal/modules/accrual/usecase"
	orderServ "kialkuz/shop-with-loyalty/internal/modules/order/service"
	"sync"
	"time"

	"go.uber.org/zap"
)

var errorPrefix = errors.New("error scheduler")

type AccrualScheduler struct {
	config            *config.Config
	sugar             *zap.SugaredLogger
	orderService      *orderServ.OrderService
	accrualService    *accrualServ.AccrualService
	usecaseAddAccrual *accrualUsecase.AddAccrual
	m                 sync.Mutex
}

func New(
	config *config.Config,
	sugar *zap.SugaredLogger,
	orderService *orderServ.OrderService,
	accrualService *accrualServ.AccrualService,
	usecaseAddAccrual *accrualUsecase.AddAccrual,
) *AccrualScheduler {
	return &AccrualScheduler{
		config:            config,
		sugar:             sugar,
		orderService:      orderService,
		accrualService:    accrualService,
		usecaseAddAccrual: usecaseAddAccrual,
	}
}

func (sh *AccrualScheduler) Run(ctx context.Context) {
	go sh.GetList(ctx)
}

func (sh *AccrualScheduler) getError() error {
	return fmt.Errorf("%v %v", time.Now().Format("2006/01/02 15:04:05"), errorPrefix)
}
