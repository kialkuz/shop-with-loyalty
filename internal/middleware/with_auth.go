package middleware

import (
	"errors"
	"fmt"
	authService "kialkuz/shop-with-loyalty/internal/auth/service"
	pkgErrors "kialkuz/shop-with-loyalty/pkg/errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func WithAuth(authService *authService.AuthService, secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			fmt.Println("111111111111111111")
			c.Error(pkgErrors.ErrInvalidAuthorizationHeader)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": pkgErrors.ErrAuthNotValid.Error(),
			})
			return
		}

		parts := strings.Fields(authHeader)

		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			fmt.Println("2222222222222222222")
			c.Error(pkgErrors.ErrInvalidAuthorizationHeader)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": pkgErrors.ErrAuthNotValid.Error(),
			})
			return
		}

		tokenString := parts[1]

		claims, err := authService.ParseTokenString(tokenString, secretKey)
		if err != nil {
			fmt.Println("333333333333333333333")
			if !errors.Is(err, pkgErrors.ErrNotFound) {
				c.Error(err)
			}

			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": pkgErrors.ErrAuthNotValid.Error()})
			return
		}

		mToken, err := authService.GetTokenByJti(ctx, claims.Jti)
		if err != nil {
			fmt.Println("4444444444444444444")
			if !errors.Is(err, pkgErrors.ErrNotFound) {
				c.Error(err)
			}

			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": pkgErrors.ErrAuthNotValid.Error()})
			return
		}

		if mToken.UserID != claims.UserID {
			fmt.Println("55555555555555555555")
			c.Error(pkgErrors.ErrTokenisBelongsToAnotherUser)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": pkgErrors.ErrAuthNotValid.Error()})
			return
		}

		if mToken.ExpiresAt.Before(time.Now()) {
			fmt.Println("666666666666666666666")
			c.Error(pkgErrors.ErrTokenIsExpired)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": pkgErrors.ErrAuthNotValid.Error()})
			return
		}

		c.Set("user_id", mToken.UserID)
		c.Next()
	}
}
