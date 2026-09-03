package order

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	dtoResponce "kialkuz/shop-with-loyalty/internal/domain/order/dto/responce"
	orderModel "kialkuz/shop-with-loyalty/internal/domain/order/model"
	"kialkuz/shop-with-loyalty/internal/helper"
	pkgErrors "kialkuz/shop-with-loyalty/pkg/errors"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
)

var needRetryAddr int64
var locked = false

type getAccrual struct {
	Number  string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
	Err     error
}

type orderAccrual struct {
	Number  string
	Status  *orderModel.Status
	Accrual int
}

func (h *OrderHandler) GetList(c *gin.Context) {
	ctx := c.Request.Context()

	userID := helper.CurrentUserID(c)

	rows, err := h.orderService.GetByUserID(ctx, userID)
	if err != nil {
		if !errors.Is(err, pkgErrors.ErrNotFound) {
			c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "empty list orders"})
			return
		}

		c.JSON(http.StatusNoContent, gin.H{"error": "empty list orders"})
		return
	}

	var orders []dtoResponce.ViewOrder
	var ordersForGetAccrual = make(map[string]orderModel.Order)

	for _, row := range rows {
		if row.Status.Value == orderModel.StatusProcessed {
			orders = append(orders, dtoResponce.ViewOrder{
				Number:     row.Number.Value,
				Status:     row.Status.GetValueForView(),
				Accrual:    row.Accrual,
				UploadedAt: row.UploadedAt,
			})
		} else {
			ordersForGetAccrual[row.Number.Value] = row
		}
	}

	if len(ordersForGetAccrual) > 0 {
		ordersFromAccrual, accrualErrors := h.getFromAccrualSystem(ctx, ordersForGetAccrual)

		userBalance, err := h.balanceService.GetByUserID(ctx, userID)
		if err != nil {
			if !errors.Is(err, pkgErrors.ErrNotFound) {
				c.Error(err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "empty list orders"})
				return
			}

			c.JSON(http.StatusNoContent, gin.H{"error": "empty list orders"})
			return
		}

		for _, accrualError := range accrualErrors {
			c.Error(accrualError)
		}

		for _, order := range ordersFromAccrual {
			orderForUpdate := ordersForGetAccrual[order.Number]
			orderForUpdate.Status = order.Status
			orderForUpdate.Accrual = &order.Accrual

			userBalance.Current += order.Accrual

			err = h.orderService.UpdateAccrual(ctx, orderForUpdate, *userBalance)
			if err != nil {
				c.Error(err)
				return
			}
		}
	}

	c.JSON(http.StatusOK, orders)
}

func (h *OrderHandler) getFromAccrualSystem(
	ctx context.Context,
	ordersForGetAccrual map[string]orderModel.Order,
) ([]orderAccrual, []error) {
	numJobs := len(ordersForGetAccrual)
	var errors []error

	var m sync.Mutex

	for {
		jobs := make(chan string, numJobs)
		results := make(chan getAccrual, numJobs)

		client := &http.Client{}

		c := sync.NewCond(&m)

		for i := 0; i < numJobs; i++ {
			go h.getByOrderNumber(ctx, &m, c, client, jobs, results)
		}

		for _, order := range ordersForGetAccrual {
			jobs <- order.Number.Value
		}
		close(jobs)

		var ordersForUpdate []orderAccrual
		for i := 0; i < numJobs; i++ {
			orderFromAccrual := <-results
			if orderFromAccrual.Err != nil {
				errors = append(errors, orderFromAccrual.Err)
				continue
			}

			status, err := orderModel.CreateStatusFromView(orderFromAccrual.Status)
			if err != nil {
				errors = append(errors, orderFromAccrual.Err)
				continue
			}

			accrual := int(orderFromAccrual.Accrual * 100)

			var orderForUpdate orderAccrual

			orderForUpdate.Number = orderFromAccrual.Number
			orderForUpdate.Status = status
			orderForUpdate.Accrual = accrual

			ordersForUpdate = append(ordersForUpdate, orderForUpdate)
		}
		close(results)

		return ordersForUpdate, errors
	}
}

func (h *OrderHandler) getByOrderNumber(
	ctx context.Context,
	m *sync.Mutex,
	c *sync.Cond,
	client *http.Client,
	jobs chan string,
	results chan<- getAccrual,
) {
	var number string
	var retryNum int
	getErrPrefix := fmt.Errorf("%v %v", time.Now().Format("2006/01/02 15:04:05"), errors.New("get accrual"))

	number = <-jobs

	for {
		var timer *time.Timer

		needRetry := atomic.LoadInt64(&needRetryAddr)

		if needRetry == 1 {
			// Устанавливаем таймер для ожидания, пока сервис начисления баллов будет снова доступен.
			// locked, чтобы только одна горутина создавала таймер, т.к. горутин может быть много.
			m.Lock()
			if !locked {
				timer = time.NewTimer(time.Duration(retryNum)*time.Second + (5 * time.Second))
				locked = true
			}
			m.Unlock()

			if timer == nil {
				// Для горутин, которые просто ждут без таймера устанавливаем блокировку выполнения
				c.L.Lock()
				c.Wait()
				c.L.Unlock()

				if ctx.Err() != nil {
					results <- getAccrual{
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
						results <- getAccrual{
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
			h.config.AccrualServerHost+":"+h.config.AccrualServerPort+"/api/orders/"+number,
			nil,
		)
		resp, err := client.Do(req)
		if err != nil {
			results <- getAccrual{
				Err: fmt.Errorf("%v %s %v", getErrPrefix, number, err),
			}
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusNoContent {
			results <- getAccrual{
				Err: fmt.Errorf("%v %s %v", getErrPrefix, number, errors.New("order not registered")),
			}
			return
		}

		if resp.StatusCode == http.StatusTooManyRequests {
			retryNum, err = strconv.Atoi(resp.Header.Get("Retry-After"))
			if err != nil {
				results <- getAccrual{
					Err: fmt.Errorf("%v  %s %v", getErrPrefix, number, err),
				}
				return
			}

			atomic.CompareAndSwapInt64(&needRetryAddr, 0, 1)
			continue
		}

		var orderAccrual getAccrual

		if err := json.NewDecoder(resp.Body).Decode(&orderAccrual); err != nil {
			results <- getAccrual{
				Err: fmt.Errorf("%v  %s %v", getErrPrefix, number, err),
			}
			return
		}

		results <- orderAccrual

		return
	}
}
