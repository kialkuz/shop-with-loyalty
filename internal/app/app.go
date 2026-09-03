package app

import (
	"context"

	authRepository "kialkuz/shop-with-loyalty/internal/auth/repository"
	authServ "kialkuz/shop-with-loyalty/internal/auth/service"
	"kialkuz/shop-with-loyalty/internal/config"
	balanceRepository "kialkuz/shop-with-loyalty/internal/domain/balance/repository"
	balanceServ "kialkuz/shop-with-loyalty/internal/domain/balance/service"
	drawalRepository "kialkuz/shop-with-loyalty/internal/domain/drawal/repository"
	drawalServ "kialkuz/shop-with-loyalty/internal/domain/drawal/service"
	orderRepo "kialkuz/shop-with-loyalty/internal/domain/order/repository"
	orderServ "kialkuz/shop-with-loyalty/internal/domain/order/service"
	userRepository "kialkuz/shop-with-loyalty/internal/domain/user/repository"
	userServ "kialkuz/shop-with-loyalty/internal/domain/user/service"
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

	infrastructure.RunMigrations("pgx", config.DB.DatabaseURI)

	transactionManager := infrastructure.NewTransactionManager(pool)

	tokenService := authServ.NewTokenService(authRepository.NewTokenRepository(dbStorage))
	userService := userServ.NewUserService(userRepository.NewUserRepository(dbStorage))
	balanceService := balanceServ.NewBalanceService(balanceRepository.NewBalanceRepository(dbStorage))
	drawalService := drawalServ.NewDrawalService(drawalRepository.NewDrawalRepository(dbStorage))

	authService := authServ.NewAuthService(
		tokenService,
		userService,
		balanceService,
		transactionManager,
	)

	orderRepository := orderRepo.NewOrderRepository(dbStorage)

	orderService := orderServ.NewOrderService(
		transactionManager,
		orderRepository,
		balanceService,
		drawalService,
	)

	accrualService := orderServ.NewAccrualService(
		config.AccrualSystemAddress,
		transactionManager,
		orderService,
	)

	appHandler := NewHandler(
		config,
		authService,
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
	authService *authServ.AuthService,
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
