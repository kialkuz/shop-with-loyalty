package main

import (
	"context"
	"fmt"
	"kialkuz/shop-with-loyalty/internal/app"
	appConfig "kialkuz/shop-with-loyalty/internal/config"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()
	sugar := logger.Sugar()

	if err := run(sugar); err != nil {
		sugar.Error(err.Error())
	}
}

func run(sugar *zap.SugaredLogger) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config, err := appConfig.NewConfig()
	if err != nil {
		sugar.Error(fmt.Errorf("error get configuration: %s", err.Error()))
		return nil
	}

	pool, err := pgxpool.New(ctx, config.DB.DatabaseURI)
	if err != nil {
		sugar.Error(fmt.Errorf("error connect to db: %s", err.Error()))
		return nil
	}

	appServer, err := app.NewApp(ctx, pool, sugar, config)
	if err != nil {
		sugar.Error(err.Error())
		return nil
	}

	defer func() {
		if appServer.Pool != nil {
			appServer.Pool.Close()
		}
	}()

	sugar.Info(fmt.Printf("Server running on port: %s", config.ServerPort))
	newServer := app.NewServer(config, appServer.Router)
	err = newServer.ListenAndServe()
	if err != nil {
		sugar.Error(fmt.Errorf("error starting server: %s", err.Error()))
		return nil
	}

	return nil
}
