package handler

import (
	"kialkuz/gophermart/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
}

func NewAuthHandler(
	service *service.AuthService,
) *AuthHandler {
	return &AuthHandler{}
}

func (h *AuthHandler) Register(c *gin.Context) {

}

func (h *AuthHandler) Login(c *gin.Context) {

}

func (h *AuthHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/user/register", h.Register)
	rg.POST("/user/login", h.Login)
}
