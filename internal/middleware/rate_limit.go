package middleware

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

const (
	MaxRequests = 5
	Window      = time.Minute
	BlockTTL    = 2 * time.Minute
)

func isAdmin(c *gin.Context) bool {
	role, exists := c.Get("role")
	return exists && role == "admin"
}

func RateLimit(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {

		// ADMIN BYPASS (must come AFTER auth middleware)
		if isAdmin(c) {
			c.Next()
			return
		}

		log.Println("RATE LIMIT ROLE IN MIDDLEWARE:", c.GetString("role"))

		ip := c.ClientIP()
		ctx := c.Request.Context()

		// Check blocked IP
		blocked, err := rdb.Exists(ctx, "blocked_ip:"+ip).Result()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "rate limiter error",
			})
			return
		}

		if blocked == 1 {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "Too many requests were received from this IP. Access is temporarily blocked. Please try again later.",
			})
			return
		}

		key := "ratelimit:" + ip

		count, err := rdb.Incr(ctx, key).Result()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "rate limiter error",
			})
			return
		}

		if count == 1 {
			_ = rdb.Expire(ctx, key, Window)
		}

		if count > MaxRequests {
			blockIP(ctx, rdb, ip)
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "too many requests, try again later",
			})
			return
		}

		c.Next()
	}
}

func blockIP(ctx context.Context, rdb *redis.Client, ip string) {
	data := map[string]string{
		"reason":     "Rate limit exceeded",
		"blocked_at": time.Now().UTC().Format(time.RFC3339),
	}

	jsonData, _ := json.Marshal(data)
	_ = rdb.Set(ctx, "blocked_ip:"+ip, jsonData, BlockTTL).Err()
}
