package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// AuthMiddleware validates that a valid Authorization header is present.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader != "Bearer test" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		c.Next()
	}
}
