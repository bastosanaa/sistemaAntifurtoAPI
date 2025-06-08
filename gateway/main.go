package main

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/bastosanaa/sistemaAntifurtoAPI/gateway/handlers"
	"github.com/bastosanaa/sistemaAntifurtoAPI/gateway/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.AuthMiddleware())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "gateway OK"})
	})

	userURL := getenv("USER_SERVICE_URL")
	alarmURL := getenv("ALARM_SERVICE_URL")
	controlURL := getenv("CONTROL_SERVICE_URL")
	triggerURL := getenv("TRIGGER_SERVICE_URL")
	notificationURL := getenv("NOTIFICATION_SERVICE_URL")
	loggingURL := getenv("LOGGING_SERVICE_URL")

	// Set timeout
	http.DefaultTransport.(*http.Transport).ResponseHeaderTimeout = 3 * time.Second

	r.Any("/users/*proxyPath", handlers.ReverseProxy(userURL))
	r.Any("/alarms/*proxyPath", handlers.ReverseProxy(alarmURL))
	r.Any("/controls/*proxyPath", handlers.ReverseProxy(controlURL))
	r.Any("/triggers/*proxyPath", handlers.ReverseProxy(triggerURL))
	r.POST("/notify", handlers.ReverseProxy(notificationURL))
	r.POST("/logs", handlers.ReverseProxy(loggingURL))

	r.Run(":8000")
}

func getenv(key string) string {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		panic("missing env " + key)
	}
	return v
}
