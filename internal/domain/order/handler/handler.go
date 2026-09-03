package order

import (
	"kialkuz/shop-with-loyalty/internal/config"
	balanceServ "kialkuz/shop-with-loyalty/internal/domain/balance/service"
	orderServ "kialkuz/shop-with-loyalty/internal/domain/order/service"
	userServ "kialkuz/shop-with-loyalty/internal/domain/user/service"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	config         *config.Config
	orderService   *orderServ.OrderService
	userService    *userServ.UserService
	balanceService *balanceServ.BalanceService
}

func NewOrderHandler(
	config *config.Config,
	orderService *orderServ.OrderService,
	userService *userServ.UserService,
	balanceService *balanceServ.BalanceService,
) *OrderHandler {
	return &OrderHandler{
		config,
		orderService,
		userService,
		balanceService,
	}
}

func (h *OrderHandler) RegisterProtectedRoutes(rg *gin.RouterGroup) {
	rg.POST("/user/orders", h.LoadNew)
	rg.GET("/user/orders", h.GetList)
}
