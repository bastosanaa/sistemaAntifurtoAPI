# Documentação — **user-service**

> Micro-serviço em **Go + Gin** responsável por CRUD de usuários.  
> Porta padrão: **http://localhost:8001**

---

## Sumário

| Método | Rota                     | Descrição                     |
| ------ | ------------------------ | ----------------------------- |
| GET    | `/health`               | Verifica se o serviço está OK |
| POST   | `/users`                | Cria um novo usuário          |
| GET    | `/users/{id}`           | Busca usuário por ID          |
| PUT    | `/users/{id}`           | Atualiza usuário existente    |
| DELETE | `/users/{id}`           | Remove usuário                |

---

## Regras de negócio adicionadas
- Um mesmo telefone pode ser adicionado apenas uma vez

---

# Documentação — **alarm-service**

> Micro-serviço em **Go + Gin** responsável por CRUD de alarmes, gerenciamento de usuários autorizados e pontos monitorados.  
> Porta padrão: **http://localhost:8002**

---

## Sumário

| Método | Rota                                | Descrição                                   |
| ------ | ----------------------------------- | ------------------------------------------- |
| POST   | `/alarms`                           | Cria um novo alarme                         |
| GET    | `/alarms/{id}`                      | Busca um alarme por ID                      |
| PUT    | `/alarms/{id}`                      | Atualiza um alarme existente                |
| DELETE | `/alarms/{id}`                      | Remove um alarme e retorna o objeto excluído |
| GET    | `/alarms`                           | Lista todos os alarmes                      |
| POST   | `/alarms/{id}/users`                | Adiciona um usuário autorizado a um alarme  |
| GET    | `/alarms/{id}/users`                | Lista usuários autorizados de um alarme     |
| DELETE | `/alarms/{id}/users/{user_id}`      | Remove um usuário autorizado de um alarme   |
| POST   | `/alarms/{id}/points`               | Adiciona um ponto monitorado a um alarme    |
| GET    | `/alarms/{id}/points`               | Lista pontos monitorados de um alarme       |
| DELETE | `/alarms/{id}/points/{point_id}`    | Remove um ponto monitorado de um alarme     |

## Regras de negócio adicionadas
- O alarme deve existir para criar uma relação alarmUser
- Usuário não pode ter mais de uma autorização para o mesmo alarme
- Um ponto supervisionado não pode ser cadastrado mais de uma vez para o mesmo alarme

---


# Documentação — **trigger-service**

> Micro-serviço em **Go + Gin** responsável por registrar disparos de alarmes.
> Porta padrão: **http://localhost:8003**


---

## Sumário
| Método | Rota                               | Descrição                         |
| ------ | ---------------------------------- | --------------------------------- |
| GET    | `/health`                          | Verifica se o serviço está OK     |
| POST   | `/triggers`                        | Registra um disparo de alarme     |
| GET    | `/alarms/{alarm_id}/triggers`      | Lista disparos de um alarme       |

Disparos geram notificações via **notification-service** para todos os usuários autorizados.

## Regras de negócio adicionadas
- `event` deve ser `open` ou `presence`.
- `alarm_id` e `point` devem existir no alarm-service.
- Todos os disparos são registrados sem verificação de duplicidade.

---

# Documentação — **control-service**

> Micro-serviço em **Go + Gin** responsável por armar e desarmar alarmes.
> Porta padrão: **http://localhost:8005**


| Método | Rota                           | Descrição                             |
| ------ | ------------------------------ | ------------------------------------- |
| POST   | `/controls/{alarm_id}/arm`    | Arma o alarme indicado                |
| POST   | `/controls/{alarm_id}/disarm` | Desarma o alarme indicado             |
| GET    | `/controls/{alarm_id}/status` | Consulta o estado atual do alarme     |
| GET    | `/health`                     | Verifica se o serviço está OK         |
Ao armar ou desarmar com sucesso, o serviço envia uma notificação via **notification-service** para o usuário informado.

## Exemplo de corpo para arm/disarm

```json
{
  "user_id": 1,
  "mode": "app"
}
```

---

# Documentação — **notification-service**

> Micro-serviço em **Go + Gin** responsável por enviar notificações aos usuários.
> Porta padrão: **http://localhost:8004**

## Sumário
| Método | Rota      | Descrição                         |
| ------ | --------- | --------------------------------- |
| GET    | `/health` | Verifica se o serviço está OK     |
| POST   | `/notify` | Envia uma notificação ao usuário   |

## Exemplo de corpo para `/notify`

```json
{
  "user_id": 1,
  "alarm_id": 2,
  "event": "arm",
  "timestamp": "2025-06-06T15:00:00Z"
}
```

Mensagem gerada:

```
O Alarme de número 2 foi ativado (06/06/2025 - 15:00)
```

---

# Documentação — **logging-service**

> Micro-serviço em **Go + Gin** responsável por registrar todos os eventos de arm, disarm e trigger.
> Porta padrão: **http://localhost:8006**

## Sumário
| Método | Rota      | Descrição                         |
| ------ | --------- | --------------------------------- |
| GET    | `/health` | Verifica se o serviço está OK     |
| POST   | `/logs`   | Registra um evento de alarme      |

