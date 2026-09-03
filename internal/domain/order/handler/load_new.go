package order

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"kialkuz/shop-with-loyalty/internal/domain/order/model"
	"kialkuz/shop-with-loyalty/internal/helper"
	pkgErrors "kialkuz/shop-with-loyalty/pkg/errors"
)

func (h *OrderHandler) LoadNew(c *gin.Context) {
	ctx := c.Request.Context()

	fmt.Println("1111111111111111111111")
	body, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read body"})
		return
	}

	fmt.Println("2222222222222222222")
	number := strings.TrimSpace(string(body))
	if number == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "empty number"})
		return
	}

	fmt.Println("3333333333333333333333")
	userID := helper.CurrentUserID(c)

	status, err := model.NewStatus(model.StatusNew)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	fmt.Println("4444444444444444444444")

	newOrder := model.Order{
		ID:         uuid.New(),
		UserID:     userID,
		Number:     model.NewNumber(number),
		Status:     status,
		Accrual:    nil,
		UploadedAt: time.Now(),
	}

	if !newOrder.Number.CheckWithLuna() {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "invalid order number format"})
		return
	}

	fmt.Println("5555555555555555555")
	existOrder, err := h.orderService.GetByNumber(ctx, number)
	if err != nil && !errors.Is(err, pkgErrors.ErrNotFound) {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("66666666666666666666666")
	if existOrder != nil {
		if existOrder.IsLoadedByUser(userID) {
			c.JSON(http.StatusConflict, gin.H{"error": "order is loaded by this user"})
		} else {
			c.JSON(http.StatusConflict, gin.H{"error": "order is loaded by another user"})
		}
		return
	}

	fmt.Println("77777777777777777777")
	err = h.orderService.Add(ctx, newOrder)
	if err != nil {
		fmt.Println("888888888888888888")
		fmt.Println(err.Error())

		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("99999999999999999")
	order, err := h.orderService.GetByNumber(ctx, number)
	if err != nil && !errors.Is(err, pkgErrors.ErrNotFound) {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("1212121212121212")

	fmt.Println(order)

	c.JSON(http.StatusAccepted, gin.H{})
}
