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