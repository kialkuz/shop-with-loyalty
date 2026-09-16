package handler

import (
	tokenService "kialkuz/shop-with-loyalty/internal/auth/service"
	"kialkuz/shop-with-loyalty/internal/config"
	userService "kialkuz/shop-with-loyalty/internal/modules/user/service"
	userUsecase "kialkuz/shop-with-loyalty/internal/modules/user/usecase"
	"kialkuz/shop-with-loyalty/pkg/validator"
	"time"

	"github.com/gin-gonic/gin"
)

const TokenExp = time.Hour * 3

type UserHandler struct {
	config       *config.Config
	validator    *validator.Validator
	userService  *userService.UserService
	tokenService *tokenService.TokenService
	userUsecase  *userUsecase.RegisterAndLogin
}

func NewUserHandler(
	config *config.Config,
	validator *validator.Validator,
	userService *userService.UserService,
	tokenService *tokenService.TokenService,
	userUsecase *userUsecase.RegisterAndLogin,

) *UserHandler {
	return &UserHandler{
		config:       config,
		validator:    validator,
		userService:  userService,
		tokenService: tokenService,
		userUsecase:  userUsecase,
	}
}

func (h *UserHandler) RegisterPublicRoutes(rg *gin.RouterGroup) {
	rg.POST("/user/register", h.Register)
	rg.POST("/user/login", h.Login)
}
