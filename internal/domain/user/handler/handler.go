package handler

import (
	authService "kialkuz/shop-with-loyalty/internal/auth/service"
	"kialkuz/shop-with-loyalty/internal/config"
	userService "kialkuz/shop-with-loyalty/internal/domain/user/service"
	"kialkuz/shop-with-loyalty/pkg/validator"
	"time"

	"github.com/gin-gonic/gin"
)

const TokenExp = time.Hour * 3

type UserHandler struct {
	authService *authService.AuthService
	userService *userService.UserService
	validator   *validator.Validator
	config      *config.Config
}

func NewAuthHandler(
	authService *authService.AuthService,
	userService *userService.UserService,
	validator *validator.Validator,
	config *config.Config,
) *UserHandler {
	return &UserHandler{
		authService: authService,
		userService: userService,
		validator:   validator,
		config:      config,
	}
}

func (h *UserHandler) RegisterPublicRoutes(rg *gin.RouterGroup) {
	rg.POST("/user/register", h.Register)
	rg.POST("/user/login", h.Login)
}
