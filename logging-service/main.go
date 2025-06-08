package main

import (
	"github.com/bastosanaa/sistemaAntifurtoAPI/logging-service/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	handlers.InitDatabase()
	defer handlers.DB.Close()

	r := gin.Default()
	r.POST("/logs", handlers.CreateLog)
	r.GET("/logs", handlers.ListLogs)
	r.GET("/logs/health", handlers.Health)

	r.Run(":8006")
}
