package order

import (
	orderServ "kialkuz/shop-with-loyalty/internal/domain/order/service"
	userServ "kialkuz/shop-with-loyalty/internal/domain/user/service"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderService *orderServ.OrderService
	userService  *userServ.UserService
}

func NewOrderHandler(
	orderService *orderServ.OrderService,
	userService *userServ.UserService,
) *OrderHandler {
	return &OrderHandler{
		orderService,
		userService,
	}
}

func (h *OrderHandler) RegisterProtectedRoutes(rg *gin.RouterGroup) {
	rg.POST("/user/orders", h.LoadNew)
	rg.GET("/user/orders", h.GetList)
}
