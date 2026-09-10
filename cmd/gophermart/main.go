package main

import (
	"context"
	"fmt"
	"kialkuz/shop-with-loyalty/internal/app"
	appConfig "kialkuz/shop-with-loyalty/internal/config"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
		sugar.Error(fmt.Errorf("error get configuration: %w", err))
		return nil
	}

	pool, err := pgxpool.New(ctx, config.DB.DatabaseURI)
	if err != nil {
		sugar.Error(fmt.Errorf("error connect to db: %w", err))
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

	srv := &http.Server{
		Handler: appServer.Router,
		Addr:    ":" + config.ServerPort,
	}

	serverError := make(chan error, 1)

	go func() error {
		sugar.Info("Server running on port ", config.ServerPort)

		err := http.ListenAndServe(":"+config.ServerPort, appServer.Router)
		if err != nil && err != http.ErrServerClosed {
			serverError <- fmt.Errorf("failed to start server %s", err)
		}

		return nil
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	select {
	case <-stop:
		sugar.Info("Shutting down server...")
		ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			return fmt.Errorf("server forced to shutdown: %w", err)
		}
		sugar.Info("Server gracefully shutdown")

		return nil
	case err := <-serverError:
		return err
	}
}
