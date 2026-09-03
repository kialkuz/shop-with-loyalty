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

		fmt.Println("11111111111111111111")
		fmt.Println(authHeader)
		if authHeader == "" {
			c.Error(pkgErrors.ErrInvalidAuthorizationHeader)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": pkgErrors.ErrAuthNotValid,
			})
			return
		}

		parts := strings.Fields(authHeader)

		fmt.Println("22222222222222222222")
		fmt.Println(parts)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.Error(pkgErrors.ErrInvalidAuthorizationHeader)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": pkgErrors.ErrAuthNotValid,
			})
			return
		}

		tokenString := parts[1]

		fmt.Println("3333333333333333333333")
		fmt.Println(tokenString)
		fmt.Println(secretKey)
		claims, err := authService.ParseTokenString(tokenString, secretKey)
		if err != nil {
			if !errors.Is(err, pkgErrors.ErrNotFound) {
				c.Error(err)
			}

			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": pkgErrors.ErrAuthNotValid})
			return
		}

		fmt.Println("44444444444444444444")
		mToken, err := authService.GetTokenByJti(ctx, claims.Jti)
		if err != nil {
			if !errors.Is(err, pkgErrors.ErrNotFound) {
				c.Error(err)
			}

			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": pkgErrors.ErrAuthNotValid})
			return
		}

		fmt.Println("5555555555555555555555")
		fmt.Println(mToken.UserID)
		fmt.Println(claims.UserID)
		if mToken.UserID != claims.UserID {
			c.Error(pkgErrors.ErrTokenisBelongsToAnotherUser)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": pkgErrors.ErrAuthNotValid})
			return
		}

		fmt.Println("6666666666666666666")
		if mToken.ExpiresAt.Before(time.Now()) {
			c.Error(pkgErrors.ErrTokenIsExpired)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": pkgErrors.ErrAuthNotValid})
			return
		}

		c.Set("user_id", mToken.UserID)
		c.Next()
	}
}
