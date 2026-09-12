package balance

import (
	balanceService "kialkuz/shop-with-loyalty/internal/modules/balance/service"
	balanceUsecase "kialkuz/shop-with-loyalty/internal/modules/balance/usecase"
	orderService "kialkuz/shop-with-loyalty/internal/modules/order/service"

	"kialkuz/shop-with-loyalty/pkg/validator"

	"github.com/gin-gonic/gin"
)

type BalanceHandler struct {
	validator      *validator.Validator
	orderService   *orderService.OrderService
	balanceService *balanceService.BalanceService
	balanceUsecase *balanceUsecase.WithDraw
}

func NewBalanceHandler(
	validator *validator.Validator,
	orderService *orderService.OrderService,
	balanceService *balanceService.BalanceService,
	balanceUsecase *balanceUsecase.WithDraw,
) *BalanceHandler {
	return &BalanceHandler{
		validator:      validator,
		orderService:   orderService,
		balanceService: balanceService,
		balanceUsecase: balanceUsecase,
	}
}

func (h *BalanceHandler) RegisterProtectedRoutes(rg *gin.RouterGroup) {
	rg.GET("/user/balance", h.GetBalance)
	rg.POST("/user/balance/withdraw", h.WithDraw)
}
