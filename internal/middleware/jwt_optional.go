package middleware

import (
	"strings"

	"authentication/backend/internal/utils"

	"github.com/gin-gonic/gin"
)

func JWTOptional() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.Next()
			return
		}

		tokenStr := strings.TrimPrefix(auth, "Bearer ")

		claims, err := utils.VerifyToken(tokenStr)
		if err == nil {
			c.Set("user_id", claims.UserID)
			c.Set("role", claims.Role)
		}

		c.Next()
	}
}
