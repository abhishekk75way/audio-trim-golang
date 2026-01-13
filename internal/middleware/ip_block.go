package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func IPBlock(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {

		// ADMIN BYPASS
		if role, ok := c.Get("role"); ok && role == "admin" {
			c.Next()
			return
		}

		log.Println("IP BLOCK ROLE IN MIDDLEWARE:", c.GetString("role"))

		ip := c.ClientIP()

		blocked, _ := rdb.Exists(c, "blocked_ip:"+ip).Result()
		if blocked == 1 {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "Too many requests were received from this IP. Access is temporarily blocked. Please try again later.",
			})
			return
		}

		c.Next()
	}
}
