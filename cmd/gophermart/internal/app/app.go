package app

import (
	"context"

	"kialkuz/shop-with-loyalty/internal/config"
	"kialkuz/shop-with-loyalty/internal/infrastructure/repository/db"
	"kialkuz/shop-with-loyalty/internal/infrastructure/storage"
	"kialkuz/shop-with-loyalty/internal/middleware"
	"kialkuz/shop-with-loyalty/internal/router"
	"kialkuz/shop-with-loyalty/internal/service"

	"kialkuz/shop-with-loyalty/internal/handler"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

type App struct {
	Router *gin.Engine
	Pool   *pgxpool.Pool
}

func NewApp(ctx context.Context, pool *pgxpool.Pool, sugar *zap.SugaredLogger, config *config.Config) (*App, error) {
	dbStorage := storage.NewDB(pool)

	transactionManager := db.NewTransactionManager(pool)

	userRepository := db.NewUserRepository(dbStorage)
	userService := service.NewUserService(userRepository)

	authService := service.NewAuthService(db.NewTokenRepository(dbStorage), userRepository, transactionManager)

	appHandler := handler.NewHandler(
		config,
		authService,
		service.NewBalanceService(db.NewBalanceRepository(dbStorage)),
		service.NewOrderService(db.NewOrderRepository(dbStorage)),
		userService,
	)

	router := router.Init(
		[]router.PublicRouteRegistrar{
			appHandler.Auth,
		},
		[]router.ProtectedRouteRegistrar{
			appHandler.Balance,
			appHandler.Orders,
		},
		sugar,
		middleware.WithAuth(authService, config.SecretKey),
	)

	return &App{
		Router: router,
		Pool:   pool,
	}, nil
}
