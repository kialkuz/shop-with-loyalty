package app

import (
	"kialkuz/shop-with-loyalty/internal/config"
	balance "kialkuz/shop-with-loyalty/internal/domain/balance/handler"
	drawal "kialkuz/shop-with-loyalty/internal/domain/drawal/handler"
	order "kialkuz/shop-with-loyalty/internal/domain/order/handler"
	user "kialkuz/shop-with-loyalty/internal/domain/user/handler"

	authServ "kialkuz/shop-with-loyalty/internal/auth/service"
	balanceServ "kialkuz/shop-with-loyalty/internal/domain/balance/service"
	drawalServ "kialkuz/shop-with-loyalty/internal/domain/drawal/service"
	orderServ "kialkuz/shop-with-loyalty/internal/domain/order/service"
	userServ "kialkuz/shop-with-loyalty/internal/domain/user/service"

	"kialkuz/shop-with-loyalty/pkg/validator"
)

type Handler struct {
	User    *user.UserHandler
	Balance *balance.BalanceHandler
	Order   *order.OrderHandler
	Drawal  *drawal.DrawalHandler
}

func NewHandler(
	config *config.Config,
	authService *authServ.AuthService,
	balanceService *balanceServ.BalanceService,
	orderService *orderServ.OrderService,
	userService *userServ.UserService,
	drawalService *drawalServ.DrawalService,
	accrualService *orderServ.AccrualService,
) *Handler {
	v := validator.NewValidator()

	return &Handler{
		User: user.NewAuthHandler(
			authService,
			userService,
			v,
			config,
		),
		Balance: balance.NewBalanceHandler(
			v,
			orderService,
			balanceService,
		),
		Order: order.NewOrderHandler(
			config,
			orderService,
			accrualService,
			userService,
			balanceService,
		),
		Drawal: drawal.NewDrawalHandler(
			drawalService,
		),
	}
}
