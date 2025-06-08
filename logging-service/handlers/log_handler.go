package handlers

import (
	"net/http"
	"time"

	"github.com/bastosanaa/sistemaAntifurtoAPI/logging-service/models"
	"github.com/gin-gonic/gin"
)

// CreateLog trata POST /logs e registra um novo evento.
func CreateLog(c *gin.Context) {
	var input struct {
		Service   string  `json:"service"`
		AlarmID   int64   `json:"alarm_id"`
		UserID    *int64  `json:"user_id"`
		Action    string  `json:"action"`
		Mode      *string `json:"mode"`
		Point     *string `json:"point"`
		Timestamp string  `json:"timestamp"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Service != "control" && input.Service != "trigger" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "service inválido"})
		return
	}

	if input.Action != "arm" && input.Action != "disarm" && input.Action != "trigger" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "action inválido"})
		return
	}

	ts, err := time.Parse(time.RFC3339, input.Timestamp)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "timestamp inválido"})
		return
	}

	res, err := DB.Exec(`INSERT INTO logs (service, alarm_id, user_id, action, mode, point, timestamp)
        VALUES (?, ?, ?, ?, ?, ?, ?)`,
		input.Service, input.AlarmID, input.UserID, input.Action, input.Mode, input.Point, ts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	id, _ := res.LastInsertId()

	entry := models.LogEntry{
		ID:        id,
		Service:   input.Service,
		AlarmID:   input.AlarmID,
		UserID:    input.UserID,
		Action:    input.Action,
		Mode:      input.Mode,
		Point:     input.Point,
		Timestamp: ts,
	}
	c.JSON(http.StatusCreated, entry)
}

// Health responde ao GET /health.
func Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "logging-service OK"})
}
