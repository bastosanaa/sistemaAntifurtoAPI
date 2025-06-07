package handlers

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDatabase() {
	var err error
	DB, err = sql.Open("sqlite3", "./controls.db")
	if err != nil {
		log.Fatalf("Falha ao abrir banco de dados: %v", err)
	}

	createTable := `
    CREATE TABLE IF NOT EXISTS controls (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        alarm_id INTEGER NOT NULL,
        user_id INTEGER,
        source TEXT,
        mode TEXT NOT NULL,
        action TEXT NOT NULL,
        timestamp DATETIME NOT NULL,
        result TEXT NOT NULL
    );`
	if _, err := DB.Exec(createTable); err != nil {
		log.Fatalf("Erro ao criar tabela controls: %v", err)
	}

	log.Println("Banco de controles inicializado em ./control-service/controls.db")
}
