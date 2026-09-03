package handler

import (
	modelBalance "kialkuz/shop-with-loyalty/internal/domain/balance/model"
	requestDto "kialkuz/shop-with-loyalty/internal/domain/user/dto/request"
	modelUser "kialkuz/shop-with-loyalty/internal/domain/user/model"
	pkgErrors "kialkuz/shop-with-loyalty/pkg/errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *UserHandler) Register(c *gin.Context) {
	ctx := c.Request.Context()

	var req requestDto.AuthUser

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := h.validator.ValidateStructDTO(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	isLoginExists, err := h.userService.CheckExistLogin(ctx, req.Login)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if isLoginExists {
		c.JSON(http.StatusConflict, gin.H{"error": pkgErrors.ErrLoginBusy.Error()})
		return
	}

	hashedPassword, err := h.authService.HashPassword(req.Password)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "incorrect password"})
		return
	}

	var user modelUser.User
	user.ID = uuid.New()
	user.Login = req.Login
	user.Password = string(hashedPassword)

	token := h.authService.MakeToken(user.ID, TokenExp)

	tokenString, err := h.authService.GenerateTokenString(ctx, token, h.config.SecretKey)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error, please write to support"})
		return
	}

	var balance modelBalance.Balance
	balance.ID = uuid.New()
	balance.UserID = user.ID
	balance.Current = 0
	balance.WithDrawn = 0

	err = h.authService.RegisterAndLogin(ctx, user, token, balance)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register"})
		return
	}

	c.Header("Authorization", "Bearer "+tokenString)
	c.JSON(http.StatusOK, gin.H{})
}
