package middleware

import (
	"errors"
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
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "",
			})
			return
		}

		parts := strings.Fields(authHeader)

		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid Authorization header",
			})
			return
		}

		tokenString := parts[1]

		claims, err := authService.ParseTokenString(tokenString, secretKey)
		if err != nil {
			if !errors.Is(err, pkgErrors.ErrNotFound) {
				c.Error(err)
			}

			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
			return
		}

		mToken, err := authService.GetTokenByJti(ctx, claims.Jti)
		if mToken.UserId != claims.UserId {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user not valid"})
			return
		}

		if mToken.ExpiresAt.Before(time.Now()) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token is expired"})
			return
		}

		c.Set("user_id", mToken.UserId)
		c.Next()
	}
}
