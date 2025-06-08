package main

import (
	"github.com/bastosanaa/sistemaAntifurtoAPI/control-service/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	handlers.InitDatabase()
	defer handlers.DB.Close()

	r := gin.Default()
	r.Use(handlers.AuthMiddleware(), handlers.LoggingMiddleware())

	r.POST("/controls/:id/arm", handlers.ArmAlarm)
	r.POST("/controls/:id/disarm", handlers.DisarmAlarm)
	r.GET("/controls/:id/status", handlers.GetStatus)
	r.GET("/controls/health", handlers.Health)

	r.Run(":8005")
}
