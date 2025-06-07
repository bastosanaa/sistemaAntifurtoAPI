package handlers

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// InitDatabase abre ou cria triggers.db e garante a tabela triggers.
func InitDatabase() {
	var err error
	DB, err = sql.Open("sqlite3", "./triggers.db")
	if err != nil {
		log.Fatalf("Falha ao abrir banco de dados: %v", err)
	}

	if _, err := DB.Exec(`PRAGMA foreign_keys = ON;`); err != nil {
		log.Fatalf("Erro ao habilitar foreign_keys: %v", err)
	}

	createTable := `
    CREATE TABLE IF NOT EXISTS triggers (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        alarm_id INTEGER NOT NULL,
        point TEXT NOT NULL,
        event TEXT NOT NULL,
        timestamp DATETIME NOT NULL,
        FOREIGN KEY(alarm_id) REFERENCES alarms(id) ON DELETE CASCADE
    );`
	if _, err := DB.Exec(createTable); err != nil {
		log.Fatalf("Erro ao criar tabela triggers: %v", err)
	}

	log.Println("Banco triggers.db inicializado e tabela 'triggers' verificada.")
}