## Exemplo de corpo para `/logs`

```json
{
  "service": "control",
  "alarm_id": 1,
  "user_id": 2,
  "action": "arm",
  "mode": "app",
  "timestamp": "2025-06-08T14:30:00Z"
}
```

---

# Documentação — **gateway**

> API Gateway em **Go + Gin** responsável por unificar o acesso aos microserviços.
> Porta padrão: **http://localhost:8000**

## Sumário

| Método | Rota                      | Encaminha para               |
| ------ | ------------------------- | ---------------------------- |
| GET    | `/health`                | Status do gateway            |
| ALL    | `/users/*`               | user-service                 |
| ALL    | `/alarms/*`              | alarm-service                |
| ALL    | `/controls/*`            | control-service              |
| ALL    | `/triggers/*`            | trigger-service              |
| POST   | `/notify`                | notification-service         |
| POST   | `/logs`                  | logging-service              |


+154
-104

# Guia de Uso via Gateway

Este documento apresenta as rotas disponíveis através do **API Gateway** do sistema Antifurto. Todas as requisições devem ser enviadas para `http://localhost:8000`.

## Autenticação

Inclua o cabeçalho abaixo em todas as chamadas:

```
Authorization: Bearer secret
```

Para as rotas do **control-service** também é necessário:

```
X-Token: secret
```

---

## 1. Usuários

### POST `/users`
Cria um novo usuário.

**Corpo de exemplo**
```json
{
  "name": "Maria",
  "phone": "551199999999"
}
```

### GET `/users/{id}`
Busca um usuário por ID.

### PUT `/users/{id}`
Atualiza um usuário existente.

**Corpo de exemplo**
```json
{
  "name": "Maria",
  "phone": "551188888888"
}
```

### DELETE `/users/{id}`
Remove um usuário.

### GET `/health`
Verifica se o serviço de usuários está disponível.

---

## 2. Alarmes

### POST `/alarms`
Cria um novo alarme.

**Corpo de exemplo**
```json
{
  "location": "Escritório"
}
```

### GET `/alarms/{id}`
Busca um alarme por ID.

### PUT `/alarms/{id}`
Atualiza um alarme existente.

**Corpo de exemplo**
```json
{
  "location": "Loja"
}
```

### DELETE `/alarms/{id}`
Remove um alarme.

### GET `/alarms`
Lista todos os alarmes.

### POST `/alarms/{id}/users`
Autoriza um usuário a operar o alarme.

**Corpo de exemplo**
```json
{
  "user_id": 1
}
```

### GET `/alarms/{id}/users`
Lista usuários autorizados de um alarme.

### DELETE `/alarms/{id}/users/{user_id}`
Remove a autorização de um usuário.

### POST `/alarms/{id}/points`
Cadastra um ponto monitorado.

**Corpo de exemplo**
```json
{
  "point": "Sala"
}
```

### GET `/alarms/{id}/points`
Lista pontos monitorados de um alarme.

### DELETE `/alarms/{id}/points/{point}`
Remove um ponto do alarme.

### GET `/health`
Verifica se o serviço de alarmes está disponível.

---

## 3. Disparos

### POST `/triggers`
Registra um disparo de alarme.

**Corpo de exemplo**
```json
{
  "alarm_id": 1,
  "point": "Sala",
  "event": "open"
}
```

### GET `/alarms/{alarm_id}/triggers`
Lista disparos de um alarme.

### GET `/health`
Verifica se o serviço de disparos está disponível.

---

## 4. Controle de Alarmes

As rotas abaixo exigem também o cabeçalho `X-Token: secret`.

### POST `/controls/{alarm_id}/arm`
Arma o alarme indicado.

**Corpo de exemplo**
```json
{
  "user_id": 1,
  "mode": "app"
}
```

OBS: campo app, serve para identificar se o envio foi feito via central ou via app (ou outro)

### POST `/controls/{alarm_id}/disarm`
Desarma o alarme indicado.

**Corpo de exemplo**
```json
{
  "user_id": 1,
  "mode": "app"
}
```

### GET `/controls/{alarm_id}/status`
Consulta o estado atual do alarme.

### GET `/health`
Verifica se o serviço de controle está disponível.

---

## 5. Notificações

### POST `/notify`
Envia uma notificação para um usuário.

**Corpo de exemplo**
```json
{
  "user_id": 1,
  "alarm_id": 2,
  "event": "arm",
  "timestamp": "2025-06-06T15:00:00Z"
}
```

### GET `/health`
Verifica se o serviço de notificações está disponível.

---

## 6. Logs

### POST `/logs`
Registra um evento no serviço de logs.

**Corpo de exemplo**
```json
{
  "service": "control",
  "alarm_id": 1,
  "user_id": 2,
  "action": "arm",
  "mode": "app",
  "timestamp": "2025-06-08T14:30:00Z"
}
```

### GET `/health`
Verifica se o serviço de logs está disponível.

---

## 7. Gateway

### GET `/health`
Retorna o status do próprio gateway.
