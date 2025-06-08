package main

import (
	"net/http"
	"time"

	"github.com/bastosanaa/sistemaAntifurtoAPI/gateway/handlers"
	"github.com/bastosanaa/sistemaAntifurtoAPI/gateway/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	r.Use(middleware.CORSMiddleware())
	// r.Use(middleware.AuthMiddleware())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "gateway OK"})
	})

	userURL := "http://localhost:8001"
	alarmURL := "http://localhost:8002"
	controlURL := "http://localhost:8005"
	triggerURL := "http://localhost:8003"
	notificationURL := "http://localhost:8004"
	loggingURL := "http://localhost:8006"

	// Set timeout
	http.DefaultTransport.(*http.Transport).ResponseHeaderTimeout = 3 * time.Second

	r.Any("/users", handlers.ReverseProxy(userURL))
	r.Any("/users/*proxyPath", handlers.ReverseProxy(userURL))
	r.Any("/alarms", handlers.ReverseProxy(alarmURL))
	r.Any("/alarms/*proxyPath", handlers.ReverseProxy(alarmURL))
	r.Any("/controls", handlers.ReverseProxy(controlURL))
	r.Any("/controls/*proxyPath", handlers.ReverseProxy(controlURL))
	r.Any("/triggers", handlers.ReverseProxy(triggerURL))
	r.Any("/triggers/*proxyPath", handlers.ReverseProxy(triggerURL))
	r.POST("/notify", handlers.ReverseProxy(notificationURL))
	r.Any("/logs", handlers.ReverseProxy(loggingURL))

	r.Run(":8000")
}
