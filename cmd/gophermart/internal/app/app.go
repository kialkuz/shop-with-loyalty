package app

import (
	"context"

	authRepository "kialkuz/shop-with-loyalty/internal/auth/repository"
	authService "kialkuz/shop-with-loyalty/internal/auth/service"
	"kialkuz/shop-with-loyalty/internal/config"
	balanceRepository "kialkuz/shop-with-loyalty/internal/domain/balance/repository"
	balanceService "kialkuz/shop-with-loyalty/internal/domain/balance/service"
	drawalRepository "kialkuz/shop-with-loyalty/internal/domain/drawal/repository"
	drawalService "kialkuz/shop-with-loyalty/internal/domain/drawal/service"
	orderRepository "kialkuz/shop-with-loyalty/internal/domain/order/repository"
	orderService "kialkuz/shop-with-loyalty/internal/domain/order/service"
	userRepository "kialkuz/shop-with-loyalty/internal/domain/user/repository"
	userService "kialkuz/shop-with-loyalty/internal/domain/user/service"
	"kialkuz/shop-with-loyalty/internal/infrastructure"
	"kialkuz/shop-with-loyalty/internal/middleware"
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

	transactionManager := infrastructure.NewTransactionManager(pool)

	tokenService := authService.NewTokenService(authRepository.NewTokenRepository(dbStorage))
	userService := userService.NewUserService(userRepository.NewUserRepository(dbStorage))
	balanceService := balanceService.NewBalanceService(balanceRepository.NewBalanceRepository(dbStorage))
	drawalService := drawalService.NewDrawalService(drawalRepository.NewDrawalRepository(dbStorage))

	authService := authService.NewAuthService(
		tokenService,
		userService,
		balanceService,
		transactionManager,
	)

	orderService := orderService.NewOrderService(
		transactionManager,
		orderRepository.NewOrderRepository(dbStorage),
		balanceService,
		drawalService,
	)

	appHandler := NewHandler(
		config,
		authService,
		balanceService,
		orderService,
		userService,
		drawalService,
	)

	schedulerRunner := scheduler.New(
		config,
		sugar,
		transactionManager,
		orderService,
		balanceService,
	)

	return &App{
		Router:          getRouter(appHandler, sugar, authService, config.SecretKey),
		Pool:            pool,
		SchedulerRunner: schedulerRunner,
	}, nil
}

func getRouter(
	appHandler *Handler,
	sugar *zap.SugaredLogger,
	authService *authService.AuthService,
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
		middleware.WithAuth(authService, secretKey),
	)
}
