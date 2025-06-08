package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/bastosanaa/sistemaAntifurtoAPI/alarm-service/models"
	"github.com/gin-gonic/gin"
)

// CreateAlarm trata POST /alarms
// Body JSON: { "location": "Local de Instalação" }
func CreateAlarm(c *gin.Context) {
	var input models.Alarm
	// Lê e valida JSON recebido
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Insere nova linha em alarms
	result, err := DB.Exec("INSERT INTO alarms (location) VALUES (?)", input.Location)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	id, _ := result.LastInsertId()
	input.ID = id
	// Retorna HTTP 201 Created com o alarme criado
	c.JSON(http.StatusCreated, input)
}

// GetAlarm trata GET /alarms/:id
func GetAlarm(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var alarm models.Alarm
	err = DB.QueryRow("SELECT id, location FROM alarms WHERE id = ?", id).
		Scan(&alarm.ID, &alarm.Location)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Alarme não encontrado"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, alarm)
}

// UpdateAlarm trata PUT /alarms/:id
// Body JSON: { "location": "Novo Local" }
func UpdateAlarm(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Verifica existência
	var exists int
	if err := DB.QueryRow("SELECT 1 FROM alarms WHERE id = ?", id).Scan(&exists); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Alarme não encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var input models.Alarm
	// Lê e valida JSON recebido
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Atualiza registro no banco
	_, err = DB.Exec("UPDATE alarms SET location = ? WHERE id = ?", input.Location, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	input.ID = id
	// Retorna HTTP 200 OK com dados atualizados
	c.JSON(http.StatusOK, input)
}

// DeleteAlarm trata DELETE /alarms/:id
func DeleteAlarm(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Busca antes de excluir para retornar
	var alarm models.Alarm
	err = DB.QueryRow("SELECT id, location FROM alarms WHERE id = ?", id).
		Scan(&alarm.ID, &alarm.Location)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Alarme não encontrado"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	_, err = DB.Exec("DELETE FROM alarms WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, alarm)
}

// AddUserToAlarm trata POST /alarms/:id/users
// Body JSON: { "user_id": 42 }
func AddUserToAlarm(c *gin.Context) {
	idParam := c.Param("id")
	alarmID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Alarme inválido"})
		return
	}

	// Verifica se o alarme existe
	var tmp int
	if err := DB.QueryRow("SELECT 1 FROM alarms WHERE id = ?", alarmID).Scan(&tmp); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Alarme não encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Faz bind do body
	var input struct {
		UserID int64 `json:"user_id"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 3) Verifica se a relação já existe (usar variável nova, ex: relExists)
	var relExists int
	err = DB.QueryRow(
		"SELECT 1 FROM alarm_users WHERE alarm_id = ? AND user_id = ?",
		alarmID, input.UserID,
	).Scan(&relExists)
	if err != nil && err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if relExists == 1 {
		c.JSON(http.StatusConflict, gin.H{"error": "Usuário já autorizado para este alarme"})
		return
	}

	// Insere relacionamento no banco
	result, err := DB.Exec("INSERT INTO alarm_users (alarm_id, user_id) VALUES (?, ?)", alarmID, input.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	relID, _ := result.LastInsertId()
	// Retorna HTTP 201 Created com a relação criada
	c.JSON(http.StatusCreated, gin.H{
		"id":       relID,
		"alarm_id": alarmID,
		"user_id":  input.UserID,
	})
}

// RemoveUserFromAlarm trata DELETE /alarms/:id/users/:user_id
// Exemplo de url para apagar pontos com espaço - Porta%20principal
func RemoveUserFromAlarm(c *gin.Context) {
	alarmParam := c.Param("id")
	alarmID, err := strconv.ParseInt(alarmParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Alarme inválido"})
		return
	}
	userParam := c.Param("user_id")
	userID, err := strconv.ParseInt(userParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Usuário inválido"})
		return
	}

	// Verifica se existe o relacionamento
	var exists int
	query := "SELECT id FROM alarm_users WHERE alarm_id = ? AND user_id = ?"
	if err := DB.QueryRow(query, alarmID, userID).Scan(&exists); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Relacionamento não encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Remove registro da tabela alarm_users
	_, err = DB.Exec("DELETE FROM alarm_users WHERE alarm_id = ? AND user_id = ?", alarmID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Retorna HTTP 200 OK com IDs removidos
	c.JSON(http.StatusOK, gin.H{"alarm_id": alarmID, "user_id": userID})
}

// AddAlarmPoint trata POST /alarms/:id/points
// Body JSON: { "point": "Sala" }
func AddAlarmPoint(c *gin.Context) {
	idParam := c.Param("id")
	alarmID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Alarme inválido"})
		return
	}

	// 1) Verifica se o alarme existe
	var tmp int
	if err := DB.QueryRow("SELECT 1 FROM alarms WHERE id = ?", alarmID).Scan(&tmp); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Alarme não encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Lê e valida JSON recebido
	var input struct {
		Point string `json:"point"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 3) Verifica se já existe ponto com mesmo nome no mesmo alarme
	var exists int
	err = DB.QueryRow(
		"SELECT 1 FROM alarm_points WHERE alarm_id = ? AND point = ?",
		alarmID, input.Point,
	).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if exists == 1 {
		c.JSON(http.StatusConflict, gin.H{"error": "Ponto já cadastrado para este alarme"})
		return
	}

	// Insere novo ponto monitorado no banco
	result, err := DB.Exec(
		"INSERT INTO alarm_points (alarm_id, point) VALUES (?, ?)",
		alarmID, input.Point,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	relID, _ := result.LastInsertId()
	// Retorna HTTP 201 Created com o ponto criado
	c.JSON(http.StatusCreated, gin.H{
		"id":       relID,
		"alarm_id": alarmID,
		"point":    input.Point,
	})
}

