package handler

import (
	"kialkuz/shop-with-loyalty/internal/service"

	"github.com/gin-gonic/gin"
)

type OrdersHandler struct {
}

func NewOrdersHandler(
	service *service.OrderService,
) *OrdersHandler {
	return &OrdersHandler{}
}

func (h *OrdersHandler) LoadNew(c *gin.Context) {

}

func (h *OrdersHandler) GetList(c *gin.Context) {

}

func (h *OrdersHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/user/orders", h.LoadNew)
	rg.GET("/user/orders", h.GetList)
}
