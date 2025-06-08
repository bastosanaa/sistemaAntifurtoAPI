package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/bastosanaa/sistemaAntifurtoAPI/user-service/models"
	"github.com/gin-gonic/gin"
)

// CreateUser trata POST /users e cria um novo usuário.
func CreateUser(c *gin.Context) {
	// Lê e valida JSON recebido
	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. Verifica se já existe usuário com este telefone
	var exists int
	err := DB.QueryRow("SELECT 1 FROM users WHERE phone = ?", input.Phone).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		// Algo deu errado na consulta (exceto "nenhuma linha"), retornar 500
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if exists == 1 {
		// Já existe usuário com esse telefone, retorno 409 Conflict
		c.JSON(http.StatusConflict, gin.H{"error": "Telefone já cadastrado"})
		return
	}

	// Insere registro no banco de dados
	result, err := DB.Exec("INSERT INTO users (name, phone) VALUES (?, ?)", input.Name, input.Phone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 4. Obtém o ID gerado pelo banco para o novo registro.
	id, _ := result.LastInsertId()
	input.ID = id

	// Retorna HTTP 201 Created com o objeto criado
	c.JSON(http.StatusCreated, input)
}

// GetUser trata GET /users/:id e retorna os dados de um usuário pelo ID.
func GetUser(c *gin.Context) {
	// 1. Extrai o parâmetro ":id" da URL e converte para inteiro.
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// 2. Declara uma variável para armazenar o resultado.
	var user models.User

	// 3. Consulta o banco pelo ID.
	err = DB.QueryRow("SELECT id, name, phone FROM users WHERE id = ?", id).
		Scan(&user.ID, &user.Name, &user.Phone)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuário não encontrado"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 4. Se tudo OK, retorna o usuário em JSON.
	c.JSON(http.StatusOK, user)
}

// UpdateUser trata PUT /users/:id e atualiza os dados de um usuário.
func UpdateUser(c *gin.Context) {
	// 1. Lê o "id" da URL.
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Lê e valida JSON recebido
	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 3. Verifica se o usuário existe.
	var exists int
	if err := DB.QueryRow("SELECT 1 FROM users WHERE id = ?", id).Scan(&exists); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Usuário não encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 4. Executa o UPDATE no banco (name e phone).
	_, err = DB.Exec("UPDATE users SET name = ?, phone = ? WHERE id = ?", input.Name, input.Phone, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 5. Retorna o objeto atualizado (com o mesmo ID).
	input.ID = id
	c.JSON(http.StatusOK, input)
}

// DeleteUser trata DELETE /users/:id e remove o usuário.
func DeleteUser(c *gin.Context) {
	// 1. Lê o "id" da URL.
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// 2. Primeiro, busca o usuário no banco para poder retorná-lo depois.
	var user models.User
	err = DB.QueryRow("SELECT id, name, phone FROM users WHERE id = ?", id).
		Scan(&user.ID, &user.Name, &user.Phone)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Usuário não encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 3. Executa o DELETE no banco.
	_, err = DB.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Retorna HTTP 200 OK com usuário removido
	c.JSON(http.StatusOK, user)
}

// Health retorna um JSON simples para sabermos que o user-service está no ar.
// Rota: GET /health
func Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "user-service OK"})
}
