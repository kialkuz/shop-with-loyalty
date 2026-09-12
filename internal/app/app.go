package app

import (
	"context"

	tokenRepository "kialkuz/shop-with-loyalty/internal/auth/repository"
	tokenServ "kialkuz/shop-with-loyalty/internal/auth/service"
	"kialkuz/shop-with-loyalty/internal/config"
	"kialkuz/shop-with-loyalty/internal/infrastructure"
	"kialkuz/shop-with-loyalty/internal/middleware"
	accrualServ "kialkuz/shop-with-loyalty/internal/modules/accrual/service"
	balanceRepository "kialkuz/shop-with-loyalty/internal/modules/balance/repository"
	balanceServ "kialkuz/shop-with-loyalty/internal/modules/balance/service"
	drawalRepository "kialkuz/shop-with-loyalty/internal/modules/drawal/repository"
	drawalServ "kialkuz/shop-with-loyalty/internal/modules/drawal/service"
	orderRepo "kialkuz/shop-with-loyalty/internal/modules/order/repository"
	orderServ "kialkuz/shop-with-loyalty/internal/modules/order/service"
	userRepository "kialkuz/shop-with-loyalty/internal/modules/user/repository"
	userServ "kialkuz/shop-with-loyalty/internal/modules/user/service"
	"kialkuz/shop-with-loyalty/internal/router"
	"kialkuz/shop-with-loyalty/internal/scheduler"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

type App struct {
	Router          *gin.Engine
	Pool            *pgxpool.Pool
	SchedulerRunner *scheduler.Runner
}

func NewApp(ctx context.Context, pool *pgxpool.Pool, sugar *zap.SugaredLogger, config *config.Config) (*App, error) {
	dbStorage := infrastructure.NewDB(pool)

	infrastructure.RunMigrations("pgx", config.DB.DatabaseURI)

	transactionManager := infrastructure.NewTransactionManager(pool)

	tokenService := tokenServ.NewTokenService(tokenRepository.NewTokenRepository(dbStorage))
	userService := userServ.NewUserService(userRepository.NewUserRepository(dbStorage))
	balanceService := balanceServ.NewBalanceService(balanceRepository.NewBalanceRepository(dbStorage))
	drawalService := drawalServ.NewDrawalService(drawalRepository.NewDrawalRepository(dbStorage))
	orderService := orderServ.NewOrderService(orderRepo.NewOrderRepository(dbStorage))
	accrualService := accrualServ.NewAccrualService(config.AccrualSystemAddress)

	appHandler := NewHandler(
		config,
		transactionManager,
		tokenService,
		balanceService,
		orderService,
		userService,
		drawalService,
		accrualService,
	)

	schedulerRunner := scheduler.New(
		config,
		sugar,
		transactionManager,
		orderService,
		balanceService,
		accrualService,
	)

	return &App{
		Router:          getRouter(appHandler, sugar, tokenService, config.SecretKey),
		Pool:            pool,
		SchedulerRunner: schedulerRunner,
	}, nil
}

func getRouter(
	appHandler *Handler,
	sugar *zap.SugaredLogger,
	tokenService *tokenServ.TokenService,
	secretKey string,
) *gin.Engine {
	return router.Init(
		[]router.PublicRouteRegistrar{
			appHandler.User,
		},
		[]router.ProtectedRouteRegistrar{
			appHandler.Balance,
			appHandler.Order,
			appHandler.Drawal,
		},
		sugar,
		middleware.WithAuth(tokenService, secretKey),
	)
}
