package middleware

import (
	"errors"
	tokenServ "kialkuz/shop-with-loyalty/internal/auth/service"
	pkgErrors "kialkuz/shop-with-loyalty/pkg/errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func WithAuth(tokenService *tokenServ.TokenService, secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.Error(pkgErrors.ErrInvalidAuthorizationHeader)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": pkgErrors.ErrAuthNotValid.Error(),
			})
			return
		}

		parts := strings.Fields(authHeader)

		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.Error(pkgErrors.ErrInvalidAuthorizationHeader)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": pkgErrors.ErrAuthNotValid.Error(),
			})
			return
		}

		tokenString := parts[1]

		claims, err := tokenService.ParseTokenString(tokenString, secretKey)
		if err != nil {
			if !errors.Is(err, pkgErrors.ErrNotFound) {
				c.Error(err)
			}

			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": pkgErrors.ErrAuthNotValid.Error()})
			return
		}

		mToken, err := tokenService.GetTokenByJti(ctx, claims.Jti)
		if err != nil {
			if !errors.Is(err, pkgErrors.ErrNotFound) {
				c.Error(err)
			}

			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": pkgErrors.ErrAuthNotValid.Error()})
			return
		}

		if mToken.UserID != claims.UserID {
			c.Error(pkgErrors.ErrTokenisBelongsToAnotherUser)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": pkgErrors.ErrAuthNotValid.Error()})
			return
		}

		if mToken.ExpiresAt.Before(time.Now()) {
			c.Error(pkgErrors.ErrTokenIsExpired)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": pkgErrors.ErrAuthNotValid.Error()})
			return
		}

		c.Set("user_id", mToken.UserID)
		c.Next()
	}
}
