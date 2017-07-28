package ginx

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// Gin logging middleware that logs to the given standard lib Logger.
// Based on gin.Logger
func RequestLogger(logger *log.Logger) gin.HandlerFunc {
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

		logger.Printf(
			"[http] | %3d | %13v | %s | %-7s %q\n%s",
			statusCode,
			latency,
			clientIP,
			method,
			path,
			comment,
		)
	}
}
