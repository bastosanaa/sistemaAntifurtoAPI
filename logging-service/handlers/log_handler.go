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
	// Lê e valida JSON recebido
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verifica que service é "control" ou "trigger"
	if input.Service != "control" && input.Service != "trigger" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "service inválido"})
		return
	}

	// Validações do campo action
	if input.Action != "arm" && input.Action != "disarm" && input.Action != "trigger" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "action inválido"})
		return
	}

	// Converte timestamp recebido para time.Time
	ts, err := time.Parse(time.RFC3339, input.Timestamp)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "timestamp inválido"})
		return
	}

	// Insere registro no banco local logs.db
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
	// Retorna HTTP 201 Created com o log gravado
	c.JSON(http.StatusCreated, entry)
}

// ListLogs trata GET /logs e retorna todos os eventos registrados.
func ListLogs(c *gin.Context) {
	// Consulta todos os logs no banco SQLite
	rows, err := DB.Query(`
        SELECT
            id, service, alarm_id, user_id, action, mode, point, timestamp
        FROM logs
        ORDER BY timestamp DESC
    `)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var entries []models.LogEntry
	for rows.Next() {
		var e models.LogEntry
		// Faz o scan de cada linha para struct LogEntry
		if err := rows.Scan(
			&e.ID,
			&e.Service,
			&e.AlarmID,
			&e.UserID,
			&e.Action,
			&e.Mode,
			&e.Point,
			&e.Timestamp,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		entries = append(entries, e)
	}
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Retorna HTTP 200 OK com a lista de logs
	c.JSON(http.StatusOK, entries)
}

// Health responde ao GET /health.
func Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "logging-service OK"})
}
