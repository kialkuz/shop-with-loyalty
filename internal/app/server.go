package app

import (
	"kialkuz/shop-with-loyalty/internal/config"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func NewServer(
	config *config.Config,
	router *gin.Engine,
) *http.Server {
	return &http.Server{
		Addr:         ":" + config.ServerPort,
		Handler:      router,
		ReadTimeout:  time.Duration(5 * time.Second),
		WriteTimeout: time.Duration(10 * time.Second),
		IdleTimeout:  time.Duration(15 * time.Second),
	}
}
