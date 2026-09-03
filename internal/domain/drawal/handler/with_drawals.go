package balance

import (
	"errors"
	"fmt"
	dtoResponce "kialkuz/shop-with-loyalty/internal/domain/drawal/dto/responce"
	"kialkuz/shop-with-loyalty/internal/helper"
	pkgErrors "kialkuz/shop-with-loyalty/pkg/errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *DrawalHandler) WithDrawals(c *gin.Context) {
	ctx := c.Request.Context()

	userID := helper.CurrentUserID(c)

	rows, err := h.drawalService.GetByUserID(ctx, userID)
	if err != nil {
		if !errors.Is(err, pkgErrors.ErrNotFound) {
			c.Error(err)
		}

		c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Errorf("error get configuration: %w", err)})
		return
	}

	var drawals []dtoResponce.ViewDrawals
	for _, row := range rows {
		drawals = append(drawals, dtoResponce.ViewDrawals{
			Number:      row.OrderNumber,
			Sum:         float64(row.Sum) / 100,
			ProcessedAt: row.ProcessedAt,
		})
	}

	c.JSON(http.StatusOK, drawals)
}
