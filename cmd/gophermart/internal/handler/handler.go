package handler

import (
	"kialkuz/shop-with-loyalty/internal/config"
	"kialkuz/shop-with-loyalty/internal/handler/auth"
	"kialkuz/shop-with-loyalty/internal/handler/order"
	"kialkuz/shop-with-loyalty/internal/service"
	"kialkuz/shop-with-loyalty/pkg/validator"
)

type Handler struct {
	Auth    *auth.AuthHandler
	Balance *BalanceHandler
	Orders  *order.OrdersHandler
}

func NewHandler(
	config *config.Config,
	authService *service.AuthService,
	balanceService *service.BalanceService,
	ordersService *service.OrderService,
	userService *service.UserService,
) *Handler {
	v := validator.NewValidator()

	return &Handler{
		Auth: auth.NewAuthHandler(
			authService,
			userService,
			v,
			config,
		),
		Balance: NewBalanceHandler(
			balanceService,
		),
		Orders: order.NewOrdersHandler(
			ordersService,
			userService,
		),
	}
}
