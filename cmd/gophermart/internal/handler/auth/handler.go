package auth

import (
	"kialkuz/shop-with-loyalty/internal/config"
	"kialkuz/shop-with-loyalty/internal/service"
	"kialkuz/shop-with-loyalty/pkg/validator"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
	userService *service.UserService
	validator   *validator.Validator
	config      *config.Config
}

func NewAuthHandler(
	authService *service.AuthService,
	userService *service.UserService,
	validator *validator.Validator,
	config *config.Config,
) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		userService: userService,
		validator:   validator,
		config:      config,
	}
}

func (h *AuthHandler) RegisterPublicRoutes(rg *gin.RouterGroup) {
	rg.POST("/user/register", h.Register)
	rg.POST("/user/login", h.Login)
}
