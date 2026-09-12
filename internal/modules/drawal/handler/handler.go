package balance

import (
	drawalService "kialkuz/shop-with-loyalty/internal/modules/drawal/service"

	"github.com/gin-gonic/gin"
)

type DrawalHandler struct {
	drawalService *drawalService.DrawalService
}

func NewDrawalHandler(
	drawalService *drawalService.DrawalService,
) *DrawalHandler {
	return &DrawalHandler{
		drawalService: drawalService,
	}
}

func (h *DrawalHandler) RegisterProtectedRoutes(rg *gin.RouterGroup) {
	rg.GET("/user/withdrawals", h.WithDrawals)
}
