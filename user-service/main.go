package main

import (
    "github.com/gin-gonic/gin"
    "github.com/bastosanaa/sistemaAntifurtoAPI/user-service/handlers"
)

func main() {
    handlers.InitDatabase()
    defer handlers.DB.Close()

    r := gin.Default()

    r.POST("/users", handlers.CreateUser)
    r.GET("/users/:id", handlers.GetUser)
    r.PUT("/users/:id", handlers.UpdateUser)
    r.DELETE("/users/:id", handlers.DeleteUser)
    r.GET("/users/health", handlers.Health)

    r.Run(":8001")
}
