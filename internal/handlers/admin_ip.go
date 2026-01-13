package handlers

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func ListBlockedIPs(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		var cursor uint64
		var blockedIPs []map[string]interface{}

		for {
			keys, nextCursor, err := rdb.Scan(ctx, cursor, "blocked_ip:*", 100).Result()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "redis scan failed",
				})
				return
			}

			for _, key := range keys {
				ip := strings.TrimPrefix(key, "blocked_ip:")
				ttl, _ := rdb.TTL(ctx, key).Result()

				blockedIPs = append(blockedIPs, map[string]interface{}{
					"ip":  ip,
					"ttl": int(ttl.Seconds()),
				})
			}

			cursor = nextCursor
			if cursor == 0 {
				break
			}
		}

		// RETURN ARRAY
		c.JSON(http.StatusOK, blockedIPs)
	}
}

func UnblockIP(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.Param("ip")
		rdb.Del(c, "blocked_ip:"+ip)
		c.JSON(200, gin.H{"message": "IP unblocked"})
	}
}
