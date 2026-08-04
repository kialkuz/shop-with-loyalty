package handler

import "kialkuz/gophermart/internal/service"

type Handler struct {
	Auth    *AuthHandler
	Balance *BalanceHandler
	Orders  *OrdersHandler
}

func NewHandler(
	authService *service.AuthService,
	balanceService *service.BalanceService,
	ordersService *service.OrdersService,
) *Handler {
	return &Handler{
		Auth:    NewAuthHandler(authService),
		Balance: NewBalanceHandler(balanceService),
		Orders:  NewOrdersHandler(ordersService),
	}
}
