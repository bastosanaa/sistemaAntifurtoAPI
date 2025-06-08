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


---

## Sumário
| Método | Rota                               | Descrição                         |
| ------ | ---------------------------------- | --------------------------------- |
| GET    | `/health`                          | Verifica se o serviço está OK     |
| POST   | `/triggers`                        | Registra um disparo de alarme     |
| GET    | `/alarms/{alarm_id}/triggers`      | Lista disparos de um alarme       |

## Regras de negócio adicionadas
- `event` deve ser `open` ou `presence`.
- `alarm_id` e `point` devem existir no alarm-service.
- Todos os disparos são registrados sem verificação de duplicidade.

---

# Documentação — **control-service**

> Micro-serviço em **Go + Gin** responsável por armar e desarmar alarmes.
> Porta padrão: **http://localhost:8003**


| Método | Rota                           | Descrição                             |
| ------ | ------------------------------ | ------------------------------------- |
| POST   | `/controls/{alarm_id}/arm`    | Arma o alarme indicado                |
| POST   | `/controls/{alarm_id}/disarm` | Desarma o alarme indicado             |
| GET    | `/controls/{alarm_id}/status` | Consulta o estado atual do alarme     |
| GET    | `/health`                     | Verifica se o serviço está OK         |

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
