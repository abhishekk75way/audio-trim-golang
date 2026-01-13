package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		// recover from panics (no app crash)
		defer func() {
			if rec := recover(); rec != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": "internal server error",
				})
			}
		}()

		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last()

		status := http.StatusInternalServerError

		if c.Writer.Status() != http.StatusOK {
			status = c.Writer.Status()
		}

		// send unified response
		c.AbortWithStatusJSON(status, gin.H{
			"error": err.Error(),
		})
	}
}
