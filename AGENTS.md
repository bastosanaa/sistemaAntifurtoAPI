# AGENTS.md

> **Objetivo:**  
> Este documento descreve o contexto, convenções e práticas do projeto **sistemaAntifurtoAPI**, para que o Codex (ou qualquer outro assistente de código) possa atuar como agente auxiliar de forma eficiente.

---

## 1. Visão Geral do Projeto

- **Propósito:** Backend para um sistema de controle e monitoramento de alarmes antifurto.  
- **Tecnologia:**  
  - **Linguagem:** Go (>= 1.20)  
  - **Framework HTTP:** Gin  
  - **Banco de dados:** SQLite (cada microserviço com seu próprio arquivo .db)  

- **Microserviços:**  
  1. **user-service** (porta 8001)  
     - CRUD de usuários  
     - Unicidade de `phone`  
     - `/health`, `/users`  
  2. **alarm-service** (porta 8002)
     - CRUD de alarmes (`location`)
     - Relacionamento `alarm_users` (usuários autorizados)
     - Relacionamento `alarm_points` (pontos monitorados)
     - Comunicação HTTP com `user-service` para buscar nomes dos usuários

  3. **trigger-service** (porta 8003)
     - Registro de disparos de alarme
     - `/triggers`, `/alarms/{id}/triggers`
     - Chama `logging-service` e `notification-service`

  4. **control-service** (porta 8003)
     - Armar e desarmar alarmes
     - Consulta de status `/controls/{alarm_id}/status`

- **Gateway (opcional):**  
  - Proxy REST que encaminha `/users` e `/alarms` entre frontends e microserviços.

---

## 2. Estrutura de Diretórios

sistemaAntifurtoAPI/
├── go.mod
├── gateway/
│   └── main.go
├── user-service/

│ ├── handlers/
│ │ ├── database.go
│ │ └── userHandler.go
│ ├── models/
│ │ └── user.go
│ └── main.go
├── alarm-service/
│ ├── handlers/
│ │ ├── database.go
│ │ └── alarmHandler.go
│ ├── models/
│ │ └── alarm.go
│ └── main.go
└── trigger-service/
    ├── handlers/
    │ ├── database.go
    │ └── trigger_handler.go
    ├── models/
    │ └── trigger.go

│   ├── handlers/
│   │   ├── database.go
│   │   └── userHandler.go
│   ├── models/
│   │   └── user.go
│   └── main.go
├── alarm-service/
│   ├── handlers/
│   │   ├── database.go
│   │   └── alarmHandler.go
│   ├── models/
│   │   └── alarm.go
│   └── main.go
└── control-service/
    ├── handlers/
    │   ├── controlHandler.go
    │   ├── database.go
    │   └── middlewares.go
    ├── models/
    │   └── control.go
    └── main.go


---

## 3. Convenções de Código

- **Pacote e Imports:**  
  - Cada pasta com `main.go` deve usar `package main`.  
  - Handlers em subpastas usam `package handlers` e importam `handlers.DB`.  
  - Models em `package models` estruturam apenas as `structs` e tags JSON.  

- **Banco de Dados:**  
  - Cada microserviço inicializa seu `.db` em `InitDatabase()`.  
  - Habilitar `PRAGMA foreign_keys = ON;` no alarm-service.  
  - Usar `?` para placeholders em `DB.Exec` e `DB.QueryRow`.

- **Erros e Status HTTP:**  
  - Parâmetro inválido → `400 Bad Request`.  
  - Recurso não encontrado → `404 Not Found`.  
  - Conflito de unicidade → `409 Conflict`.  
  - Erro de banco ou decode JSON → `500 Internal Server Error`.  
  - Sucesso de criação → `201 Created`.  
  - Sucesso geral de leitura/atualização → `200 OK`.  
  - Sucesso de deleção (quando retorna objeto) → `200 OK`; sem corpo → `204 No Content`.

- **Naming:**  
  - Handlers: `CreateXxx`, `GetXxx`, `UpdateXxx`, `DeleteXxx`, `ListXxx`.  
  - Rotas REST: plural (`/users`, `/alarms`, `/alarms/:id/users`).

---

## 4. Como o Agente Deve Ajudar

1. **Gerar Código Novo:**  
   - Use as convenções acima, gere somente o que falta (handlers, migrations, middlewares).  
2. **Refatorações:**  
   - Quando solicitado, agrupe imports, corrija variáveis duplicadas e garanta padrões de erro.  
3. **Documentação Automática:**  
   - Ao criar rotas, atualize README.md e AGENTS.md com tabelas e exemplos de requests/responses.  
4. **Testes Manuais:**  
   - Sugira comandos `curl` ou trechos de Postman/Insomnia para validar endpoints.  
5. **Manutenção de Contexto:**  
   - Sempre verifique a estrutura de pastas e o `go.mod` antes de propor alterações.  
   - Caso crie novas tabelas, atualize `InitDatabase()` e o modelo correspondente.  
6. **Boas Práticas:**  
   - Não acople pacotes uns aos outros; respeite a separação de responsabilidades.  
   - Prefira middleware para autenticação/CORS quando o projeto evoluir.

---

## 5. Exemplos de Uso

- **Criar um Handler de Novo Microserviço:**
  ```go
  func CreateThing(c *gin.Context) {
      var input models.Thing
      if err := c.ShouldBindJSON(&input); err != nil {
          c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
          return
      }
      // ...
  }

Dica: Sempre consulte este arquivo antes de gerar ou refatorar código, para manter alinhamento com as convenções do projeto.








