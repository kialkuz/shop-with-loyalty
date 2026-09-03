package order

import (
	"errors"
	dtoResponce "kialkuz/shop-with-loyalty/internal/domain/order/dto/responce"
	"kialkuz/shop-with-loyalty/internal/helper"
	pkgErrors "kialkuz/shop-with-loyalty/pkg/errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

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
	for _, row := range rows {
		orders = append(orders, dtoResponce.ViewOrder{
			Number:     row.Number.Value,
			Status:     row.Status.GetValueForView(),
			Accrual:    row.Accrual,
			UploadedAt: row.UploadedAt,
		})
	}

	c.JSON(http.StatusOK, orders)
}
