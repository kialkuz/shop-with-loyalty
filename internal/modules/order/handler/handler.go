package order

import (
	"kialkuz/shop-with-loyalty/internal/config"
	accrualServ "kialkuz/shop-with-loyalty/internal/modules/accrual/service"
	accrualUsecase "kialkuz/shop-with-loyalty/internal/modules/accrual/usecase"
	orderServ "kialkuz/shop-with-loyalty/internal/modules/order/service"
	userServ "kialkuz/shop-with-loyalty/internal/modules/user/service"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	config            *config.Config
	orderService      *orderServ.OrderService
	accrualService    *accrualServ.AccrualService
	userService       *userServ.UserService
	usecaseAddAccrual *accrualUsecase.AddAccrual
}

func NewOrderHandler(
	config *config.Config,
	orderService *orderServ.OrderService,
	accrualService *accrualServ.AccrualService,
	userService *userServ.UserService,
	usecaseAddAccrual *accrualUsecase.AddAccrual,
) *OrderHandler {
	return &OrderHandler{
		config,
		orderService,
		accrualService,
		userService,
		usecaseAddAccrual,
	}
}

func (h *OrderHandler) RegisterProtectedRoutes(rg *gin.RouterGroup) {
	rg.POST("/user/orders", h.LoadNew)
	rg.GET("/user/orders", h.GetList)
}
