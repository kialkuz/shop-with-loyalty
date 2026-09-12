package app

import (
	"kialkuz/shop-with-loyalty/internal/config"
	"kialkuz/shop-with-loyalty/internal/infrastructure"
	balance "kialkuz/shop-with-loyalty/internal/modules/balance/handler"
	drawal "kialkuz/shop-with-loyalty/internal/modules/drawal/handler"
	order "kialkuz/shop-with-loyalty/internal/modules/order/handler"
	user "kialkuz/shop-with-loyalty/internal/modules/user/handler"

	tokenServ "kialkuz/shop-with-loyalty/internal/auth/service"
	accrualServ "kialkuz/shop-with-loyalty/internal/modules/accrual/service"
	accrualUsecase "kialkuz/shop-with-loyalty/internal/modules/accrual/usecase"
	balanceServ "kialkuz/shop-with-loyalty/internal/modules/balance/service"
	balanceUsecase "kialkuz/shop-with-loyalty/internal/modules/balance/usecase"
	drawalServ "kialkuz/shop-with-loyalty/internal/modules/drawal/service"
	orderServ "kialkuz/shop-with-loyalty/internal/modules/order/service"
	userServ "kialkuz/shop-with-loyalty/internal/modules/user/service"
	userUsecase "kialkuz/shop-with-loyalty/internal/modules/user/usecase"

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
	transactionManager infrastructure.Transaction,
	tokenService *tokenServ.TokenService,
	balanceService *balanceServ.BalanceService,
	orderService *orderServ.OrderService,
	userService *userServ.UserService,
	drawalService *drawalServ.DrawalService,
	accrualService *accrualServ.AccrualService,
) *Handler {
	v := validator.NewValidator()

	return &Handler{
		User: user.NewUserHandler(
			config,
			v,
			userService,
			tokenService,
			userUsecase.NewRegisterAndLoginUseCase(userService, tokenService, balanceService, transactionManager),
		),
		Balance: balance.NewBalanceHandler(
			v,
			orderService,
			balanceService,
			balanceUsecase.NewWithDrawUseCase(transactionManager, drawalService, balanceService),
		),
		Order: order.NewOrderHandler(
			config,
			orderService,
			accrualService,
			userService,
			accrualUsecase.NewAddAccrualUseCase(transactionManager, orderService, balanceService),
		),
		Drawal: drawal.NewDrawalHandler(
			drawalService,
		),
	}
}
