package utils

import (
	"github.com/gin-gonic/gin"
)

// CorsMiddleware sets the necessary headers to support Cross-Origin Resource Sharing (CORS).
func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Allow all origins, methods, and headers
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204) 
			return
		}
		
		// Continue to the next middleware
		c.Next()
	}
}