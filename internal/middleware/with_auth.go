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
			c.Error(pkgErrors.ErrInvalidAuthorizationHeader)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": pkgErrors.ErrAuthNotValid,
			})
			return
		}

		parts := strings.Fields(authHeader)

		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.Error(pkgErrors.ErrInvalidAuthorizationHeader)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": pkgErrors.ErrAuthNotValid,
			})
			return
		}

		tokenString := parts[1]

		claims, err := authService.ParseTokenString(tokenString, secretKey)
		if err != nil {
			if !errors.Is(err, pkgErrors.ErrNotFound) {
				c.Error(err)
			}

			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": pkgErrors.ErrAuthNotValid})
			return
		}

		mToken, err := authService.GetTokenByJti(ctx, claims.Jti)
		if err != nil {
			if !errors.Is(err, pkgErrors.ErrNotFound) {
				c.Error(err)
			}

			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": pkgErrors.ErrAuthNotValid})
			return
		}

		if mToken.UserID != claims.UserID {
			c.Error(pkgErrors.ErrTokenisBelongsToAnotherUser)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": pkgErrors.ErrAuthNotValid})
			return
		}

		if mToken.ExpiresAt.Before(time.Now()) {
			c.Error(pkgErrors.ErrTokenIsExpired)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": pkgErrors.ErrAuthNotValid})
			return
		}

		c.Set("user_id", mToken.UserID)
		c.Next()
	}
}
