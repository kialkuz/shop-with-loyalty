package router

import (
	"kialkuz/shop-with-loyalty/internal/middleware"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type PublicRouteRegistrar interface {
	RegisterPublicRoutes(*gin.RouterGroup)
}

type ProtectedRouteRegistrar interface {
	RegisterProtectedRoutes(*gin.RouterGroup)
}

func Init(
	publicHandlers []PublicRouteRegistrar,
	protectedHandlers []ProtectedRouteRegistrar,
	sugar *zap.SugaredLogger,
	authMiddleware gin.HandlerFunc,
) *gin.Engine {
	router := gin.New()

	router.Use(middleware.WithLogging(sugar))

	api := router.Group("/api")

	for _, h := range publicHandlers {
		h.RegisterPublicRoutes(api)
	}

	protected := api.Group("")
	protected.Use(authMiddleware)

	for _, h := range protectedHandlers {
		h.RegisterProtectedRoutes(protected)
	}

	return router
}
