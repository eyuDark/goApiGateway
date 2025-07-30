package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/eyuDark/api-gateway/internal/log"
	"github.com/eyuDark/api-gateway/internal/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	if err := os.MkdirAll("logs", 0755); err != nil {
		panic(fmt.Sprintf("Failed to create logs directory: %v", err))
	}

	rotatingLogger, err := log.NewRotatingLogger("logs/gateway.log", 10)
	if err != nil {
		panic(fmt.Sprintf("Logger init failed: %v\n"+
			"Hint: Check if 'logs' directory exists and is writable", err))
	}
	defer rotatingLogger.Close()
	r := gin.New()

	r.Use(middleware.Logger(rotatingLogger))
	r.GET("/health", health)

	r.Run(":8080")
}
func health(c *gin.Context) {
	c.String(http.StatusOK, "hello world!")
}
