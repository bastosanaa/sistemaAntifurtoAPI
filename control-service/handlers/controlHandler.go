package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/bastosanaa/sistemaAntifurtoAPI/control-service/models"
	"github.com/gin-gonic/gin"
)

const notificationServiceURL = "http://localhost:8004/notify"

// ControlRequest representa o body para armar/desarmar
type ControlRequest struct {
	UserID *int64 `json:"user_id,omitempty"`
	Source string `json:"source,omitempty"`
	Mode   string `json:"mode"`
}

func ArmAlarm(c *gin.Context) { handleControl(c, "arm") }

func DisarmAlarm(c *gin.Context) { handleControl(c, "disarm") }

func handleControl(c *gin.Context, action string) {
	idParam := c.Param("id")
	alarmID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Alarme inválido"})
		return
	}

	// verifica se alarme existe
	if !alarmExists(alarmID) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Alarme não encontrado"})
		return
	}

	var req ControlRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// verifica permissão se user_id informado
	if req.UserID != nil {
		if !userAllowed(alarmID, *req.UserID) {
			recordAttempt(alarmID, req.UserID, req.Source, req.Mode, action, "failure")
			c.JSON(http.StatusForbidden, gin.H{"error": "Sem permissão"})
			return
		}
	}

	currentState, _ := getCurrentState(alarmID)

	if action == "arm" && currentState == "armed" {
		recordAttempt(alarmID, req.UserID, req.Source, req.Mode, action, "failure")
		c.JSON(http.StatusConflict, gin.H{"error": "Alarme já armado"})
		return
	}
	if action == "disarm" && currentState == "disarmed" {
		recordAttempt(alarmID, req.UserID, req.Source, req.Mode, action, "failure")
		c.JSON(http.StatusConflict, gin.H{"error": "Alarme já desarmado"})
		return
	}

	ctrl := recordAttempt(alarmID, req.UserID, req.Source, req.Mode, action, "success")
	sendLog(ctrl)
	sendNotification(ctrl)

	c.JSON(http.StatusOK, ctrl)
}

func GetStatus(c *gin.Context) {
	idParam := c.Param("id")
	alarmID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Alarme inválido"})
		return
	}

	state, ts := getCurrentState(alarmID)
	c.JSON(http.StatusOK, gin.H{"alarm_id": alarmID, "state": state, "timestamp": ts})
}

func Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "control-service OK"})
}

func alarmExists(id int64) bool {
	resp, err := http.Get(fmt.Sprintf("http://localhost:8002/alarms/%d", id))
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

func userAllowed(alarmID, userID int64) bool {
	resp, err := http.Get(fmt.Sprintf("http://localhost:8002/alarms/%d/users", alarmID))
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return false
	}
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Users []int64 `json:"users"`
	}
	if err := json.Unmarshal(body, &data); err != nil {
		return false
	}
	for _, id := range data.Users {
		if id == userID {
			return true
		}
	}
	return false
}

func getCurrentState(alarmID int64) (string, string) {
	var action string
	var ts string
	query := `SELECT action, timestamp FROM controls WHERE alarm_id = ? AND result = 'success' ORDER BY timestamp DESC LIMIT 1`
	err := DB.QueryRow(query, alarmID).Scan(&action, &ts)
	if err == sql.ErrNoRows {
		return "disarmed", ""
	} else if err != nil {
		return "disarmed", ""
	}
	if action == "arm" {
		return "armed", ts
	}
	return "disarmed", ts
}

func recordAttempt(alarmID int64, userID *int64, source, mode, action, result string) models.Control {
	ts := time.Now().UTC().Format(time.RFC3339)

	_, err := DB.Exec(`INSERT INTO controls (alarm_id, user_id, source, mode, action, timestamp, result) VALUES (?, ?, ?, ?, ?, ?, ?)`, alarmID, userID, source, mode, action, ts, result)
	if err != nil {
		// apenas loga o erro
		fmt.Println("Erro ao registrar controle:", err)
	}

	ctrl := models.Control{
		AlarmID:   alarmID,
		UserID:    userID,
		Source:    source,
		Mode:      mode,
		Action:    action,
		Timestamp: ts,
		Result:    result,
	}
	return ctrl
}

func sendLog(ctrl models.Control) {
	payload, _ := json.Marshal(gin.H{
		"service":   "control",
		"alarm_id":  ctrl.AlarmID,
		"user_id":   ctrl.UserID,
		"action":    ctrl.Action,
		"mode":      ctrl.Mode,
		"timestamp": ctrl.Timestamp,
		"result":    ctrl.Result,
	})
	http.Post("http://localhost:8004/logs", "application/json", bytes.NewBuffer(payload))
}

func sendNotification(ctrl models.Control) {
	if ctrl.Result != "success" || ctrl.UserID == nil {
		return
	}
	payload, _ := json.Marshal(gin.H{
		"user_id":   *ctrl.UserID,
		"alarm_id":  ctrl.AlarmID,
		"event":     ctrl.Action,
		"timestamp": ctrl.Timestamp,
	})
	http.Post(notificationServiceURL, "application/json", bytes.NewBuffer(payload))
}
