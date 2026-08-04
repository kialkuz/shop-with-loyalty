package app

import (
	"kialkuz/gophermart/internal/config"
	"kialkuz/gophermart/internal/handler"
	"kialkuz/gophermart/internal/router"
	"net/http"
	"time"
)

func NewServer(handler *handler.Handler, config *config.Config) *http.Server {
	return &http.Server{
		Addr: ":" + config.ServerPort,
		Handler: router.Init(
			handler.Auth,
			handler.Balance,
			handler.Orders,
		),
		ReadTimeout:  time.Duration(5 * time.Second),
		WriteTimeout: time.Duration(10 * time.Second),
		IdleTimeout:  time.Duration(15 * time.Second),
	}
}
