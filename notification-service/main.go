package main

import (
	"github.com/bastosanaa/sistemaAntifurtoAPI/notification-service/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	handlers.InitDatabase()

	r := gin.Default()

	r.GET("/health", handlers.Health)
	r.POST("/notify", handlers.Notify)

	r.Run(":8004")
}
