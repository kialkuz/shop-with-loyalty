package order

import (
	"context"
	"errors"
	"kialkuz/shop-with-loyalty/internal/helper"
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

	orders, err := h.orderService.GetByUserID(ctx, userID)
	if err != nil {
		if !errors.Is(err, pkgErrors.ErrNotFound) {
			c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}

		c.JSON(http.StatusNoContent, gin.H{"error": pkgErrors.ErrEmptyListOrders.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}
