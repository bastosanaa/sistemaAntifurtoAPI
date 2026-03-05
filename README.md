# Sistema de Monitoramento Antifurto (Backend API)
Este projeto consiste em um ecossistema de microsserviços desenvolvido em Go para o controle, monitoramento e gerenciamento de sistemas de segurança residencial e comercial. A arquitetura foi projetada para ser acessada tanto por aplicativos móveis quanto por painéis físicos de centrais de alarme, garantindo escalabilidade e isolamento de responsabilidades.

<img width="1314" height="720" alt="image" src="https://github.com/user-attachments/assets/64f2c276-e7d6-4044-a11f-1166cf70a2dd" />


## 🏛️ Arquitetura do Sistema
O backend é estruturado em torno de um API Gateway central que gerencia as requisições e as distribui entre seis microsserviços especializados, cada um com sua própria base de dados SQLite independente:

- **User Service (Porta 8001)**: Responsável pelo ciclo de vida dos usuários, mantendo registros de nomes e números de celular para notificações.

- **Alarm Service (Porta 8002)**: Gerencia o inventário de alarmes, incluindo locais de instalação, permissões de acesso por usuário e o mapeamento de pontos monitorados (ex: sala, quarto, porta principal).

- **Control Service (Porta 8005)**: Centraliza a inteligência de acionamento (armar/desarmar). Valida permissões e estados atuais antes de registrar mudanças solicitadas via aplicativo ou central física.

- **Trigger Service (Porta 8003)**: Processa e registra eventos de disparo originados por sensores de abertura ou presença em pontos específicos.

- **Notification Service (Porta 8004)**: Atua de forma stateless para simular o envio de alertas críticos aos usuários via terminal sempre que ocorre um acionamento, desativação ou disparo.

- **Logging Service (Porta 8006)**: Consolida o histórico completo de auditoria do sistema, persistindo cada interação de controle e cada alerta gerado para consultas futuras.


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
