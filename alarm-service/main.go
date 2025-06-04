package main

import (
    "github.com/gin-gonic/gin"
    "github.com/bastosanaa/sistemaAntifurtoAPI/alarm-service/handlers"
)

func main() {
    handlers.InitDatabase()
    defer handlers.DB.Close()

    r := gin.Default()

    
    // Rotas para CRUD de Alarmes
    r.POST("/alarms", handlers.CreateAlarm)          // Criar um alarme
    r.GET("/alarms/:id", handlers.GetAlarm)          // Ler um alarme por ID
    r.PUT("/alarms/:id", handlers.UpdateAlarm)       // Atualizar um alarme
    r.DELETE("/alarms/:id", handlers.DeleteAlarm)    // Excluir um alarme
    r.GET("/alarms", handlers.ListAlarms)            // Listar todos os alarmes

    // Rotas para gerenciar usuários autorizados em um alarme
    r.POST("/alarms/:id/users", handlers.AddUserToAlarm)            // Adicionar usuário a alarme
    r.GET("/alarms/:id/users", handlers.ListUsersForAlarm)          // Listar usuários de um alarme
    r.DELETE("/alarms/:id/users/:user_id", handlers.RemoveUserFromAlarm) // Remover usuário de um alarme

    // Rotas para gerenciar pontos monitorados de um alarme
    r.POST("/alarms/:id/points", handlers.AddAlarmPoint)                 // Adicionar ponto a alarme
    r.GET("/alarms/:id/points", handlers.ListPointsForAlarm)             // Listar pontos de um alarme
    r.DELETE("/alarms/:id/points/:point", handlers.RemoveAlarmPoint)  // Remover ponto de um alarme

    r.GET("/health", handlers.Health)



    r.Run(":8002")
}
