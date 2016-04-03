package ginx

import (
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

// Gin logging middleware that logs to the standard lib log package
// Based on gin.Logger
func StdLogLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path

		// Process request
		c.Next()

		// Stop timer
		end := time.Now()
		latency := end.Sub(start)

		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()

		comment := c.Errors.ByType(gin.ErrorTypePrivate).String()

		log.Printf(
			"[GIN] | %3d | %13v | %s | %-7s %s\n%s",
			statusCode,
			latency,
			clientIP,
			method,
			path,
			comment,
		)
	}
}
