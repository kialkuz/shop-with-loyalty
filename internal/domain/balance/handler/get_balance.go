package balance

import (
	"errors"
	"fmt"
	responceDto "kialkuz/shop-with-loyalty/internal/domain/balance/dto/responce"
	"kialkuz/shop-with-loyalty/internal/helper"
	pkgErrors "kialkuz/shop-with-loyalty/pkg/errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *BalanceHandler) GetBalance(c *gin.Context) {
	ctx := c.Request.Context()

	userID := helper.CurrentUserID(c)

	balance, err := h.balanceService.GetByUserID(ctx, userID)
	if err != nil {
		if !errors.Is(err, pkgErrors.ErrNotFound) {
			c.Error(err)
		}

		c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Errorf("user %v", pkgErrors.ErrNotFound)})
		return
	}

	c.JSON(http.StatusOK, responceDto.ViewBalance{
		Current:   balance.Current,
		WithDrawn: balance.WithDrawn,
	})
}
