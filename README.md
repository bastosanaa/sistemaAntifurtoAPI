
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
