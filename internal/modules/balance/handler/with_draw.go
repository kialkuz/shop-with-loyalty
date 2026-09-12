package balance

import (
	"errors"
	"kialkuz/shop-with-loyalty/internal/helper"
	requestDto "kialkuz/shop-with-loyalty/internal/modules/balance/dto/request"
	drawalModel "kialkuz/shop-with-loyalty/internal/modules/drawal/model"
	orderModel "kialkuz/shop-with-loyalty/internal/modules/order/model"
	pkgErrors "kialkuz/shop-with-loyalty/pkg/errors"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

func (h *BalanceHandler) WithDraw(c *gin.Context) {
	ctx := c.Request.Context()

	var req requestDto.WithDraw

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validator.ValidateStructDTO(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := helper.CurrentUserID(c)

	drawal := drawalModel.Drawal{
		ID:          uuid.New(),
		UserID:      userID,
		OrderNumber: orderModel.NewNumber(req.Order),
		Sum:         int(req.Sum * 100),
		ProcessedAt: time.Now(),
	}

	if !drawal.OrderNumber.CheckWithLuna() {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": pkgErrors.ErrInvalidOrderNumberFormat.Error()})
		return
	}

	err := h.balanceUsecase.Execute(ctx, userID, drawal)
	if err != nil {
		if errors.Is(err, pkgErrors.ErrLessDrawals) {
			c.JSON(http.StatusPaymentRequired, gin.H{"error": err.Error()})
			return
		}

		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
