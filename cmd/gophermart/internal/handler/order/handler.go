package order

import (
	"kialkuz/shop-with-loyalty/internal/service"

	"github.com/gin-gonic/gin"
)

type OrdersHandler struct {
	orderService *service.OrderService
	userService  *service.UserService
}

func NewOrdersHandler(
	orderService *service.OrderService,
	userService *service.UserService,
) *OrdersHandler {
	return &OrdersHandler{
		orderService,
		userService,
	}
}

func (h *OrdersHandler) RegisterProtectedRoutes(rg *gin.RouterGroup) {
	rg.POST("/user/orders", h.LoadNew)
	rg.GET("/user/orders", h.GetList)
}
