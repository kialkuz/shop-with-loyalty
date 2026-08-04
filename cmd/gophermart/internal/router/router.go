package router

import (
	"github.com/gin-gonic/gin"
)

type RouteRegistrar interface {
	RegisterRoutes(*gin.RouterGroup)
}

func Init(handlers ...RouteRegistrar) *gin.Engine {
	router := gin.New()

	api := router.Group("/api")

	for _, h := range handlers {
		h.RegisterRoutes(api)
	}

	return router
}
