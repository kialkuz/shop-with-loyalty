package accrual

import (
	"context"
	"errors"
	"fmt"
	orderModel "kialkuz/shop-with-loyalty/internal/modules/order/model"
	pkgErrors "kialkuz/shop-with-loyalty/pkg/errors"
	"time"
)

const PeriodForWait = 1

var ordersList map[string]orderModel.Order

var locked = false

var num int
var needRetryAddr int64

func (sh *AccrualScheduler) GetList(ctx context.Context) {
	for {
		ordersList, err := sh.orderService.GetOrdersForGetAccrual(ctx)
		if err != nil {
			if errors.Is(err, pkgErrors.ErrNotFound) {
				continue
			}

			sh.sugar.Error(fmt.Errorf("%v %v", sh.getError(), err))
			return
		}

		// Если нет необработанных заказов - ждем минуту
		numJobs := len(ordersList)
		if numJobs == 0 {
			ticker := time.NewTimer(PeriodForWait * time.Minute)
			<-ticker.C
			continue
		}

		ordersForUpdate, errList := sh.accrualService.GetFromAccrualSystem(ctx, ordersList)
		for err := range errList {
			sh.sugar.Error(fmt.Errorf("%v %v", sh.getError(), err))
		}

		for _, orderForUpdate := range ordersForUpdate {
			order := ordersList[orderForUpdate.Number]
			order.Number = orderModel.NewNumber(orderForUpdate.Number)

			status, err := orderModel.NewStatus(orderForUpdate.Status)
			if err != nil {
				sh.sugar.Error(fmt.Errorf("%v %v", sh.getError(), err))
			}
			order.Status = status
			order.Accrual = &orderForUpdate.Accrual

			err = sh.usecaseAddAccrual.Execute(ctx, order)
			if err != nil {
				sh.sugar.Error(fmt.Errorf("%v %v", sh.getError(), err))
			}
		}
	}
}
