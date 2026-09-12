package order

import (
	"context"
	"errors"
	"kialkuz/shop-with-loyalty/internal/helper"
	dtoResponce "kialkuz/shop-with-loyalty/internal/modules/order/dto/responce"
	orderModel "kialkuz/shop-with-loyalty/internal/modules/order/model"
	pkgErrors "kialkuz/shop-with-loyalty/pkg/errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const WaitGetOrdersTime = 10

func (h *OrderHandler) GetList(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), WaitGetOrdersTime*time.Second)
	defer cancel()

	userID := helper.CurrentUserID(c)

	rows, err := h.orderService.GetByUserID(ctx, userID)
	if err != nil {
		if !errors.Is(err, pkgErrors.ErrNotFound) {
			c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}

		c.JSON(http.StatusNoContent, gin.H{"error": pkgErrors.ErrEmptyListOrders.Error()})
		return
	}

	var orders []dtoResponce.ViewOrder
	var ordersForGetAccrual = make(map[string]orderModel.Order)

	for _, row := range rows {
		if row.Status.Value == orderModel.StatusProcessed {
			preparedAccrual := float64(*row.Accrual) / 100

			orders = append(orders, dtoResponce.ViewOrder{
				Number:     row.Number.Value,
				Status:     row.Status.GetValueForView(),
				Accrual:    &preparedAccrual,
				UploadedAt: row.UploadedAt,
			})
		} else {
			ordersForGetAccrual[row.Number.Value] = row
		}
	}

	if len(ordersForGetAccrual) > 0 {
		ordersForUpdate, errList := h.accrualService.GetFromAccrualSystem(ctx, ordersForGetAccrual)
		for _, err := range errList {
			c.Error(err)
		}

		/**
		ЭТОТ КОД ПРОДУБЛИРОВАЛ, Т.К. ВСЕ РАВНО ПОТОМ КАК РЕВЬЮ БУДЕТ ПРИНЯТО ОСТАВЛЮ ОБРАБОТКУ
		НАЧИСЛЕНИЯ БАЛЛОВ ТОЛЬКО В ПЛАНИРОВЩИКЕ
		*/
		for _, orderForUpdate := range ordersForUpdate {
			order := ordersForGetAccrual[orderForUpdate.Number]
			order.Number = orderModel.NewNumber(orderForUpdate.Number)

			status, err := orderModel.NewStatus(orderForUpdate.Status)
			if err != nil {
				c.Error(err)
			}
			order.Status = status
			order.Accrual = &orderForUpdate.Accrual

			err = h.usecaseAddAccrual.Execute(ctx, order)
			if err != nil {
				c.Error(err)
			}
		}

		rows, err = h.orderService.GetByUserID(ctx, userID)
		if err != nil {
			if !errors.Is(err, pkgErrors.ErrNotFound) {
				c.Error(err)
				c.JSON(http.StatusInternalServerError, gin.H{})
				return
			}

			c.JSON(http.StatusNoContent, gin.H{"error": pkgErrors.ErrEmptyListOrders.Error()})
			return
		}

		for _, row := range rows {
			var preparedAccrual float64
			if row.Accrual != nil {
				preparedAccrual = float64(*row.Accrual) / 100
			}

			orders = append(orders, dtoResponce.ViewOrder{
				Number:     row.Number.Value,
				Status:     row.Status.GetValueForView(),
				Accrual:    &preparedAccrual,
				UploadedAt: row.UploadedAt,
			})
		}
	}

	c.JSON(http.StatusOK, orders)
}
