package middleware

import "github.com/gin-gonic/gin"

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get("role")
		if role != "admin" {
			c.AbortWithStatusJSON(403, gin.H{"error": "admin access required"})
			return
		}
		c.Next()
	}
}
