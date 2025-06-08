package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/bastosanaa/sistemaAntifurtoAPI/notification-service/models"
	"github.com/gin-gonic/gin"
)

const (
	userServiceURL    = "http://localhost:8001"
	loggingServiceURL = "http://localhost:8006/logs"
)

// Health responde ao GET /health.
func Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "notification-service OK"})
}

// Notify trata POST /notify e simula o envio de uma notificação.
func Notify(c *gin.Context) {
	var input struct {
		UserID    int64  `json:"user_id"`
		AlarmID   int64  `json:"alarm_id"`
		Event     string `json:"event"`
		Timestamp string `json:"timestamp"`
	}
	
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Event != "arm" && input.Event != "disarm" && input.Event != "trigger" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "event inválido"})
		return
	}
	t, err := time.Parse(time.RFC3339, input.Timestamp)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "timestamp inválido"})
		return
	}

	phone, statusCode, err := fetchPhone(input.UserID)
	if err != nil {
		if statusCode == http.StatusBadGateway {
			c.JSON(http.StatusBadGateway, gin.H{"error": "falha ao chamar user-service"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if statusCode == http.StatusNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuário não encontrado"})
		return
	}

	msg := formatMessage(input.AlarmID, input.Event, t)
	log.Printf("[NOTIFY] Para: %s | Alarme: %d | Evento: %s | Mensagem: %s", phone, input.AlarmID, input.Event, msg)

	not := models.Notification{
		UserID:    input.UserID,
		AlarmID:   input.AlarmID,
		Event:     input.Event,
		Message:   msg,
		Timestamp: t,
	}
	sendLog(not)

	c.JSON(http.StatusOK, gin.H{"sent": true})
}

func fetchPhone(id int64) (string, int, error) {
	resp, err := http.Get(fmt.Sprintf("%s/users/%d", userServiceURL, id))
	if err != nil {
		return "", http.StatusBadGateway, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return "", http.StatusNotFound, nil
	}
	if resp.StatusCode != http.StatusOK {
		return "", http.StatusBadGateway, fmt.Errorf("user-service status %d", resp.StatusCode)
	}
	var data struct {
		Phone string `json:"phone"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", http.StatusInternalServerError, err
	}
	return data.Phone, http.StatusOK, nil
}

func sendLog(n models.Notification) {
	body, err := json.Marshal(n)
	if err != nil {
		return
	}
	http.Post(loggingServiceURL, "application/json", bytes.NewBuffer(body))
}

func formatMessage(alarmID int64, event string, t time.Time) string {
	var action string
	switch event {
	case "arm":
		action = "ativado"
	case "disarm":
		action = "desativado"
	case "trigger":
		action = "disparado"
	default:
		action = event
	}
	return fmt.Sprintf("O Alarme de número %d foi %s (%s)", alarmID, action, t.Format("02/01/2006 - 15:04"))
}
