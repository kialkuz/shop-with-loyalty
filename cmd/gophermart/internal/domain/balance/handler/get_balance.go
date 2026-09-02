package balance

import (
	"errors"
	responceDto "kialkuz/shop-with-loyalty/internal/domain/balance/dto/responce"
	"kialkuz/shop-with-loyalty/internal/helper"
	pkgErrors "kialkuz/shop-with-loyalty/pkg/errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *BalanceHandler) GetBalance(c *gin.Context) {
	ctx := c.Request.Context()

	userId := helper.CurrentUserId(c)

	balance, err := h.balanceService.GetByUserId(ctx, userId)
	if err != nil {
		if !errors.Is(err, pkgErrors.ErrNotFound) {
			c.Error(err)
		}

		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, responceDto.ViewBalance{
		Current:   balance.Current,
		WithDrawn: balance.WithDrawn,
	})
}
