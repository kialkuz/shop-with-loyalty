package accrual

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	balanceModel "kialkuz/shop-with-loyalty/internal/domain/balance/model"
	orderModel "kialkuz/shop-with-loyalty/internal/domain/order/model"
	accrualDto "kialkuz/shop-with-loyalty/internal/scheduler/accrual/dto"
	pkgErrors "kialkuz/shop-with-loyalty/pkg/errors"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

const PeriodForWait = 1

var ordersList map[string]orderModel.Order

var locked = false

var num int
var needRetryAddr int64

func (sh *AccrualScheduler) GetList(ctx context.Context) {
	var err error
	var usersID []uuid.UUID

	for {
		ordersList, err = sh.orderService.GetOrdersForGetAccrual(ctx)
		if err != nil {
			if errors.Is(err, pkgErrors.ErrNotFound) {
				continue
			}

			sh.sugar.Error(err)
			return
		}

		numJobs := len(ordersList)
		if numJobs == 0 {
			ticker := time.NewTimer(PeriodForWait * time.Minute)
			<-ticker.C
			continue
		}

		for _, order := range ordersList {
			usersID = append(usersID, order.UserID)
		}

		usersBalance, err := sh.balanceService.GetByUsersID(ctx, usersID)
		if err != nil {
			sh.sugar.Error(err)
			return
		}

		jobs := make(chan string, numJobs)
		results := make(chan accrualDto.GetAccrual, numJobs)

		client := &http.Client{}

		c := sync.NewCond(&sh.m)

		for i := 0; i < numJobs; i++ {
			go sh.getByOrderNumber(ctx, c, client, jobs, results)
		}

		for _, order := range ordersList {
			jobs <- order.Number.Value
		}
		close(jobs)

		for i := 0; i < numJobs; i++ {
			orderAccrual := <-results
			if orderAccrual.Err != nil {
				sh.sugar.Error(err)
				return
			}

			orderForUpdate := ordersList[orderAccrual.Number]
			status, err := orderModel.CreateStatusFromView(orderAccrual.Status)
			if err != nil {
				sh.sugar.Error(err)
				return
			}

			accrual := int(orderAccrual.Accrual * 100)

			orderForUpdate.Status = status
			orderForUpdate.Accrual = &accrual

			userBalance := usersBalance[orderForUpdate.UserID]
			userBalance.Current += accrual

			usersBalance[orderForUpdate.UserID] = userBalance

			err = sh.updateAccrual(ctx, orderForUpdate, userBalance)
			if err != nil {
				sh.sugar.Error(err)
				return
			}
		}
		close(results)
	}
}

func (sh *AccrualScheduler) getByOrderNumber(
	ctx context.Context,
	c *sync.Cond,
	client *http.Client,
	jobs chan string,
	results chan<- accrualDto.GetAccrual,
) {
	var number string
	var retryNum int
	getErrPrefix := fmt.Errorf("%v %v", time.Now().Format("2006/01/02 15:04:05"), errorPrefix)

	number = <-jobs

	for {
		var timer *time.Timer

		needRetry := atomic.LoadInt64(&needRetryAddr)

		if needRetry == 1 {
			// Устанавливаем таймер для ожидания, пока сервис начисления баллов будет снова доступен.
			// locked, чтобы только одна горутина создавала таймер, т.к. горутин может быть много.
			sh.m.Lock()
			if !locked {
				timer = time.NewTimer(time.Duration(retryNum)*time.Second + (5 * time.Second))
				locked = true
			}
			sh.m.Unlock()

			if timer == nil {
				// Для горутин, которые просто ждут без таймера устанавливаем блокировку выполнения
				c.L.Lock()
				c.Wait()
				c.L.Unlock()

				if ctx.Err() != nil {
					results <- accrualDto.GetAccrual{
						Err: fmt.Errorf("%v %v", getErrPrefix, ctx.Err()),
					}
					return
				}
			} else {
				// Для горутины с таймером ждем, пока время ожидания внешнего сервиса пройдет
				// и освобождаем все остальные горутины
				select {
				case <-timer.C:
					locked = false
					atomic.StoreInt64(&needRetryAddr, 0)
					c.Broadcast()
				case <-ctx.Done():
					timer.Stop()
					c.Broadcast()

					if ctx.Err() != nil {
						results <- accrualDto.GetAccrual{
							Err: fmt.Errorf("%v %v", getErrPrefix, ctx.Err()),
						}
						return
					}
					return
				}
			}
		}

		if number == "" {
			return
		}

		req, _ := http.NewRequestWithContext(
			ctx,
			http.MethodGet,
			"http://"+sh.config.AccrualServerHost+":"+sh.config.AccrualServerPort+"/api/orders/"+number,
			nil,
		)
		resp, err := client.Do(req)
		if err != nil {
			results <- accrualDto.GetAccrual{
				Err: fmt.Errorf("%v %v", getErrPrefix, err),
			}
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusTooManyRequests {
			retryNum, err = strconv.Atoi(resp.Header.Get("Retry-After"))
			if err != nil {
				results <- accrualDto.GetAccrual{
					Err: fmt.Errorf("%v %v", getErrPrefix, err),
				}
				return
			}

			atomic.CompareAndSwapInt64(&needRetryAddr, 0, 1)
			continue
		}

		var orderAccrual accrualDto.GetAccrual

		if err := json.NewDecoder(resp.Body).Decode(&orderAccrual); err != nil {
			results <- accrualDto.GetAccrual{
				Err: fmt.Errorf("%v %v", getErrPrefix, err),
			}
			return
		}

		results <- orderAccrual

		return
	}
}

func (sh *AccrualScheduler) updateAccrual(
	ctx context.Context,
	orderForUpdate orderModel.Order,
	userBalance balanceModel.Balance,
) error {
	return sh.transactionManager.WithinTransaction(ctx, func(tx pgx.Tx) error {
		err := sh.orderService.UpdateTx(ctx, tx, orderForUpdate)
		if err != nil {
			return err
		}

		err = sh.balanceService.UpdateTx(ctx, tx, userBalance)
		if err != nil {
			return err
		}

		return nil
	})
}
