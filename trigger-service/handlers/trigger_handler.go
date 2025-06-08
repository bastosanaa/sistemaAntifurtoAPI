package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/bastosanaa/sistemaAntifurtoAPI/trigger-service/models"
	"github.com/gin-gonic/gin"
)

const (
	alarmServiceURL        = "http://localhost:8002"
	loggingServiceURL      = "http://localhost:8006/logs"
	notificationServiceURL = "http://localhost:8004/notify"
)

// CreateTrigger trata POST /triggers
func CreateTrigger(c *gin.Context) {
	var input struct {
		AlarmID int64  `json:"alarm_id"`
		Point   string `json:"point"`
		Event   string `json:"event"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Event != "open" && input.Event != "presence" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "event inválido"})
		return
	}

	// Valida existência do alarme
	resp, err := http.Get(fmt.Sprintf("%s/alarms/%d", alarmServiceURL, input.AlarmID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Alarme não encontrado"})
		return
	}
	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao validar alarme"})
		return
	}

	// Valida ponto monitorado
	resp, err = http.Get(fmt.Sprintf("%s/alarms/%d/points", alarmServiceURL, input.AlarmID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if resp.StatusCode == http.StatusNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Alarme não encontrado"})
		return
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao validar ponto"})
		return
	}
	var pointsResp struct {
		AlarmID int64    `json:"alarm_id"`
		Points  []string `json:"points"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&pointsResp); err != nil {
		resp.Body.Close()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	resp.Body.Close()
	exists := false
	for _, p := range pointsResp.Points {
		if p == input.Point {
			exists = true
			break
		}
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ponto não encontrado"})
		return
	}

	timestamp := time.Now().UTC()
	res, err := DB.Exec("INSERT INTO triggers (alarm_id, point, event, timestamp) VALUES (?, ?, ?, ?)",
		input.AlarmID, input.Point, input.Event, timestamp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	id, _ := res.LastInsertId()
	trig := models.Trigger{ID: id, AlarmID: input.AlarmID, Point: input.Point, Event: input.Event, Timestamp: timestamp}

	payload := map[string]interface{}{
		"service":   "trigger",
		"alarm_id":  input.AlarmID,
		"point":     input.Point,
		"action":    "trigger",
		"timestamp": timestamp,
	}
	body, _ := json.Marshal(payload)

	if err := postJSON(loggingServiceURL, body); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := notifyUsers(input.AlarmID, timestamp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, trig)
}

func postJSON(url string, data []byte) error {
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("serviço %s retornou status %d", url, resp.StatusCode)
	}
	return nil
}

func notifyUsers(alarmID int64, ts time.Time) error {
	resp, err := http.Get(fmt.Sprintf("%s/alarms/%d/users", alarmServiceURL, alarmID))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("alarme n\u00e3o encontrado")
		}
		return fmt.Errorf("falha ao obter usuarios")
	}
	var list struct {
		Users []int64 `json:"users"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		return err
	}
	for _, uid := range list.Users {
		payload, _ := json.Marshal(map[string]interface{}{
			"user_id":   uid,
			"alarm_id":  alarmID,
			"event":     "trigger",
			"timestamp": ts,
		})
		if err := postJSON(notificationServiceURL, payload); err != nil {
			return err
		}
	}
	return nil
}

// ListTriggers trata GET /alarms/:alarm_id/triggers
func ListTriggers(c *gin.Context) {
	idParam := c.Param("alarm_id")
	alarmID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Alarme inválido"})
		return
	}

	// valida existência do alarme
	resp, err := http.Get(fmt.Sprintf("%s/alarms/%d", alarmServiceURL, alarmID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Alarme não encontrado"})
		return
	}
	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao validar alarme"})
		return
	}

	query := "SELECT point, event, timestamp FROM triggers WHERE alarm_id = ?"
	args := []interface{}{alarmID}

	if fromStr := c.Query("from"); fromStr != "" {
		from, err := time.Parse(time.RFC3339, fromStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "from inválido"})
			return
		}
		query += " AND timestamp >= ?"
		args = append(args, from)
	}
	if toStr := c.Query("to"); toStr != "" {
		to, err := time.Parse(time.RFC3339, toStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "to inválido"})
			return
		}
		query += " AND timestamp <= ?"
		args = append(args, to)
	}

	rows, err := DB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var triggers []struct {
		Point     string    `json:"point"`
		Event     string    `json:"event"`
		Timestamp time.Time `json:"timestamp"`
	}
	for rows.Next() {
		var t struct {
			Point     string
			Event     string
			Timestamp time.Time
		}
		if err := rows.Scan(&t.Point, &t.Event, &t.Timestamp); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		triggers = append(triggers, struct {
			Point     string    `json:"point"`
			Event     string    `json:"event"`
			Timestamp time.Time `json:"timestamp"`
		}{t.Point, t.Event, t.Timestamp})
	}

	c.JSON(http.StatusOK, triggers)
}

func Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "trigger-service OK"})
}
