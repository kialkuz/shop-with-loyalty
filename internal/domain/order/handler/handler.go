package order

import (
	"kialkuz/shop-with-loyalty/internal/config"
	orderServ "kialkuz/shop-with-loyalty/internal/domain/order/service"
	userServ "kialkuz/shop-with-loyalty/internal/domain/user/service"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	config         *config.Config
	orderService   *orderServ.OrderService
	accrualService *orderServ.AccrualService
	userService    *userServ.UserService
}

func NewOrderHandler(
	config *config.Config,
	orderService *orderServ.OrderService,
	accrualService *orderServ.AccrualService,
	userService *userServ.UserService,
) *OrderHandler {
	return &OrderHandler{
		config,
		orderService,
		accrualService,
		userService,
	}
}

func (h *OrderHandler) RegisterProtectedRoutes(rg *gin.RouterGroup) {
	rg.POST("/user/orders", h.LoadNew)
	rg.GET("/user/orders", h.GetList)
}
