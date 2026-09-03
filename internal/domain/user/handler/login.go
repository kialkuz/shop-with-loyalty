package handler

import (
	"errors"
	requestDto "kialkuz/shop-with-loyalty/internal/domain/user/dto/request"
	responceDto "kialkuz/shop-with-loyalty/internal/domain/user/dto/responce"
	pkgErrors "kialkuz/shop-with-loyalty/pkg/errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func (h *UserHandler) Login(c *gin.Context) {
	ctx := c.Request.Context()

	var req requestDto.AuthUser

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validator.ValidateStructDTO(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.GetByLogin(ctx, req.Login)
	if err != nil {
		if !errors.Is(err, pkgErrors.ErrNotFound) {
			c.Error(err)
		}

		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid_password"})
			return
		}

		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": pkgErrors.ErrWriteToSupport})
		return
	}

	token := h.authService.MakeToken(user.ID, TokenExp)

	tokenString, err := h.authService.GenerateTokenString(ctx, token, h.config.SecretKey)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": pkgErrors.ErrWriteToSupport})
		return
	}

	err = h.authService.AddNewToken(ctx, token)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": pkgErrors.ErrWriteToSupport})
		return
	}

	responseDto := responceDto.Login{
		Token: tokenString,
	}

	c.JSON(http.StatusOK, responseDto)
}
