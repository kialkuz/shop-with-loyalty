package balance

import (
	balanceService "kialkuz/shop-with-loyalty/internal/domain/balance/service"
	orderService "kialkuz/shop-with-loyalty/internal/domain/order/service"

	"kialkuz/shop-with-loyalty/pkg/validator"

	"github.com/gin-gonic/gin"
)

type BalanceHandler struct {
	validator      *validator.Validator
	orderService   *orderService.OrderService
	balanceService *balanceService.BalanceService
}

func NewBalanceHandler(
	validator *validator.Validator,
	orderService *orderService.OrderService,
	balanceService *balanceService.BalanceService,
) *BalanceHandler {
	return &BalanceHandler{
		validator:      validator,
		orderService:   orderService,
		balanceService: balanceService,
	}
}

func (h *BalanceHandler) RegisterProtectedRoutes(rg *gin.RouterGroup) {
	rg.GET("/user/balance", h.GetBalance)
	rg.POST("/user/balance/withdraw", h.WithDraw)
}
