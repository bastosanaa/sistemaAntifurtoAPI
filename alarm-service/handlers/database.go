package handlers

import (
    "database/sql"
    "log"

    _ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// InitDatabase abre (ou cria) o arquivo alarms.db e configura as tabelas:
// - alarms: locais de instalação
// - alarm_users: usuários com permissão para cada alarme
// - alarm_points: pontos monitorados de cada alarme (sala, quarto, etc.)
func InitDatabase() {
    var err error
    // Abre ou cria o arquivo "alarms.db" dentro da pasta alarm-service
    DB, err = sql.Open("sqlite3", "./alarms.db")
    if err != nil {
        log.Fatalf("Falha ao abrir banco de dados: %v", err)
    }

    // Habilita enforcement de chaves estrangeiras no SQLite
    _, err = DB.Exec(`PRAGMA foreign_keys = ON;`)
    if err != nil {
        log.Fatalf("Erro ao ativar foreign_keys: %v", err)
    }

    // Tabela Alarms
    createAlarms := `
    CREATE TABLE IF NOT EXISTS alarms (
        id        INTEGER PRIMARY KEY AUTOINCREMENT,
        location  TEXT    NOT NULL
    );`
    if _, err := DB.Exec(createAlarms); err != nil {
        log.Fatalf("Erro ao criar tabela alarms: %v", err)
    }

    // tabela AlarmUsers - usuários com permissão de acesso 
    createAlarmUsers := `
    CREATE TABLE IF NOT EXISTS alarm_users (
        id       INTEGER PRIMARY KEY AUTOINCREMENT,
        alarm_id INTEGER NOT NULL,
        user_id  INTEGER NOT NULL,
        FOREIGN KEY(alarm_id) REFERENCES alarms(id) ON DELETE CASCADE
    );`
    if _, err := DB.Exec(createAlarmUsers); err != nil {
        log.Fatalf("Erro ao criar tabela alarm_users: %v", err)
    }

    // tabela AlarmPoints - pontos monitorados
    createAlarmPoints := `
    CREATE TABLE IF NOT EXISTS alarm_points (
        id       INTEGER PRIMARY KEY AUTOINCREMENT,
        alarm_id INTEGER NOT NULL,
        point     TEXT    NOT NULL,
        FOREIGN KEY(alarm_id) REFERENCES alarms(id) ON DELETE CASCADE
    );`
    if _, err := DB.Exec(createAlarmPoints); err != nil {
        log.Fatalf("Erro ao criar tabela alarm_points: %v", err)
    }

    log.Println("Banco inicializado e tabelas de alarmes, alarm_users e alarm_points criadas/validadas.")
}
