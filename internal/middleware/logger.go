package middleware

import (
	// "time"
	"fmt"
	"os"
	"time"

	"github.com/eyuDark/api-gateway/internal/log"
	"github.com/gin-gonic/gin"
)

func Logger(logger *log.RotatingLogger) gin.HandlerFunc {
    return func(c *gin.Context) {
        // TODO: 1. Record start time
        startTime := time.Now()

        // TODO: 2. Process the request (call next handler)
        c.Next()
        // TODO: 3. Calculate latency
        endTime := time.Now()
        latency := endTime.Sub(startTime)
        // TODO: 4. Get client IP
        clientIp := c.ClientIP()
        // TODO: 5. Get HTTP method
        method := c.Request.Method
        // TODO: 6. Get request path
        path := c.Request.URL.Path
        // TODO: 7. Get HTTP status
        status := c.Writer.Status()
        // TODO: 8. Format and print log
        //format date
        timeStr := time.Now().Format("2006-01-02 15:04:05")
        latencyMs := float64(latency.Microseconds()) / 1000.0

        logLine := fmt.Sprintf("[%s] %s %s %s %d %.1fms\n", 
            timeStr,
            clientIp, method, path, status, latencyMs)
        
        // Write to rotating logger (thread-safe)
        if _, err := logger.Write([]byte(logLine)); err != nil {
            // Fallback to stderr ONLY on failure
            fmt.Fprintf(os.Stderr, "Logger error: %v\n", err)
        }
    }
}