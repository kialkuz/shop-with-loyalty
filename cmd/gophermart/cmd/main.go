package main

import (
	"context"
	"fmt"
	"kialkuz/gophermart/internal/app"
	appConfig "kialkuz/gophermart/internal/config"
	"log"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config, err := appConfig.NewConfig()
	if err != nil {
		return fmt.Errorf("error get configuration: %s", err.Error())
	}

	appServer, err := app.NewApp(ctx, config)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if appServer.Pool != nil {
			appServer.Pool.Close()
		}
	}()

	fmt.Println("Server running on port ", config.ServerPort)
	newServer := app.NewServer(appServer.Handler, config)
	err = newServer.ListenAndServe()
	if err != nil {
		return fmt.Errorf("error starting server: %s", err.Error())
	}

	return nil
}
