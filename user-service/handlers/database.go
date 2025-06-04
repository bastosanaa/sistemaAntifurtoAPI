package handlers

import (
    "database/sql"
    "log"

    _ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// InitDatabase abre (ou cria) o arquivo users.db e cria a tabela se não existir.
func InitDatabase() {
    var err error
    DB, err = sql.Open("sqlite3", "./users.db")
    if err != nil {
        log.Fatalf("Falha ao abrir banco de dados: %v", err)
    }

    createTable := `
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        phone TEXT NOT NULL UNIQUE
    );`
    _, err = DB.Exec(createTable)
    if err != nil {
        log.Fatalf("Erro ao criar tabela users: %v", err)
    }

    log.Println("BD inicializado em ./user-service/users.db e tabela 'users' verificada.")
}
