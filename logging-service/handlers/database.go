package handlers

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// InitDatabase abre ou cria logs.db e cria a tabela logs.
func InitDatabase() {
	var err error
	DB, err = sql.Open("sqlite3", "./logs.db")
	if err != nil {
		log.Fatalf("Falha ao abrir banco de dados: %v", err)
	}

	createTable := `
    CREATE TABLE IF NOT EXISTS logs (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        service TEXT NOT NULL,
        alarm_id INTEGER NOT NULL,
        user_id INTEGER,
        action TEXT NOT NULL,
        mode TEXT,
        point TEXT,
        timestamp DATETIME NOT NULL
    );`

	if _, err := DB.Exec(createTable); err != nil {
		log.Fatalf("Erro ao criar tabela logs: %v", err)
	}

	log.Println("Banco logs.db inicializado e tabela 'logs' verificada.")
}
