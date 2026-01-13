package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Test: Auth only
func TestAuth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "authenticated",
		"user_id": c.GetInt("user_id"),
		"role":    c.GetString("role"),
	})
}

// Test: IPBlock + RateLimit
func TestProtected(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Protected service hits",
		"role":    c.GetString("role"),
	})
}

// Test: Admin only
func TestAdmin(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "admin access granted",
		"role":    c.GetString("role"),
	})
}