// RemoveAlarmPoint trata DELETE /alarms/:id/points/:point

func RemoveAlarmPoint(c *gin.Context) {
	// Lê o "id" do alarme
	alarmParam := c.Param("id")
	alarmID, err := strconv.ParseInt(alarmParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Alarme inválido"})
		return
	}

	// Lê o "point" (nome do ponto) da URL
	pointName := c.Param("point")
	if pointName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ponto inválido"})
		return
	}

	// Verifica se existe um ponto com esse nome para o alarme
	var existsID int64
	query := "SELECT id FROM alarm_points WHERE alarm_id = ? AND point = ?"
	if err := DB.QueryRow(query, alarmID, pointName).Scan(&existsID); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Ponto não encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Remove o ponto usando alarme e point (nome)
	_, err = DB.Exec("DELETE FROM alarm_points WHERE alarm_id = ? AND point = ?", alarmID, pointName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Retorna HTTP 200 OK após remoção
	c.JSON(http.StatusOK, gin.H{
		"alarm_id": alarmID,
		"point":    pointName,
	})
}

// ListAlarms trata GET /alarms
func ListAlarms(c *gin.Context) {
	// Consulta todos os alarmes cadastrados
	rows, err := DB.Query("SELECT id, location FROM alarms")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var alarms []models.Alarm
	for rows.Next() {
		var a models.Alarm
		if err := rows.Scan(&a.ID, &a.Location); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		alarms = append(alarms, a)
	}
	// Retorna HTTP 200 OK com lista de alarmes
	c.JSON(http.StatusOK, alarms)
}

// ListUsersForAlarm trata GET /alarms/:id/users
func ListUsersForAlarm(c *gin.Context) {
	alarmParam := c.Param("id")
	alarmID, err := strconv.ParseInt(alarmParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Alarme inválido"})
		return
	}

	// Verifica se o alarme existe
	var exists int
	if err := DB.QueryRow("SELECT 1 FROM alarms WHERE id = ?", alarmID).Scan(&exists); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Alarme não encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Busca IDs de usuários autorizados
	rows, err := DB.Query("SELECT user_id FROM alarm_users WHERE alarm_id = ?", alarmID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var users []int64
	for rows.Next() {
		var uid int64
		if err := rows.Scan(&uid); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		users = append(users, uid)
	}
	// Retorna HTTP 200 OK com lista de usuários
	c.JSON(http.StatusOK, gin.H{"alarm_id": alarmID, "users": users})
}

// ListPointsForAlarm trata GET /alarms/:id/points
// Retorna apenas uma lista de strings com o campo "point" para o alarme indicado.
func ListPointsForAlarm(c *gin.Context) {
	alarmParam := c.Param("id")
	alarmID, err := strconv.ParseInt(alarmParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Alarme inválido"})
		return
	}

	// Verifica se o alarme existe
	var tmp int
	if err := DB.QueryRow("SELECT 1 FROM alarms WHERE id = ?", alarmID).Scan(&tmp); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Alarme não encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Busca apenas o campo "point" (string) em alarm_points
	rows, err := DB.Query("SELECT point FROM alarm_points WHERE alarm_id = ?", alarmID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var points []string
	for rows.Next() {
		var p string
		if err := rows.Scan(&p); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		points = append(points, p)
	}

	// Retorna HTTP 200 OK com lista de pontos
	c.JSON(http.StatusOK, gin.H{
		"alarm_id": alarmID,
		"points":   points,
	})
}

func Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "alarm-service OK"})
}
