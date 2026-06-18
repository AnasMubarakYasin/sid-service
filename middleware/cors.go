package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CORS returns a Cross-Origin Resource Sharing middleware.
func CORS() gin.HandlerFunc {
	// allowedOrigin := "*"
	allowedOrigin := "http://localhost:3000"
	return func(c *gin.Context) {
		// Allow requests from any origin.
		c.Writer.Header().Set("Access-Control-Allow-Origin", allowedOrigin)

		// Allowed HTTP methods.
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")

		// Allowed headers.
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Content-Length, Accept, Authorization, Connection, Cache-Control")

		// Allow credentials (cookies, auth headers, etc.).
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		// Handle preflight OPTIONS request.
		if c.Request.Method == "OPTIONS" {

			// Cache preflight response for 24 hours.
			c.Writer.Header().Set("Access-Control-Max-Age", "86400")

			// Return 204 No Content for preflight.
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		// Pass control to the next middleware/handler.
		c.Next()
	}
}
