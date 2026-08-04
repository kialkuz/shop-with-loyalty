package app

import (
	"context"

	"kialkuz/gophermart/internal/config"
	"kialkuz/gophermart/internal/infrastructure/repository/db"
	"kialkuz/gophermart/internal/infrastructure/storage"
	"kialkuz/gophermart/internal/service"

	// "kialkuz/gophermart/internal/model"
	"kialkuz/gophermart/internal/handler"
	// pkgContracts "kialkuz/service-metrics-and-alerting/pkg/contracts"

	"github.com/jackc/pgx/v5/pgxpool"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type App struct {
	Handler *handler.Handler
	Pool    *pgxpool.Pool
}

func NewApp(ctx context.Context, cfg *config.Config) (*App, error) {
	pool, err := pgxpool.New(ctx, cfg.DB.DatabaseURI)
	if err != nil {
		return nil, err
	}

	dbStorage := storage.NewDB(pool)

	handler := handler.NewHandler(
		service.NewAuthService(db.NewUsersRepository(dbStorage)),
		service.NewBalanceService(db.NewBalanceRepository(dbStorage)),
		service.NewOrdersService(db.NewOrdersRepository(dbStorage)),
	)

	return &App{
		Handler: handler,
		Pool:    pool,
	}, nil
}
