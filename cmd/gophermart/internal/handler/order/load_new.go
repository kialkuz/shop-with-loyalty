package order

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"kialkuz/shop-with-loyalty/internal/helper"
	modelOrder "kialkuz/shop-with-loyalty/internal/model/order"
	pkgErrors "kialkuz/shop-with-loyalty/pkg/errors"
)

func (h *OrdersHandler) LoadNew(c *gin.Context) {
	ctx := c.Request.Context()

	body, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read body"})
		return
	}

	number := strings.TrimSpace(string(body))

	if number == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "empty number"})
		return
	}

	userId := helper.CurrentUserId(c)

	newOrder := modelOrder.Order{
		ID:         uuid.New(),
		UserId:     userId,
		Number:     modelOrder.NewNumber(number),
		Status:     modelOrder.NewStatus(modelOrder.STATUS_PROCESSING),
		Accrual:    nil,
		UploadedAt: time.Now(),
	}

	if !newOrder.Number.CheckWithLuna() {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "invalid order number format"})
		return
	}

	existOrder, err := h.orderService.GetByNumber(ctx, number)
	if err != nil && !errors.Is(err, pkgErrors.ErrNotFound) {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if existOrder != nil {
		if existOrder.IsLoadedByUser(userId) {
			c.JSON(http.StatusConflict, gin.H{"error": "order is loaded by this user"})
		} else {
			c.JSON(http.StatusConflict, gin.H{"error": "order is loaded by another user"})
		}
		return
	}

	err = h.orderService.Add(ctx, newOrder)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{})
}
