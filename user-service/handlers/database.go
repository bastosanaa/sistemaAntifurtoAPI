package handlers

import (
    "database/sql"
    "log"

    // importa o driver SQLite, mas o "_" faz apenas o registro do driver
    _ "github.com/mattn/go-sqlite3"
)

// DB é a conexão global com o banco SQLite que os handlers vão usar.
var DB *sql.DB

// InitDatabase abre (ou cria) o arquivo users.db e cria a tabela se não existir.
// Deve ser chamada *antes* de registrar rotas, no main().
func InitDatabase() {
    var err error
    // Ao abrir, se o arquivo não existir, o SQLite criará automaticamente.
    DB, err = sql.Open("sqlite3", "./users.db")
    if err != nil {
        log.Fatalf("Falha ao abrir banco de dados: %v", err)
    }

    // Cria a tabela `users` caso ela ainda não exista:
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

    // Se quiser ver um log sempre que conectar (opcional):
    log.Println("BD inicializado em ./user-service/users.db e tabela 'users' verificada.")
}
