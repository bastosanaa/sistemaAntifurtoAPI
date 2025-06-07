package main

import (
	"github.com/bastosanaa/sistemaAntifurtoAPI/trigger-service/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	handlers.InitDatabase()
	defer handlers.DB.Close()

	r := gin.Default()

	r.POST("/triggers", handlers.CreateTrigger)
	r.GET("/alarms/:alarm_id/triggers", handlers.ListTriggers)
	r.GET("/health", handlers.Health)

	r.Run(":8003")
}
