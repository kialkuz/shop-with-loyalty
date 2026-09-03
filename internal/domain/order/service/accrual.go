package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	balanceModel "kialkuz/shop-with-loyalty/internal/domain/balance/model"
	orderModel "kialkuz/shop-with-loyalty/internal/domain/order/model"
	infrastructureInterfaces "kialkuz/shop-with-loyalty/internal/infrastructure/interfaces"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

const WaitAccrualRequest = 2

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

type AccrualService struct {
	m                    sync.Mutex
	accrualSystemAddress string
	transactionManager   infrastructureInterfaces.TransactionManager
	orderService         *OrderService
}

func NewAccrualService(
	accrualSystemAddress string,
	transactionManager infrastructureInterfaces.TransactionManager,
	orderService *OrderService,
) *AccrualService {
	return &AccrualService{
		accrualSystemAddress: accrualSystemAddress,
		transactionManager:   transactionManager,
		orderService:         orderService,
	}
}

func (s *AccrualService) GetAndUpdateAccruals(
	ctx context.Context,
	ordersForGetAccrual map[string]orderModel.Order,
	userBalance *balanceModel.Balance,
) []error {
	var errorsList []error

	ordersFromAccrual, accrualErrors := s.getFromAccrualSystem(ctx, ordersForGetAccrual)

	for _, accrualError := range accrualErrors {
		errorsList = append(errorsList, accrualError)
	}

	for _, order := range ordersFromAccrual {
		orderForUpdate := ordersForGetAccrual[order.Number]
		orderForUpdate.Status = order.Status
		orderForUpdate.Accrual = &order.Accrual

		userBalance.Current += order.Accrual

		err := s.orderService.UpdateAccrual(ctx, orderForUpdate, *userBalance)
		if err != nil {
			errorsList = append(errorsList, err)
		}
	}

	return errorsList
}

func (s *AccrualService) getFromAccrualSystem(
	ctx context.Context,
	ordersForGetAccrual map[string]orderModel.Order,
) ([]orderAccrual, []error) {
	numJobs := len(ordersForGetAccrual)
	var errors []error

	for {
		jobs := make(chan string, numJobs)
		results := make(chan getAccrual, numJobs)

		client := &http.Client{
			Timeout: WaitAccrualRequest * time.Second,
		}

		c := sync.NewCond(&s.m)

		for i := 0; i < numJobs; i++ {
			go s.getByOrderNumber(ctx, c, client, jobs, results)
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

func (s *AccrualService) getByOrderNumber(
	ctx context.Context,
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
			s.m.Lock()
			if !locked {
				timer = time.NewTimer(time.Duration(retryNum)*time.Second + (5 * time.Second))
				locked = true
			}
			s.m.Unlock()

			if timer == nil {
				// Для горутин, которые просто ждут без таймера устанавливаем блокировку выполнения
				c.L.Lock()
				c.Wait()
				c.L.Unlock()

				if ctx.Err() != nil {
					results <- getAccrual{
						Err: fmt.Errorf("%v %s %v", getErrPrefix, number, ctx.Err()),
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
							Err: fmt.Errorf("%v %s %v", getErrPrefix, number, ctx.Err()),
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
			s.accrualSystemAddress+"/api/orders/"+number,
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
					Err: fmt.Errorf("%v %s %v", getErrPrefix, number, err),
				}
				return
			}

			atomic.CompareAndSwapInt64(&needRetryAddr, 0, 1)
			continue
		}

		var orderAccrual getAccrual

		if err := json.NewDecoder(resp.Body).Decode(&orderAccrual); err != nil {
			results <- getAccrual{
				Err: fmt.Errorf("%v %s %v", getErrPrefix, number, err),
			}
			return
		}

		results <- orderAccrual

		return
	}
}
